module github.com/tuupke/pixie/server

go 1.18

require (
	github.com/google/flatbuffers v2.0.6+incompatible
	github.com/google/uuid v1.3.0
	github.com/hashicorp/mdns v1.0.5
	github.com/kevinburke/go.uuid v1.2.0
	github.com/nats-io/nats-server/v2 v2.8.4
	github.com/nats-io/nats.go v1.16.0
	github.com/rs/zerolog v1.27.0
	github.com/tuupke/pixie v0.0.0-20220503214133-75c988fe240a
	github.com/valyala/fasthttp v1.39.0
	gorm.io/driver/sqlite v1.3.6
	gorm.io/gorm v1.23.8
	openticket.tech/crud v1.2.0
	openticket.tech/db v1.3.2
	openticket.tech/env v1.5.0
	openticket.tech/lifecycle/v2 v2.0.1
	openticket.tech/null v1.1.0
	openticket.tech/rest/v3 v3.2.0
)

require (
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20220708163326-82d177caec6e // indirect
	github.com/adhityaramadhanus/fasthttpcors v0.0.0-20170121111917-d4c07198763a // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/fasthttp/router v1.4.11 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gocql/gocql v1.2.0 // indirect
	github.com/golang-migrate/migrate/v4 v4.15.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gookit/color v1.5.1 // indirect
	github.com/gookit/filter v1.1.3 // indirect
	github.com/gookit/goutil v0.5.8 // indirect
	github.com/gookit/validate v1.4.2 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.10.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.3.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/simukti/sqldb-logger v0.0.0-20220521163925-faf2f2be0eb6 // indirect
	github.com/simukti/sqldb-logger/logadapter/zerologadapter v0.0.0-20220521163925-faf2f2be0eb6 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/quicktemplate v1.7.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220812174116-3211cb980234 // indirect
	golang.org/x/sys v0.0.0-20220817070843-5a390386f1f2 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	golang.org/x/tools v0.1.12 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	openticket.tech/buffers v1.4.1 // indirect
	openticket.tech/certloader v1.0.0 // indirect
	openticket.tech/db/migrations v1.3.2 // indirect
	openticket.tech/gotrus/v3 v3.0.6 // indirect
	openticket.tech/iso8601 v1.1.1 // indirect
	openticket.tech/list v1.2.0 // indirect
	openticket.tech/log/v2 v2.2.0 // indirect
	openticket.tech/pubsub v1.3.2 // indirect
)

replace openticket.tech/db/migrations => ./migrations

replace openticket.tech/crud => ../../../go/src/openticket.tech/crud

replace openticket.tech/null => ../../../go/src/openticket.tech/null

replace github.com/tuupke/pixie => ../
