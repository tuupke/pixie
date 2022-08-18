package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/rs/zerolog/log"
	"github.com/skip2/go-qrcode"

	"github.com/tuupke/pixie/packets"
)

func (s *settings) banner(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	id := b.CreateString(s.identifier.String())
	hn := b.CreateString(s.hn)

	packets.BannerStart(b)
	packets.BannerAddIdentifier(b, id)
	packets.BannerAddHostname(b, hn)
	return packets.BannerEnd(b)
}

func (s *settings) register(b *flatbuffers.Builder) flatbuffers.UOffsetT {
	banner := s.banner(b)

	var ipsVector flatbuffers.UOffsetT
	if s.networks != nil {
		netLen := len(s.networks)
		ips := make([]flatbuffers.UOffsetT, netLen)
		for k := 0; k < netLen; k++ {
			fmt.Println("Added ip", s.networks[k].ip)
			ip := b.CreateSharedString(s.networks[k].ip)
			mac := b.CreateSharedString(s.networks[k].hwAddr)
			name := b.CreateString(s.networks[k].name)

			packets.IPStart(b)
			packets.IPAddIp(b, ip)
			packets.IPAddNetmask(b, int32(s.networks[k].netmaskBytes))
			packets.IPAddName(b, name)
			packets.IPAddMac(b, mac)

			ips[k] = packets.IPEnd(b)
		}

		packets.RegisterStartIpsVector(b, netLen)
		for _, v := range ips {
			b.PrependUOffsetT(v)
		}

		ipsVector = b.EndVector(netLen)
	}

	packets.RegisterStart(b)
	packets.RegisterAddBanner(b, banner)
	packets.RegisterAddIps(b, ipsVector)

	return packets.RegisterEnd(b)
}

func (s *settings) start() {
	a := app.New()
	w := a.NewWindow("Pixie registration")

	buf := bytes.NewBuffer(nil)

	b := flatbuffers.NewBuilder(128)
	b.Finish(s.banner(b))

	base64.NewEncoder(base64.StdEncoding, buf).Write(b.FinishedBytes())

	png, err := qrcode.New(buf.String(), qrcode.Highest)
	log.Err(err).Msg("qr rendered")

	qr := canvas.NewImageFromImage(png.Image(600))
	qr.SetMinSize(fyne.NewSize(600, 600))
	qr.FillMode = canvas.ImageFillContain

	var nets = make([]*widget.AccordionItem, 0, len(s.networks))
	for _, n := range s.networks {
		ips := make([]fyne.CanvasObject, 0, 2)

		if n.hwAddr != "" {
			ips = append(ips, widget.NewLabel(n.hwAddr))
		}

		ips = append(ips, widget.NewLabel(n.ip))
		nets = append(nets, widget.NewAccordionItem(n.name, container.NewVBox(ips...)))
	}

	connCheck := widget.NewCheckWithData("", s.connection.connected)
	connCheck.Disable()

	regCheck := widget.NewCheckWithData("", s.connection.registered)
	regCheck.Disable()

	w.SetContent(container.NewHSplit(
		container.NewVScroll(widget.NewAccordion(nets...)),
		container.NewVBox(
			widget.NewLabel("Welcome to pixie, scan the QR to register host"),
			qr,
			container.NewHBox(
				container.NewVBox(
					widget.NewForm(
						widget.NewFormItem("Hostname", widget.NewLabel(s.hn)),
						widget.NewFormItem("Identifier", widget.NewLabel(s.identifier.String())),
						widget.NewFormItem("Environment", widget.NewLabel(s.connection.environment)),
						widget.NewFormItem("Multicast", widget.NewLabelWithData(s.connection.multicast)),
					),
				),
				container.NewVBox(
					widget.NewForm(
						widget.NewFormItem("Connected", connCheck),
						widget.NewFormItem("Registered", regCheck),
						widget.NewFormItem("Connecting to", widget.NewLabelWithData(s.connection.connecting)),
					),
				),
			),
		),
	))

	w.CenterOnScreen()

	w.ShowAndRun()
}
