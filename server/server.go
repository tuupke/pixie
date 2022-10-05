package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fasthttp/router"
	"github.com/go-gormigrate/gormigrate/v2"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"github.com/hashicorp/mdns"
	"github.com/nats-io/nats-server/v2/server"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/tuupke/pixie"
	"github.com/tuupke/pixie/crud"
	"github.com/tuupke/pixie/lifecycle"
	"github.com/tuupke/pixie/packets"
)

var listenAddr = envStringFb("LISTEN_ADDR", ":4000")

type wString string

func (w wString) String() string { return string(w) }

type Problem struct {
	Id       string   `gorm:"primaryKey" json:"id"`
	Rgb      *string  `json:"rgb"`
	Location Location `gorm:"embedded;embeddedPrefix:loc_" json:"location"`
}

func (p Problem) Identifier() fmt.Stringer {
	return wString(p.Id)
}

type ExternalData struct {
	Guid     crud.UUID `gorm:"primaryKey" json:"guid"`
	Username string    `json:"username"`
	UserId   string    `json:"id"`
	Teamname *string   `json:"team"`
	TeamId   *string   `gorm:"uniqueIndex" json:"team_id"`
	HostId   *string   `gorm:"index:" json:"host_id"`
	Location Rotated   `gorm:"embedded;embeddedPrefix:loc_" json:"location"`
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
	Guid       crud.UUID `gorm:"primaryKey,uniqueIndex" json:"guid"`
	Hostname   string    `json:"hostname"`
	PrimaryIp  string    `json:"primary_ip"`
	PrimaryMac string    `json:"primary_mac"`
	Data       []byte    `json:"data"`
	LastSeen   time.Time `json:"last_seen"`
}

var natsHost, natsPort string
var natsPortInt int
var orm *gorm.DB
var broadcastIPnet *net.IPNet

func djUrl() string {
	return strings.TrimSuffix(settings.Retrieve("domjudge"), "/") + "/"
}

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

