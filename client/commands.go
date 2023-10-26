package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/kevinburke/go.uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"

	"github.com/tuupke/pixie/packets"
)

var (
	// TODO load from env
	configFileLocation = "/etc/lightdm/lightdm-qt5-greeter.conf"
	user               = "mart"

	conn *dbus.Conn
)

func notify(header, body string) {
	err := exec.Command("sudo", "-u", user, "DISPLAY=:0", "DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/"+"1000"+"/bus", "notify-send", header, body).Run()
	if err == nil {
		return
	}

	fmt.Printf("notify-err ('%v', '%v'), '%v' %v, err: %v\n", user, 1000, header, body, err)
}

// reboot forces an unexpected reboot
func reboot(in time.Duration) {
	if in >= time.Second*3 {
		notify("Logging out", fmt.Sprintf("Machine will log out in %s", in))
	}

	// Sleep for the required amount of time, sleep(n), n <= 0, returns immediately
	time.Sleep(in)

	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	log.Err(err).Msg("sent reboot")
}

func off(in time.Duration) {
	if in >= time.Second*3 {
		notify("Shutting down", fmt.Sprintf("Machine will shut down in %s", in))
	}

	// Sleep for the required amount of time, sleep(n), n <= 0, returns immediately
	time.Sleep(in)

	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	log.Err(err).Msg("sent power off")
}

// backup forces a manual backup
func backup(backupLocation string) {
	// TODO
}

// logout can be used to forcefully log-out the user
func logout(in time.Duration) {
	if in >= time.Second*3 {
		notify("Logging out", fmt.Sprintf("Machine will log out in %s", in))
	}

	// Sleep for the required amount of time, sleep(n), n <= 0, returns immediately
	time.Sleep(in)

	resChan := make(chan string)
	defer close(resChan)
	// Restart the window manager, TODO fill the mode
	res, err := conn.RestartUnitContext(context.Background(), "lightdm-greeter", "", resChan)
	log.Err(err).Int("jobId", res).Msg("issued restart")
	if err != nil {
		reportResult("", err)
		return
	}

	reportResult(<-resChan, nil)
}

func reportResult(res string, err error) {
	// TODO send this to beanstalk
}

type greeterOptions struct{ username, password, contestId, background, apiUrl, chainString string }

// setGreeter can be used to override the credentials in the greeter
func setGreeter(username, password, contestId, background, apiUrl string, chain []string, reload bool) {
	for k := range chain {
		chain[k] = strings.TrimSpace(chain[k])
	}

	chainString := strings.Join(chain, ",")
	if chainString == "" {
		chainString = "Up,Up,Down,Down,A,B"
	}

	f, err := os.OpenFile(configFileLocation, os.O_TRUNC|os.O_WRONLY, 0755)
	log.Err(err).Msg("opening config file")
	if err != nil {
		reportResult("", err)
		return
	}

	greeterOptions{
		username,
		password,
		contestId,
		background,
		apiUrl,
		chainString,
	}.WriteGenerate(f)

	defer deferLog(zerolog.DefaultContextLogger, f, "closing file")
	if reload {
		logout(time.Second)
	}

	reportResult("OK", nil)
}

func deferLog(log *zerolog.Logger, closer io.Closer, msg string) {
	log.Err(closer.Close()).Msg("string")
}

var banner = struct {
	Identifier uuid.UUID `json:"identifier"`
	Hostname   string    `json:"hostname"`
}{
	uuid.NewV4(),
	"",
}

func buildBanner(b *flatbuffers.Builder) flatbuffers.UOffsetT {

	hn := b.CreateSharedString(banner.Hostname)
	identifier := b.CreateByteVector(banner.Identifier.Bytes())

	packets.BannerStart(b)
	packets.BannerAddHostname(b, hn)
	packets.BannerAddIdentifier(b, identifier)

	return packets.BannerEnd(b)
}

func buildRegister(log zerolog.Logger, b *flatbuffers.Builder) flatbuffers.UOffsetT {
	banner := buildBanner(b)

	addrs, err := net.InterfaceAddrs()
	log.Err(err).Msg("retrieved IPs")

	var ipSet = make([]flatbuffers.UOffsetT, 0, len(addrs))
	for _, addr := range addrs {
		ipSet = append(ipSet, b.CreateSharedString(addr.String()))
	}

	packets.RegisterStartIpsVector(b, len(ipSet))
	for _, ip := range ipSet {
		b.PrependUOffsetT(ip)
	}

	ips := b.EndVector(len(ipSet))

	packets.RegisterStart(b)
	packets.RegisterAddBanner(b, banner)
	packets.RegisterAddIps(b, ips)

	return packets.RegisterEnd(b)
}

func renderQr(data string) {
	png, err := qrcode.Encode(data, qrcode.Highest, 600)
	log.Err(err).Msg("encoded qr")

	cmd := exec.Command("imv-x11", "-", "-f")
	cmd.Env = append(cmd.Env,
		"DISPLAY=:0",
		"XAUTHORITY=/var/lib/lightdm/.Xauthority")
	cmd.Stdin = bytes.NewReader(png)
	outp, err := cmd.CombinedOutput()
	fmt.Println(err, string(outp))
}
