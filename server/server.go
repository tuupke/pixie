package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fasthttp/router"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"github.com/hashicorp/mdns"
	"github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"openticket.tech/crud"
	"openticket.tech/db"
	"openticket.tech/env"
	"openticket.tech/lifecycle/v2"
	"openticket.tech/null"
	"openticket.tech/rest/v3"

	"github.com/tuupke/pixie/packets"
)

type wString string

func (w wString) String() string { return string(w) }

type Problem struct {
	Id       string      `json:"id"`
	Rgb      null.String `json:"rgb"`
	Location Location    `json:"location" gorm:"embedded;embeddedPrefix:loc_"`
}

func (p Problem) Identifier() fmt.Stringer {
	return wString(p.Id)
}

type ExternalData struct {
	crud.SimpleBaseModel
	Username string      `json:"username"`
	UserId   string      `json:"id"`
	Teamname null.String `json:"team"`
	TeamId   null.String `json:"team_id"`
	HostId   null.String `json:"host_id"`
	Location Rotated     `json:"location" gorm:"embedded;embeddedPrefix:loc_"`
}

type Rotated struct {
	Location
	Rotation float64 `json:"rotation"`
}

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Host struct {
	crud.SimpleBaseModel
	Hostname   string    `json:"hostname"`
	PrimaryIp  string    `json:"primary_ip"`
	PrimaryMac string    `json:"primary_mac"`
	Data       []byte    `json:"data"`
	LastSeen   time.Time `json:"last_seen"`
}

var url = env.Get("DOMJUDGE_API_URL", "https://www.domjudge.org/demoweb/api/v4/")
var natsHost, natsPort string
var natsPortInt int
var orm *gorm.DB
var broadcastIPnet *net.IPNet

func init() {
	var err error
	addrs, err := net.InterfaceAddrs()
	log.Err(err).Msg("retrieved interfaces")
	if err != nil {
		log.Fatal().Msg("interfaces is required")
	}

	var natsAddr string

	// Two cases, either an nats-addr is provided, or we need to choose one
	// if it is not provided load the first private ip
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			natsAddr = ipNet.IP.String() + ":4222"
			broadcastIPnet = ipNet
			break
		}
	}

	natsHost, natsPort, err = net.SplitHostPort(natsAddr)
	if err != nil {
		log.Fatal().Err(err).Str("natsAddr", natsAddr).Msg("could not parse natsAddr")
	}

	natsPortInt, err = strconv.Atoi(natsPort)
	if err != nil {
		log.Fatal().Err(err).Str("natsPort", natsPort).Str("natsAddr", natsAddr).Msg("could not convert natsPort to int")
	}

	// For reasons of selecting an ip subnet natsHost must be an ip
	if net.ParseIP(natsHost) == nil {
		log.Fatal().Str("natsHost", natsHost).Msg("natsHost does not appear to be an ip")
	}

	mdnsServe(natsAddr)
}

var settings *Settings

