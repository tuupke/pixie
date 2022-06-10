module github.com/tuupke/pixie/client

go 1.18

replace github.com/tuupke/pixie => ../

require (
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/google/flatbuffers v2.0.6+incompatible
	github.com/hashicorp/mdns v1.0.5
	github.com/kevinburke/go.uuid v1.2.0
	github.com/rs/zerolog v1.26.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/tuupke/pixie v0.0.0-20220417215231-7cf4054b3a86
	github.com/valyala/quicktemplate v1.7.0
	go.dedis.ch/kyber/v3 v3.0.13
)

require (
	github.com/beanstalkd/go-beanstalk v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/miekg/dns v1.1.41 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	go.dedis.ch/fixbuf v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e // indirect
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d // indirect
	golang.org/x/sys v0.0.0-20210809222454-d867a43fc93e // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
