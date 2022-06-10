package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/hashicorp/mdns"
	"github.com/rs/zerolog/log"
	"go.dedis.ch/kyber/v3/pairing/bn256"
	"go.dedis.ch/kyber/v3/sign/bls"
	"go.dedis.ch/kyber/v3/util/random"

	"github.com/tuupke/pixie/beanstalk"
	"github.com/tuupke/pixie/packets"
)

var (
	reserveDuration = time.Second * 5

	hostname, _ = os.Hostname()
)

func main() {

	suite := bn256.NewSuite()
	private, public := bls.NewKeyPair(suite, random.New())

	fmt.Println(private.MarshalBinary())
	fmt.Println(public.MarshalBinary())

	n := time.Now()
	sign, err := bls.Sign(suite, private, []byte("aaa"))
	sin := time.Since(n)
	fmt.Println(sign, len(sign), err, sin)

	n = time.Now()
	for i := 0; i < 1000; i++ {
		err = bls.Verify(suite, public, []byte("aaa"), sign)
	}
	sin = time.Since(n)
	fmt.Println(err, sin/1000)

	_, err = json.Marshal(banner)
	if err != nil {
		log.Fatal().Err(err).Msg("could not marshal banner")
	}

	bsAddr, err := attemptDiscovery()
	log.Err(err).Str("bsAddr", bsAddr).Msg("discovery Result")
	requestTubes := []string{"allClients", hostname}

	bConn := beanstalk.Connect("tcp", bsAddr, requestTubes...)
	tubes, err := bConn.ListTubes()
	log.Info().Err(err).Strs("requested", requestTubes).Strs("tubes", tubes).Msg("connected to beanstalk")

	if err != nil {
		log.Fatal().Msg("requires beanstalk connection")
	}

	queueListen(bConn)
}

func attemptDiscovery() (bsAddr string, err error) {
	serviceName := "_beanstalk._pixie._tcp"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		entriesCh := make(chan *mdns.ServiceEntry)
		p := mdns.DefaultParams(serviceName)
		p.Domain = "progcont."
		p.Timeout = time.Minute / 12
		p.Entries = entriesCh
		go func() {
			err = mdns.Query(p)
			log.Err(err).Msg("queried for mDNS")
			close(p.Entries)
		}()

		defer wg.Done()

		for entry := range entriesCh {
			if strings.Contains(entry.Name, serviceName) {
				bsAddr = entry.Info
				return
			} else {
				log.Debug().Str("name", entry.Name).Msg("received other mDNS")
			}
		}
	}()

	wg.Wait()

	if err == nil && bsAddr == "" {
		err = errors.New("could not connect to beanstalk")
	}

	return
}

func queueListen(conn beanstalk.Connection) {
	for {
		id, body, err := conn.Reserve(reserveDuration)
		if err != nil {
			if !beanstalk.TimeoutErr(err) {
				log.Err(err).Msg("subscribing")
			}

			continue
		}

		var ee interface{}
		// Attempt to send over the chan, attempt to recover a failed chan send. Usually
		// due to a closed channel
		func(b []byte) {
			defer func() {
				ee = recover()
			}()

			unionTable := new(flatbuffers.Table)
			request := packets.GetRootAsCommand(b, 0)
			if request.Command(unionTable) {
				log.Info().Stringer("command", request.CommandType()).Msg("received command")
				switch typ := request.CommandType(); typ {
				case packets.CmdNONE:
					// no-op
				case packets.CmdReboot:
					var r = new(packets.Logout)
					r.Init(unionTable.Bytes, unionTable.Pos)
					reboot(time.Duration(r.In()))
				case packets.CmdAnsible:
					// todo
				case packets.CmdLogout:
					var l = new(packets.Logout)
					l.Init(unionTable.Bytes, unionTable.Pos)

					logout(time.Duration(l.In()))
				case packets.CmdNotify:
					var n = new(packets.Notify)
					n.Init(unionTable.Bytes, unionTable.Pos)

					notify(string(n.Header()), string(n.Body()))
				}
			}
		}(body)

		log.Err(err).Msg("handled message")
		if ee != nil {
			log.Err(conn.Release(id, 0, time.Second)).Msg("released")
		} else {
			log.Err(conn.Delete(id)).Msg("deleted")
		}

	}
}
