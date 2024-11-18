package commentthan

import (
	"database/sql"

	"github.com/spotlightpa/moreofa/internal/db"
)

type service struct {
	db *sql.DB
	q  *db.Queries
}

func (app *appEnv) newService() error {
	if err := db.Migrate(app.dbname); err != nil {
		return err
	}
	dbase, err := db.Open(app.dbname)
	if err != nil {
		return err
	}
	app.srv = &service{
		db: dbase,
		q:  db.New(dbase),
	}
	return nil
}
