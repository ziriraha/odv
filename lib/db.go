package lib

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func GetDefaultDBPrefix() string {
	return GetConfig().Database.Prefix
}

func runOdooDBCommand(args ...string) error {
	_, err := runOdooCommand(append([]string{"db"}, args...)...)
	return err
}

func getConn(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, GetConfig().Database.DSN)
}

func DropDBs(dbNames ...string) error {
	for _, dbName := range dbNames {
		if dbName == "" {
			continue
		}
		err := runOdooDBCommand("drop", dbName)
		if err != nil {
			return fmt.Errorf("failed to drop odoo database %s: %w", dbName, err)
		}
	}
	return nil
}

func ListDBs(prefix string) ([]string, error) {
	ctx := context.Background()
	conn, err := getConn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, "SELECT datname FROM pg_database WHERE datistemplate = false AND datname LIKE $1", prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		dbs = append(dbs, strings.TrimSpace(name))
	}
	return dbs, nil
}
