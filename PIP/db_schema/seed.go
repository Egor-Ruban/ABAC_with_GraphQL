package schema

import (
	"github.com/jmoiron/sqlx"
)

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Using a constant in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.

const seeds = `
INSERT INTO subjects (subject_id, name, age, role) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'Ivan', 50, 'ADMIN'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'Anders', 15, 'User')
	ON CONFLICT DO NOTHING;

INSERT INTO objects (object_id, name, owner_id) VALUES
	('e4e308e2-f307-11eb-9a03-0242ac130003', 'Comic Books', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b'),
	('efb32a04-f307-11eb-9a03-0242ac130003', 'McDonalds Toys', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b')
	ON CONFLICT DO NOTHING;
`

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
