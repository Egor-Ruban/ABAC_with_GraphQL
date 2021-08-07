package schema

import (
	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add subjects",
		Script: `
CREATE TABLE subjects (
	subject_id   UUID,
	name         TEXT,
	age         INT,
	role     TEXT,

	PRIMARY KEY (subject_id)
);`,
	},
	{
		Version:     2,
		Description: "Add objects",
		Script: `
CREATE TABLE objects (
	object_id   UUID,
	name         TEXT,
	owner_id         UUID,

	PRIMARY KEY (object_id)
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {

	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
