package functions

import (
	"database/sql"

	"github.com/appointy/idgen"
	"github.com/lib/pq"
)

func UpsertModuleRole(db *sql.DB, moduleName, displayName string, patterns []string) error {

	const get = `SELECT id FROM appointy_module_v1.module_role WHERE "name"=$1 AND is_default=true`
	const insert = `INSERT INTO appointy_module_v1.module_role (id, "name", display_name, pattern, is_default) VALUES ($1, $2, $3, $4, true)`
	const update = `UPDATE appointy_module_v1.module_role SET pattern=$1 WHERE id=$2`

	var mrID string

	if err := db.QueryRow(get, moduleName).Scan(&mrID); err != nil {

		if err == sql.ErrNoRows { // insert
			mrID = idgen.New("mdr")
			if _, err := db.Exec(insert, mrID, moduleName, displayName, pq.Array(patterns)); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		// update
		if _, err := db.Exec(update, pq.Array(patterns), mrID); err != nil {
			return err
		}
	}

	return nil
}
