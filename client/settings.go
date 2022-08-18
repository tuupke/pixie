package main

import (
	"net"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/rs/zerolog/log"
)

func initializeSettings() settings {
	hn, _ := os.Hostname()
	intfs, _ := net.Interfaces()

	s := settings{
		hn:         hn,
		identifier: uuid.NewV1(),
		connection: connection{
			connected:   binding.NewBool(),
			multicast:   binding.NewString(),
			environment: os.Getenv("PIXIE_ADDR"),
			connecting:  binding.NewString(),
			registered:  binding.NewBool(),
		},
		networks: make([]network, 0, len(intfs)),
	}

	// mp := binding.NewUntypedMap()

	s.connection.multicast.Set("undiscovered")

	if s.connection.environment == "" {
		s.connection.environment = "not set"
	}

	for _, intf := range intfs {
		if addrs, err := intf.Addrs(); err == nil {
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ipSplit := strings.SplitN(addr.String(), "/", 2)
					if len(ipSplit) != 2 {
						continue
					}

					netmask, err := strconv.Atoi(ipSplit[1])
					log.Err(err).Str("netmaskbytes", ipSplit[1]).Msg("parsed netmask bytes to int")
					if err != nil {
						continue
					}

					s.networks = append(s.networks, network{
						hwAddr:       intf.HardwareAddr.String(),
						name:         intf.Name,
						netmaskBytes: netmask,
						ip:           ipSplit[0],
					})
				}
			}
		}
	}

	return s
}

type network struct {
	hwAddr       string
	name         string
	netmaskBytes int
	ip           string
}

type connection struct {
	environment string
	connected   binding.Bool
	multicast   binding.String
	connecting  binding.String
	registered  binding.Bool
}

type settings struct {
	hn string

	connection connection

	networks   []network
	identifier uuid.UUID
}
