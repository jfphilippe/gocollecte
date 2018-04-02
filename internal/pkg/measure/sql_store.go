//  Copyright jean-françois PHILIPPE 2014-2018

package measure

import (
	"database/sql"
	"log"

	// Import de la librarie postgres
	_ "github.com/lib/pq"
)

const (
	// PGInsertStmt Insert for Postgresql
	PGInsertStmt = "INSERT INTO raw_datas(ts, node_id, sensor_id, value) VALUES (to_timestamp($1) AT TIME ZONE 'UTC',$2,$3,$4)"
)

// DbValuesWriter pour ecrire en base
type DbValuesWriter struct {
	insert string
	db     *sql.DB
}

// Close ferme le composant
func (w *DbValuesWriter) Close() error {
	log.Println("DbValuesWriter::Close")
	return w.db.Close()
}

// Write  Implements interface ValuesWriter
func (w *DbValuesWriter) Write(values *NodeValues) error {
	err := w.db.Ping()
	if err == nil {
		tx, err := w.db.Begin()
		if err == nil {
			stmt, err := tx.Prepare(w.insert)
			if err == nil {
				defer stmt.Close()
				for _, val := range values.Values {
					_, err = stmt.Exec(values.When, values.Node, val.Sensor, val.Value)
					if err != nil {
						_ = tx.Rollback()
						return err
					}
				}
				err = tx.Commit()
				return err
			}
		} else {
			_ = tx.Rollback()
		}
	}
	return err
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */
