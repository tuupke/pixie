package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/hashicorp/mdns"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/tuupke/pixie/packets"

	_ "openticket.tech/log/v2"
)

func main() {
	s := initializeSettings()

	go func() {
		bsAddr, err := attemptDiscovery()
		log.Err(err).Str("bsAddr", bsAddr).Msg("discovery Result")

		if bsAddr != "" {
			s.connection.multicast.Set(bsAddr)
		}
	}()

	go func() {
		var connected bool
		for !connected {
			ips := make([]string, 0, 2)
			if ip := s.connection.environment; ip != "not set" {
				ips = append(ips, ip)
			}

			if ip, err := s.connection.multicast.Get(); err == nil && ip != "undiscovered" {
				ips = append(ips, ip)
			}

			for _, ip := range ips {
				time.Sleep(time.Second)
				host, port, err := net.SplitHostPort(ip)
				log.Err(err).Str("ip", ip).Str("host", host).Str("port", port).Msg("split ip")
				if err != nil {
					continue
				}

				s.connection.connecting.Set(ip)

				nc, err := nats.Connect(ip)
				log.Err(err).Msg("connected to nats")
				if err == nil {
					s.connection.connected.Set(true)
					connected = true

					var wg sync.WaitGroup
					wg.Add(1)
					_, err := nc.Subscribe(s.identifier.String()+"_welcome", func(msg *nats.Msg) {
						ping := packets.GetRootAsPing(msg.Data, 0)

						fmt.Println(string(ping.Identifier()))

						msg.Sub.Unsubscribe()
						wg.Done()
					})

					log.Err(err).Str("subject", s.identifier.String()).Msg("subscribed")

					b := flatbuffers.NewBuilder(256)
					fmt.Println("pre register")
					b.Finish(s.register(b))
					fmt.Println("post register")

					bts := b.FinishedBytes()

					nc.Publish("register-a-new-host", bts)
					log.Debug().Bytes("registration", bts).Msg("sending registration")
					log.Err(err).Str("subject", "register-a-new-host").Msg("published")

					wg.Wait()
					log.Info().Msg("received reply")

					fyne.CurrentApp().Settings().SetTheme(newCustomTheme())

					break
				}
				s.connection.connecting.Set("Could not connect to " + ip)
				time.Sleep(time.Second)
			}
		}
	}()

	s.start()

	return
}

func attemptDiscovery() (bsAddr string, err error) {
	serviceName := "_pixie._tcp"

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {

		p := mdns.DefaultParams(serviceName)
		p.Domain = "progcont"
		p.Timeout = time.Second * 3
		p.WantUnicastResponse = false
		// p.DisableIPv6 = true
		// p.Service = serviceName
		defer wg.Done()

		for {
			entriesCh := make(chan *mdns.ServiceEntry, 10)
			p.Entries = entriesCh

			go func() {
				err = mdns.Query(p)
				log.Err(err).Msg("queried for mDNS")
				time.Sleep(time.Millisecond * 10)
				close(entriesCh)
			}()

			for entry := range entriesCh {
				if strings.Contains(entry.Name, serviceName) {
					bsAddr = entry.Info

					return
				}

				log.Debug().Str("name", entry.Name).Msg("received other mDNS")
			}

			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	if err == nil && bsAddr == "" {
		err = errors.New("could not connect to beanstalk")
	}

	return
}