var settings *pixie.Settings

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

	orm = pixie.Orm()
	log.Err(gormigrate.New(orm, gormigrate.DefaultOptions, migrations()).Migrate()).Msg("migrated")

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

		// Find team
		var team = new(ExternalData)
		err = orm.Model(team).Where(crud.PrimaryKeyExpression(banner.Identifier())).First(&team).Error
		log.Err(err).Msg("retrieved team")

		packets.PingStart(b)
		packets.PingAddHostname(b, hn)
		packets.PingAddIdentifier(b, guid)
		bannerOffset := packets.PingEnd(b)

		var hasTeam bool
		var teamId, teamName flatbuffers.UOffsetT
		if err != nil && team.TeamId != nil && team.Teamname != nil {
			teamId = b.CreateSharedString(*team.TeamId)
			teamName = b.CreateSharedString(*team.Teamname)
			hasTeam = true
		}

		packets.WelcomeStart(b)
		packets.WelcomeAddBanner(b, bannerOffset)
		packets.WelcomeAddTeamId(b, teamId)
		packets.WelcomeAddTeamName(b, teamName)
		packets.WelcomeAddHasTeam(b, hasTeam)

		b.Finish(packets.WelcomeEnd(b))

		var affected int64
		err = orm.Transaction(func(tx *gorm.DB) error {
			scoped := tx.Model(&Host{}).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "guid"}},
				DoUpdates: clause.AssignmentColumns([]string{"last_seen"}),
			}).Create(&Host{
				Guid:       crud.UUID(id),
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
			err = nc.Publish(string(banner.Identifier())+".welcome", b.FinishedBytes())
			log.Err(err).Str("guid", id.String()).Msg("responded")
		} else {
			log.Warn().Msg("skipped welcoming this host")
		}
	})

	log.Err(err).Msg("subscribed to register-a-new-host")

	settings = pixie.LoadSettings(orm)
	log.Err(err).Msg("loaded settings")
	if err != nil {
		log.Fatal().Msg("settings must load")
	}

	edc := crud.New[ExternalData](orm)

	rtr := router.New()
	api := rtr.Group("/api/")

	ed := api.Group("external_data/")
	ed.GET("", edc.List)
	ed.GET("{guid}/", edc.Get)
	ed.PATCH("{guid}/", edc.Partial)

	pbc := crud.New[Problem](orm)
	pb := api.Group("problem/")
	pb.GET("", pbc.List)

	stc := settings.CrudController()
	st := api.Group("setting/")
	st.GET("", stc.List)
	st.GET("{guid}/", stc.Get)
	st.PATCH("{guid}/", stc.Partial)

	hoc := crud.New[Host](orm)
	ho := api.Group("host/")
	ho.GET("", hoc.List)
	ho.GET("{guid}/", hoc.Get)
	ho.POST("{guid}/window/", func(ctx *fasthttp.RequestCtx) {
		guid := ctx.UserValue("guid").(string)
		lg := crud.LoggerFromRequest(ctx)
		lg.Err(nc.Publish(guid+"."+pixie.Show, nil)).Str("guid", guid).Msg("sent show")
	})
	ho.DELETE("{guid}/window/", func(ctx *fasthttp.RequestCtx) {
		guid := ctx.UserValue("guid").(string)
		lg := crud.LoggerFromRequest(ctx)
		lg.Err(nc.Publish(guid+"."+pixie.Hide, nil)).Str("guid", guid).Msg("sent hide")
	})

	var pathHandler fasthttp.RequestHandler
	if !envBoolFb("IS_DEV", false) {
		dashboardLocation := "fe/dist/"
		pathHandler = (&fasthttp.FS{
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
		}).NewRequestHandler()
	} else {
		u, err := url.Parse("http://127.0.0.1:5173")
		log.Err(err).Msg("parsed reverse proxy url")
		if err != nil {
			log.Fatal().Msg("need reverse proxy in dev mode")
		}
		frontendProxy := httputil.NewSingleHostReverseProxy(u)
		pathHandler = fasthttpadaptor.NewFastHTTPHandler(frontendProxy)
	}

	rtr.GET("/{path:*}", pathHandler)

	api.GET("inventory", func(ctx *fasthttp.RequestCtx) {
		var ansible = bytes.NewBufferString(`clients:
  vars:
    ansible_user: root
  hosts:`)

		var all []struct {
			IpAddress string `gorm:"column:primary_ip"`
			TeamName  string `gorm:"column:teamname"`
			Host      string `gorm:"column:user_id"`
		}

		orm.Table("hosts").
			Joins("join external_data  on external_data.host_id = hosts.guid").
			Select("hosts.primary_ip", "external_data.user_id", "external_data.teamname").
			Find(&all)

		for _, e := range all {
			ansible.WriteString(fmt.Sprintf(`
    %v:
      ansible_host: %v
      team_name_dj: "%v"`, e.Host, e.IpAddress, e.TeamName))
		}

		ctx.WriteString(ansible.String())
	})

	api.GET("contests", func(ctx *fasthttp.RequestCtx) {
		lg := crud.LoggerFromRequest(ctx)

		req, _ := http.NewRequest(http.MethodGet, djUrl()+"contests", nil)
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

	api.GET("djTeam", djTeamLoad)
	api.GET("djProblem", djProblemLoad)

	api.POST("tim-json", func(ctx *fasthttp.RequestCtx) {
		lg := crud.LoggerFromRequest(ctx)

		mpf, err := ctx.MultipartForm()
		lg.Err(err).Msg("decoded upload as form")
		crud.HandleError(ctx, http.StatusNotAcceptable, err)

		files := mpf.File["files"]
		if fl := len(files); fl <= 0 {
			crud.HandleError(ctx, http.StatusNotAcceptable, fmt.Errorf("too few files being uploaded, %d being uploaded, while at least 0 are needed", fl))
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
						Guid     crud.UUID `gorm:"primaryKey" json:"guid"`
						Tid      string    `json:"id" gorm:"column:team_id"`
						Location Rotated   `json:"location" gorm:"embedded;embeddedPrefix:loc_"`
					}

					err := json.NewDecoder(f).Decode(&locations)
					log.Err(err).Msg("decoded teams")
					if err != nil {
						return err
					}

					scoped := orm.Table("external_data").Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "team_id"}},
						DoUpdates: clause.AssignmentColumns([]string{"loc_x", "loc_y", "loc_rotation"}),
					}).Create(locations)

					err = scoped.Error
					rowsAffected := scoped.RowsAffected

					log.Err(err).Int64("affected", rowsAffected).Msg("upserted")
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
			crud.HandleError(ctx, http.StatusFailedDependency, errors.New("some error occured"))
		}

	})

	rtr.RedirectTrailingSlash = true
	rtr.SaveMatchedRoutePath = true
	// r.GlobalOPTIONS = nil
	rtr.GlobalOPTIONS = func(ctx *fasthttp.RequestCtx) {

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

	rtr.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
		lg := crud.LoggerFromRequest(ctx)
		lg.Error().Interface("error", i).Msg("panic received")
	}

	rtr.NotFound = func(ctx *fasthttp.RequestCtx) {
		lg := crud.LoggerFromRequest(ctx)
		lg.Warn().Msg("not found")
	}

	rtr.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		lg := crud.LoggerFromRequest(ctx)
		lg.Warn().Msg("not allowed")
	}

	err = fasthttp.ListenAndServe(listenAddr, basicAuth(rtr.Handler))

	log.Err(err).Msg("started rest")
	if err != nil {
		log.Fatal().Msg("bye")
	}

	privateRtr := router.New()
	privateRtr.ANY("/{path:*}", CupsHandler)

	err = fasthttp.ListenAndServe(cupsListen, privateRtr.Handler)
	log.Err(err).Msg("started cups proxy")

	log.Info().Msg("Booted")
	lifecycle.Finally(func() { log.Warn().Msg("Stopping") })
	lifecycle.StopListener()
}

