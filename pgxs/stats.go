package pgxs

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"strings"
)

var (
	ErrNotExist     = fmt.Errorf("schema or table does not exist")
	ErrAlreadyExist = fmt.Errorf("schema or table already exist")
)

// TableStats describes some statistics for a table.
type TableStats struct {
	Table       string
	TableType   string
	SizeTotal   int32
	SizeIndexes int32
	SizeTable   int32
	Rows        int32
}

// DBStats describes some statistics for a database.
type DBStats struct {
	Name         string
	CountTables  int32
	CountRows    int32
	SizeTotal    int64
	SizeIndexes  int64
	SizeSchema   int64
	CountIndexes int32
}

// Schemas returns a sorted list of PostgreSQL schema names.
func (db *Repo) Schemas(ctx context.Context) ([]string, error) {
	sql := "SELECT schema_name FROM information_schema.schemata ORDER BY schema_name"
	rows, err := db.Pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0, 2)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}

		if strings.HasPrefix(name, "pg_") || name == "information_schema" {
			continue
		}

		res = append(res, name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Tables returns a sorted list of specified schema PostgreSQL table names.
func (db *Repo) Tables(ctx context.Context, schemaName string) ([]string, error) {
	q := "SELECT table_name FROM information_schema.tables WHERE table_schema = $1 ORDER BY table_name"
	rows, err := db.Pool.Query(ctx, q, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0, 2)
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}

		res = append(res, name)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateSchema creates a new PostgreSQL schema.
//
// It returns ErrAlreadyExist if schema already exist.
func (db *Repo) CreateSchema(ctx context.Context, schemaName string) error {
	q := `CREATE SCHEMA ` + pgx.Identifier{schemaName}.Sanitize()
	_, err := db.Pool.Exec(ctx, q)

	if e, ok := err.(*pgconn.PgError); ok && e.Code == pgerrcode.DuplicateSchema {
		return ErrAlreadyExist
	}

	return err
}

// DropSchema drops PostgreSQL schema.
//
// It returns ErrNotExist if schema does not exist.
func (db *Repo) DropSchema(ctx context.Context, schemaName string) error {
	sql := `DROP SCHEMA ` + pgx.Identifier{schemaName}.Sanitize() + ` CASCADE`
	_, err := db.Pool.Exec(ctx, sql)

	if e, ok := err.(*pgconn.PgError); ok && e.Code == pgerrcode.InvalidSchemaName {
		return ErrNotExist
	}

	return err
}

// CreateTable creates PostgreSQL jsonb table.
//
// It returns ErrAlreadyExist if table already exist.
func (db *Repo) CreateTable(ctx context.Context, schemaName, tableName string) error {
	sql := `CREATE TABLE ` + pgx.Identifier{schemaName, tableName}.Sanitize()
	_, err := db.Pool.Exec(ctx, sql)

	if e, ok := err.(*pgconn.PgError); ok && e.Code == pgerrcode.DuplicateTable {
		return ErrAlreadyExist
	}

	return err
}

// DropTable drops PostgreSQL table.
//
// It returns ErrNotExist is table does not exist.
func (db *Repo) DropTable(ctx context.Context, schemaName, tableName string, useCascade bool) error {
	q := `DROP TABLE ` + pgx.Identifier{schemaName, tableName}.Sanitize()
	if useCascade {
		q += ` CASCADE`
	}
	_, err := db.Pool.Exec(ctx, q)

	if e, ok := err.(*pgconn.PgError); ok && e.Code == pgerrcode.UndefinedTable {
		return ErrNotExist
	}

	return err
}

// TableStats returns a set of statistics for specified schema table.
func (db *Repo) TableStats(ctx context.Context, schemaName, tableName string) (*TableStats, error) {
	res := new(TableStats)
	sql := `SELECT table_name, table_type,
           pg_total_relation_size('"'||t.table_schema||'"."'||t.table_name||'"'),
           pg_indexes_size('"'||t.table_schema||'"."'||t.table_name||'"'),
           pg_relation_size('"'||t.table_schema||'"."'||t.table_name||'"'),
           COALESCE(s.n_live_tup, 0)
      FROM information_schema.tables AS t
      LEFT OUTER
      JOIN pg_stat_user_tables AS s ON s.schemaname = t.table_schema
                                      and s.relname = t.table_name
     WHERE t.table_schema = $1
       AND t.table_name = $2`

	err := db.Pool.QueryRow(ctx, sql, schemaName, tableName).
		Scan(&res.Table, &res.TableType, &res.SizeTotal, &res.SizeIndexes, &res.SizeTable, &res.Rows)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DBStats returns a set of statistics for a specified schema.
func (db *Repo) DBStats(ctx context.Context, schemaName string) (*DBStats, error) {
	res := new(DBStats)
	sql := `SELECT COUNT(distinct t.table_name)                                                     AS CountTables,
           COALESCE(SUM(s.n_live_tup), 0)                                                           AS CountRows,
           COALESCE(SUM(pg_total_relation_size('"'||t.table_schema||'"."'||t.table_name||'"')), 0)  AS SizeTotal,
           COALESCE(SUM(pg_indexes_size('"'||t.table_schema||'"."'||t.table_name||'"')), 0)         AS SizeIndexes,
           COALESCE(SUM(pg_relation_size('"'||t.table_schema||'"."'||t.table_name||'"')), 0)        AS SizeSchema,
           COUNT(distinct i.indexname)                                                              AS CountIndexes
      FROM information_schema.tables AS t
      LEFT OUTER
      JOIN pg_stat_user_tables       AS s ON s.schemaname = t.table_schema
                                         AND s.relname = t.table_name
      LEFT OUTER
      JOIN pg_indexes                AS i ON i.schemaname = t.table_schema
                                         AND i.tablename = t.table_name
     WHERE t.table_schema = $1`

	res.Name = schemaName
	err := db.Pool.QueryRow(ctx, sql, schemaName).
		Scan(&res.CountTables, &res.CountRows, &res.SizeTotal, &res.SizeIndexes, &res.SizeSchema, &res.CountIndexes)
	if err != nil {
		return nil, err
	}

	return res, nil
}