func main() {
	lifecycle.Finally(func() { log.Warn().Msg("Exited") })
	defer lifecycle.PanicHandler()

	ns, err := server.NewServer(&server.Options{
		Host: natsHost,
		Port: natsPortInt,
		// JetStream: true,
	})
	log.Err(err).Str("host", natsHost).Str("port", natsPort).Msg("created nats")
	if err == nil {
		go ns.Start()
		lifecycle.Finally(ns.Shutdown)
	}

	// Give nats some time to start
	time.Sleep(time.Millisecond * 500)

	// Nats connect
	nc, err := nats.Connect(natsHost + ":" + natsPort)
	log.Err(err).Msg("connected nats")
	if err != nil {
		log.Fatal().Msg("nats needs to be connectable")
	}

	sql, _ := db.Conn("DB")
	sql.Close()

	orm, err = gorm.Open(sqlite.Open("./pixie.sqlite"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	log.Err(err).Msg("loaded gorm")
	if err != nil {
		log.Fatal().Msg("gorm must boot")
	}

	_, err = nc.Subscribe("register-a-new-host", func(msg *nats.Msg) {
		log.Debug().Bytes("registration", msg.Data).Msg("received registration")
		log.Info().Msg("handling registration")

		reg := packets.GetRootAsRegister(msg.Data, 0)
		banner := reg.Banner(nil)

		b := flatbuffers.NewBuilder(128)
		hn := b.CreateByteString(banner.Hostname())
		guid := b.CreateByteString(banner.Identifier())

		id, err := uuid.Parse(string(banner.Identifier()))
		log.Err(err).Msg("decoded uuid bytes")
		if err != nil {
			return
		}

		var primaryIp, primaryMac string
		// Check the main ip which matches whatever we are bound to
		var ipp = new(packets.IP)
		for i := 0; i < reg.IpsLength(); i++ {
			if !reg.Ips(ipp, i) {
				log.Warn().Msg("could not load ip")
				continue
			}

			ip := net.ParseIP(string(ipp.Ip()))
			if ip == nil {
				log.Warn().Msg("could not parse ip")
			}

			if broadcastIPnet.Contains(ip) {
				primaryIp, primaryMac = string(ipp.Ip()), string(ipp.Mac())
				log.Info().IPAddr("mainIp", ip).Msg("found the IP")
				break
			}
		}

		packets.PingStart(b)
		packets.PingAddHostname(b, hn)
		packets.PingAddIdentifier(b, guid)
		b.Finish(packets.PingEnd(b))

		var affected int64
		err = orm.Transaction(func(tx *gorm.DB) error {
			scoped := tx.Model(&Host{}).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "guid"}},
				DoUpdates: clause.AssignmentColumns([]string{"last_seen"}),
			}).Create(&Host{
				SimpleBaseModel: crud.SimpleBaseModel{
					Guid: crud.UUID(id),
				},
				Hostname:   string(banner.Hostname()),
				PrimaryIp:  primaryIp,
				PrimaryMac: primaryMac,
				Data:       msg.Data,
				LastSeen:   time.Now(),
			})

			affected = scoped.RowsAffected
			return scoped.Error
		})

		log.Err(err).Str("guid", id.String()).Int64("affected", affected).Msg("inserted host")
		if err == nil {
			err = nc.Publish(string(banner.Identifier())+"_welcome", b.FinishedBytes())
			log.Err(err).Str("guid", id.String()).Msg("responded")
		} else {
			log.Warn().Msg("skipped welcoming this host")
		}
	})

	log.Err(err).Msg("subscribed to register-a-new-host")

	settings = LoadSettings(orm)
	log.Err(err).Msg("loaded settings")
	if err != nil {
		log.Fatal().Msg("settings must load")
	}

	cc, err := settings.CrudController()
	log.Err(err).Msg("booted settings-controller")
	if err != nil {
		log.Fatal().Err(err).Msg("settings-controller required")
	}

	hostsController, err := crud.Controller(Host{}, crud.Specification{
		NotPaginated: true,
		Orm:          orm,
		ModelRoutes: []crud.RouteOperation{
			crud.DefaultOperation(crud.Single),
			crud.DefaultOperation(crud.List),
		},
		RelationRoutes: []crud.ModelRelationRoute{
			{
				Model: ExternalData{},
				Routes: []crud.RelationRouteOperation{
					crud.DefaultRelationOperation(crud.List),
					crud.DefaultRelationOperation(crud.Create),
				},
				NotPaginated: true,
			},
		},
	})
	log.Err(err).Msg("booted hosts-controller")
	if err != nil {
		log.Fatal().Err(err).Msg("hosts-controller required")
	}
	problemController, err := crud.Controller(Problem{}, crud.Specification{
		NotPaginated: true,
		Orm:          orm,
		ModelRoutes: []crud.RouteOperation{
			crud.DefaultOperation(crud.List),
		},
	})
	log.Err(err).Msg("booted problem-controller")
	if err != nil {
		log.Fatal().Err(err).Msg("hosts-controller required")
	}

	teamsController, err := crud.Controller(ExternalData{}, crud.Specification{
		NotPaginated: true,
		Orm:          orm,
		ModelRoutes: []crud.RouteOperation{
			crud.DefaultOperation(crud.Single),
			crud.DefaultOperation(crud.List),
			crud.DefaultOperation(crud.Partial),
		},
	})

	log.Err(err).Msg("booted teams-controller")
	if err != nil {
		log.Fatal().Err(err).Msg("settings-controller required")
	}

	router := rest.NewRouterWithSettings(crud.PanicHandler, func(r *router.Router) {
		r.RedirectTrailingSlash = true
		r.SaveMatchedRoutePath = true
		// r.GlobalOPTIONS = nil
		r.GlobalOPTIONS = func(ctx *fasthttp.RequestCtx) {

			headers := []string{
				"Access-Control-Request-Headers",
				"Access-Control-Request-Method",
			}

			for _, header := range headers {
				if h := ctx.Request.Header.Peek(header); h != nil {
					header = strings.Replace(header, "Request", "Allow", -1)
					header = strings.Replace(header, "Method", "Methods", -1)

					ctx.Response.Header.Set(header, string(h))
				}
			}

			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		}
		// r.HandleOPTIONS = true
	})

	router.Group("", func() {
		dashboardLocation := "fe/dist"
		dashboardFs := fasthttp.FS{
			Root:       dashboardLocation,
			IndexNames: []string{"index.html"},
			PathNotFound: func(ctx *fasthttp.RequestCtx) {
				// Attempt to load the index.html on 404
				f, err := os.Open(fmt.Sprintf("%v/index.html", dashboardLocation))
				if err != nil {
					ctx.SetStatusCode(500)
					log.Err(err).Msg("cannot find index file")
					return
				}

				defer f.Close()
				ctx.Response.Header.Set("content-type", "text/html")
				ctx.SetStatusCode(200)
				io.Copy(ctx, f)
			},
		}

		router.Get("{path:*}", dashboardFs.NewRequestHandler())

		router.Group("api", func() {

			cc.BuildRoutes(router)
			problemController.BuildRoutes(router)
			hostsController.BuildRoutes(router)
			teamsController.BuildRoutes(router)

			router.Get("contests", func(ctx *fasthttp.RequestCtx) {
				lg := rest.LoggerFromRequest(ctx)

				req, _ := http.NewRequest(http.MethodGet, url+"contests", nil)
				req.SetBasicAuth(ctx.UserValue("user").(string), ctx.UserValue("pass").(string))

				resp, err := http.DefaultClient.Do(req)
				lg.Err(err).Msg("retrieved contests at Dj")
				if err != nil {
					ctx.SetStatusCode(http.StatusInternalServerError)
					return
				}

				defer resp.Body.Close()
				io.Copy(ctx, resp.Body)
				ctx.SetStatusCode(resp.StatusCode)
			})

			router.Get("djTeam", djTeamLoad)
			router.Get("djProblem", djProblemLoad)

			router.Post("tim-json", func(ctx *fasthttp.RequestCtx) {
				lg := rest.LoggerFromRequest(ctx)

				mpf, err := ctx.MultipartForm()
				lg.Err(err).Msg("decoded upload as form")
				if err != nil {
					panic(crud.HttpError{
						Code: http.StatusNotAcceptable,
						Body: "not a form submit",
					})
				}

				files := mpf.File["files"]
				if fl := len(files); fl <= 0 {
					panic(crud.HttpError{
						Code: http.StatusNotAcceptable,
						Body: fmt.Sprintf("too few files being uploaded, %d being uploaded, while at least 0 are needed", fl),
					})
				}

				var hasErr bool
				for _, file := range files {
					var allowedFilePrefixes = map[string]func(zerolog.Logger, *fasthttp.RequestCtx, io.Reader) error{
						"map": func(lg zerolog.Logger, requestCtx *fasthttp.RequestCtx, f io.Reader) error {
							bts, err := io.ReadAll(f)
							lg.Err(err).Msg("copied all map-bytes")
							if err != nil {
								return err
							}

							settings.Set("map", bts)
							return nil
						},
						"problems": func(lg zerolog.Logger, requestCtx *fasthttp.RequestCtx, f io.Reader) error {
							djProblemLoad(ctx)

							var problems []Problem
							err := json.NewDecoder(f).Decode(&problems)
							lg.Err(err).Msg("decoded problems")
							if err != nil {
								return err
							}

							var affected int64
							err = orm.Transaction(func(tx *gorm.DB) error {
								scoped := tx.Model(&Problem{}).Clauses(clause.OnConflict{
									Columns:   []clause.Column{{Name: "id"}},
									DoUpdates: clause.AssignmentColumns([]string{"loc_x", "loc_y"}),
								}).Create(problems)

								affected = scoped.RowsAffected
								return scoped.Error
							})

							lg.Err(err).Int64("rowsaffected", affected).Msg("upserted")
							return err
						},
						"teams": func(log zerolog.Logger, requestCtx *fasthttp.RequestCtx, f io.Reader) error {
							// Load the teams from DOMjudge first
							djTeamLoad(ctx)

							var locations []struct {
								Tid string `json:"id" gorm:"team_id"`
								ExternalData
							}

							err := json.NewDecoder(f).Decode(&locations)
							log.Err(err).Msg("decoded teams")
							if err != nil {
								return err
							}

							var teams = make([]ExternalData, len(locations))
							for k := range locations {
								teams[k] = locations[k].ExternalData
								teams[k].TeamId = null.Filled(locations[k].Tid)
							}

							var affected int64
							err = orm.Transaction(func(tx *gorm.DB) error {
								scoped := tx.Model(&ExternalData{}).Clauses(clause.OnConflict{
									Columns:   []clause.Column{{Name: "team_id"}},
									DoUpdates: clause.AssignmentColumns([]string{"loc_x", "loc_y", "loc_rotation"}),
								}).Create(teams)

								affected = scoped.RowsAffected
								return scoped.Error
							})

							log.Err(err).Int64("rowsaffected", affected).Msg("upserted")
							return err
						},
					}

					for k, cb := range allowedFilePrefixes {
						f, err := file.Open()
						lg.Err(err).Msg("openend file")
						if err != nil {
							continue
						}

						if strings.HasPrefix(file.Filename, k) {
							err = cb(lg.With().Str("type", k).Logger(), ctx, f)
							hasErr = hasErr || err != nil
						}

						f.Close()
					}
				}

				if hasErr {
					panic(crud.HttpError{
						Code: http.StatusFailedDependency,
						Body: "some error occured",
					})
				}

			})
		})
	}, basicAuth)

	_, err = rest.New(rest.Config{
		Listen: ":4000",
		// Customize: func(server *fasthttp.Server) error {
		// },
		// TLS: &rest.TLS{
		// 	PrivateKey:  "/home/mart/eventix/goproxy/ssl/eventix.key.pem",
		// 	Certificate: "/home/mart/eventix/goproxy/ssl/eventix.cert.pem",
		// },
	}, router)
	log.Err(err).Msg("started rest")
	if err != nil {
		log.Fatal().Msg("bye")
	}

	log.Info().Msg("Booted")
	lifecycle.Finally(func() { log.Warn().Msg("Stopping") })
	lifecycle.StopListener()
}