var numRequests sync.Map
var checkEvery = envUInt64Fb("CHECK_PASSWORD_EVERY", 10)

func basicAuth(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	var basicAuthPrefix = []byte("Basic ")

	return func(ctx *fasthttp.RequestCtx) { // func(w http.ResponseWriter, r *http.Request) {
		auth := ctx.Request.Header.Peek("Authorization")
		if envBoolFb("IS_DEV", false) {
			auth = []byte("Basic YWRtaW46YWRtaW4=")
		}

		if bytes.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(string(auth[len(basicAuthPrefix):]))
			log.Err(err).Msg("decoded")
			if err == nil {
				pair := strings.SplitN(string(payload), ":", 2)
				if len(pair) >= 2 {
					user, pass := pair[0], pair[1]
					log.Debug().Str("user", user).Str("pass", pass).Str("dj_url", djUrl()).Msg("auth")

					// TODO optimize this code, all these conversions are unneeded
					combined := user + ":" + pass
					var hasher = sha1.New()
					hasher.Write([]byte(combined))

					combined = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
					numReqs, exists := numRequests.LoadOrStore(combined, uint64(0))

					var notOk bool

					// Check if we have it cached
					if !exists || (numReqs.(uint64))%checkEvery == 0 {
						notOk = true
						// Do a check in Dj to see if valid
						req, _ := http.NewRequest(http.MethodGet, djUrl()+"users", nil)
						req.SetBasicAuth(user, pass)

						resp, err := http.DefaultClient.Do(req)
						log.Err(err).Msg("checked credentials at Dj")

						if err == nil {
							defer resp.Body.Close()
						}

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
		} else {
			log.Debug().Msg("no basic auth")
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
	lg := crud.LoggerFromRequest(ctx)

	req, _ := http.NewRequest(http.MethodGet, djUrl()+"users", nil)
	req.SetBasicAuth(ctx.UserValue("user").(string), ctx.UserValue("pass").(string))
	req.Header.Add("Accept-Charset", "utf-8")

	resp, err := http.DefaultClient.Do(req)
	log.Err(err).Msg("loaded from Dj")
	crud.HandleError(ctx, http.StatusInternalServerError, err)

	ctx.SetStatusCode(resp.StatusCode)
	log.Debug().Int("status", resp.StatusCode).Msg("Dj response")
	defer resp.Body.Close()

	var m []ExternalData
	err = json.NewDecoder(resp.Body).Decode(&m)
	log.Err(err).Msg("decoded Dj response")

	// Users without a team are useless to us, this also eliminates all non-team users
	var mm = make([]ExternalData, 0, len(m))
	for _, v := range m {
		if v.TeamId != nil {
			mm = append(mm, v)
		}
	}
	m = mm

	var affected int64
	err = orm.Transaction(func(tx *gorm.DB) error {
		scoped := tx.Model(&ExternalData{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"username", "user_id"}),
		}, clause.OnConflict{
			Columns:   []clause.Column{{Name: "team_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"teamname"}),
		}).Create(m)

		affected = scoped.RowsAffected
		return scoped.Error
	})

	lg.Err(err).Int64("rowsaffected", affected).Msg("inserted dj")
	crud.HandleError(ctx, http.StatusInternalServerError, err)
}

func djProblemLoad(ctx *fasthttp.RequestCtx) {
	lg := crud.LoggerFromRequest(ctx)

	req, _ := http.NewRequest(http.MethodGet, djUrl()+"contests/"+settings.Retrieve("contest")+"/problems", nil)
	req.SetBasicAuth(ctx.UserValue("user").(string), ctx.UserValue("pass").(string))
	req.Header.Add("Accept-Charset", "utf-8")

	resp, err := http.DefaultClient.Do(req)
	lg.Err(err).Msg("loaded from Dj")
	crud.HandleError(ctx, http.StatusInternalServerError, err)

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
	crud.HandleError(ctx, http.StatusInternalServerError, err)
	crud.Respond(ctx, probs)
}
