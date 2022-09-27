package pixie

import (
	"database/sql"

	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/tuupke/pixie/lifecycle"
)

var orm *gorm.DB

func Orm() *gorm.DB {
	if orm == nil {
		var err error

		// lifecycle.EFinally(sql.Close)
		//

		// sql.Conn()
		db, err := sql.Open("sqlite", envString("DB_DSN"))
		if err != nil {
			log.Fatal().Err(err).Msg("loading database")
		}

		lifecycle.EFinally(db.Close)
		orm, err = gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
		log.Err(err).Msg("loaded gorm")
		if err != nil {
			log.Fatal().Msg("gorm must boot")
		}
	}

	return orm
}