var numRequests sync.Map
var checkEvery = env.Uint64("CHECK_PASSWORD_EVERY", 10)

func basicAuth(next rest.Handle) rest.Handle {
	var basicAuthPrefix = []byte("Basic ")

	return func(ctx *fasthttp.RequestCtx) { // func(w http.ResponseWriter, r *http.Request) {
		auth := ctx.Request.Header.Peek("Authorization")
		if env.Bool("IS_DEV", false) {
			auth = []byte("Basic YWRtaW46YWRtaW4=")
		}

		if bytes.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
			if err == nil {
				pair := strings.SplitN(string(payload), ":", 2)
				if len(pair) >= 2 {
					user, pass := pair[0], pair[1]

					// TODO optimize this code, all these conversions are unneeded
					combined := user + ":" + pass
					var hasher = sha1.New()
					hasher.Write([]byte(combined))

					combined = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
					numReqs, exists := numRequests.LoadOrStore(combined, uint64(0))

					var notOk bool

					// Check if we have it cached
					if !exists || (numReqs.(uint64))%checkEvery == 0 {
						// Do a check in Dj to see if valid
						req, _ := http.NewRequest(http.MethodGet, url+"users", nil)
						req.SetBasicAuth(user, pass)

						resp, err := http.DefaultClient.Do(req)
						log.Err(err).Msg("checked credentials at Dj")
						defer resp.Body.Close()

						if err == nil && resp.StatusCode == 200 {
							notOk = false
						}
					}

					if !notOk {
						ctx.SetUserValue("user", user)
						ctx.SetUserValue("pass", pass)
						numRequests.Store(combined, numReqs.(uint64)+1)
						next(ctx)

						return
					}
				}
			}
		}

		// Request Basic Authentication otherwise
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted, charset=\"UTF-8\"")
		ctx.SetStatusCode(http.StatusUnauthorized)
		// ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
	}
	// }
	//
	// w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	// http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func mdnsServe(natsAddr string) {
	// Setup our service export
	host, _ := os.Hostname()

	service, err := mdns.NewMDNSService(host, "_pixie._tcp", "progcont.", "", 8000, nil, []string{natsAddr})
	log.Err(err).Msg("created mDNS service")

	// Create the mDNS server, defer shutdown
	srv, err := mdns.NewServer(&mdns.Config{Zone: service})

	log.Err(err).Msg("started mDNS advertisement")
	lifecycle.EFinally(srv.Shutdown)
}

// GetInternalIp return internal ipv4
func GetInternalIps() []*net.IPNet {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		panic(err.Error())
	}
	ips := make([]*net.IPNet, 0, len(addr))
	for _, a := range addr {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ips = append(ips, ipNet)
		}
	}

	return ips
}

func djTeamLoad(ctx *fasthttp.RequestCtx) {
	lg := rest.LoggerFromRequest(ctx)

	req, _ := http.NewRequest(http.MethodGet, url+"users", nil)
	req.SetBasicAuth(ctx.UserValue("user").(string), ctx.UserValue("pass").(string))
	req.Header.Add("Accept-Charset", "utf-8")

	resp, err := http.DefaultClient.Do(req)
	log.Err(err).Msg("loaded from Dj")
	crud.HandleError(ctx, http.StatusInternalServerError, fmt.Errorf("loading from DOMjudge; %w", err))

	ctx.SetStatusCode(resp.StatusCode)
	lg.Debug().Int("status", resp.StatusCode).Msg("Dj response")
	defer resp.Body.Close()

	var m []ExternalData
	err = json.NewDecoder(resp.Body).Decode(&m)
	lg.Err(err).Msg("decoded Dj response")

	var affected int64
	err = orm.Transaction(func(tx *gorm.DB) error {
		scoped := tx.Model(&ExternalData{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"teamname", "username", "user_id"}),
		}).Create(m)

		affected = scoped.RowsAffected
		return scoped.Error
	})

	lg.Err(err).Int64("rowsaffected", affected).Msg("inserted")
	crud.HandleError(ctx, http.StatusInternalServerError, fmt.Errorf("inserting into database; %w", err))
}

