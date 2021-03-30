package pgxs

import (
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"strings"
)

//
// according to https://github.com/jackc/pgx/blob/master/conn.go#L84
// have to watch changes, to prevent internal issues
//
func QuoteString(str string) string {
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "%", "", -1)
	return str
}

func (db *Repo) SanitizeString(str string) string {
	return QuoteString(str)
}

// handleSqlErr used to avoid not exists and already exists debug queries
func (db *Repo) DebugLogSqlErr(q string, err error) error {
	pgErr, deuce := err.(*pgconn.PgError)
	if deuce {
		if pgErr.Code == "23505" {
			deuce = false
		}
	}

	if err != pgx.ErrNoRows && !deuce {
		db.Logger.Debugf("query: \n%s", q)
	}

	return err
}