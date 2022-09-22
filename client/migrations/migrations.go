package migrations

import (
	"embed"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed */*.sql
var fs embed.FS

func GetMigrations(name string) (string, source.Driver, error) {
	d, err := iofs.New(fs, strings.ToLower(name))
	return "iofs", d, err
}