func djProblemLoad(ctx *fasthttp.RequestCtx) {
	lg := rest.LoggerFromRequest(ctx)

	req, _ := http.NewRequest(http.MethodGet, url+"contests/"+settings.Retrieve("contest")+"/problems", nil)
	req.SetBasicAuth(ctx.UserValue("user").(string), ctx.UserValue("pass").(string))
	req.Header.Add("Accept-Charset", "utf-8")

	resp, err := http.DefaultClient.Do(req)
	log.Err(err).Msg("loaded from Dj")
	crud.HandleError(ctx, http.StatusInternalServerError, fmt.Errorf("loading from DOMjudge; %w", err))

	ctx.SetStatusCode(resp.StatusCode)
	lg.Debug().Int("status", resp.StatusCode).Msg("Dj response")
	defer resp.Body.Close()

	type djProblem struct {
		Id string `json:"short_name"`
		Problem
	}

	var m []djProblem
	err = json.NewDecoder(resp.Body).Decode(&m)
	lg.Err(err).Msg("decoded Dj response")

	var probs = make([]Problem, len(m))
	for k, v := range m {
		probs[k] = v.Problem
		probs[k].Id = v.Id
	}

	var affected int64
	err = orm.Transaction(func(tx *gorm.DB) error {
		scoped := tx.Model(&Problem{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"rgb"}),
		}).Create(probs)

		affected = scoped.RowsAffected
		return scoped.Error
	})

	lg.Err(err).Int64("rowsaffected", affected).Msg("inserted problems")
	crud.HandleError(ctx, http.StatusInternalServerError, fmt.Errorf("inserting problems into database; %w", err))
	crud.Respond(ctx, probs)
}
