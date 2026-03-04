package lib

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var DBMutex sync.Mutex

func runDBCommand(name string, args ...string) (string, error) {
	DBMutex.Lock()
	defer DBMutex.Unlock()
	return runCommand(name, args...)
}

func DropDB(dbName string) error {
	_, err := runDBCommand("dropdb", "--if-exists", "--force", dbName)
	if err != nil {
		return fmt.Errorf("failed to drop database %s: %v", dbName, err)
	}
	err = os.RemoveAll(GetFilestorePath() + "/" + dbName)
	if err != nil {
		return fmt.Errorf("failed to remove filestore for database %s: %v", dbName, err)
	}
	return nil
}

func CreateDB(dbName string) error {
	_, err := runDBCommand("createdb", dbName)
	return err
}

func DuplicateDB(sourceDB, newDB string) error {
	_, err := runDBCommand("createdb", "-T", sourceDB, newDB)
	if err != nil {
		return fmt.Errorf("failed to create database %s from template %s: %v", newDB, sourceDB, err)
	}
	sourceFilestore := GetFilestorePath() + "/" + sourceDB
	newFilestore := GetFilestorePath() + "/" + newDB
	if _, err := os.Stat(newFilestore); err == nil {
		return fmt.Errorf("filestore for new database %s already exists", newDB)
	}
	return os.CopyFS(newFilestore, os.DirFS(sourceFilestore))
}

func ListDBs(prefix string) ([]string, error) {
	if prefix == "" {
		prefix = "rd-"
	}
	output, err := runDBCommand("psql", "-d", "postgres", "-t", "-c", "SELECT datname FROM pg_database WHERE datname LIKE '"+prefix+"%';")
	if err != nil {
		return nil, err
	}
	var dbs []string
	for line := range strings.SplitSeq(string(output), "\n") {
		dbName := strings.TrimSpace(line)
		if dbName != "" {
			dbs = append(dbs, dbName)
		}
	}
	return dbs, nil
}
