package plugin

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var (
	once = sync.Once{}

	db *pluginDB
)

type pluginEntry struct {
	Name           string
	IsEnabled      int
	Repository     string
	BinaryLocation string
	Version        string
	InstalledAt    time.Time
}

type dbWriter interface {
	MigratePluginTable() error
	AddPlugin(e pluginEntry) error
	DeletePlugin(name string) error
	UpdatePlugin(pluginID string, isEnabled bool) error
}

type dbReader interface {
	List() ([]pluginEntry, error)
	Get(identifier string) (*pluginEntry, error)
}

type pluginDB struct {
	conn *sql.Conn
	ctx  context.Context
}

func (p *pluginDB) isFieldExist(tableName string, columnName string) (bool, error) {
	q := `PRAGMA table_info(plugin)`
	rows, err := p.conn.QueryContext(p.ctx, q)
	if err != nil {
		log.Error().Str("module", "db").Msgf("isFieldExist Error: %s", err)
		return false, err
	}

	defer rows.Close()

	// Check the columns of the result set
	found := false
	for rows.Next() {
		var cid int
		var name string
		var typename string
		var notnull bool
		var dfltvalue sql.NullString
		var pk bool

		err = rows.Scan(&cid, &name, &typename, &notnull, &dfltvalue, &pk)
		if err != nil {
			panic(err)
		}

		if name == "is_enabled" {
			found = true
			break
		}
	}
	return found, nil
}

func (p *pluginDB) MigratePluginTable() error {
	isExist, fieldError := p.isFieldExist("plugin", "is_enabled")

	if fieldError != nil {
		log.Error().Str("module", "db").Msgf("MigratePluginTable > isFieldExist Error: %s", fieldError)
		return fieldError
	}

	if !isExist {
		opts := &sql.TxOptions{Isolation: sql.LevelDefault}
		tx, err := p.conn.BeginTx(p.ctx, opts)
		if err != nil {
			return err
		}
		q := `ALTER TABLE plugin ADD COLUMN is_enabled INTEGER DEFAULT 1`

		_, err = tx.Exec(q)
		if err != nil {
			log.Error().Str("module", "db").Msgf("MigratePluginTable > tx.Exec Error: %s", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			log.Error().Str("module", "db").Msgf("MigratePluginTable > tx.Commit Error: %s", err)
			return err
		}
	}

	return nil
}

func (p *pluginDB) AddPlugin(e pluginEntry) error {

	log.Info().Str("module", "db").Msg("AddPlugin")

	opts := &sql.TxOptions{
		Isolation: sql.LevelDefault,
	}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	q := `INSERT INTO plugin(name, is_enabled, repository, binary_location, version, installed_at) VALUES(?, ?, ?, ?, ?, ?)
`

	_, err = tx.Exec(q, e.Name, e.IsEnabled, e.Repository, e.BinaryLocation, e.Version, e.InstalledAt)
	if err != nil {
		log.Error().Str("module", "db").Msgf("AddPlugin > tx.Exec Error: %s", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Str("module", "db").Msgf("AddPlugin > tx.Commit Error: %s", err)
		return err
	}

	return nil
}

func (p *pluginDB) createPluginTable() error {

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		return err
	}

	q := `
CREATE TABLE IF NOT EXISTS plugin ( 
	name varchar(256) PRIMARY KEY,
    is_enabled int,
	repository varchar(256),
	binary_location varchar(256),
	version varchar(256),
	installed_at DATE)
`
	_, err = tx.Exec(q)
	if err != nil {
		log.Error().Str("module", "db").Msgf("createPluginTable > tx.Exec Error: %s", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Str("module", "db").Msgf("createPluginTable > tx.Commit Error: %s", err)
		return err
	}

	return nil
}

func (p *pluginDB) DeletePlugin(name string) error {
	log.Info().Str("module", "db").Msg("DeletePlugin")

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Error().Str("module", "db").Msgf("DeletePlugin > conn.BeginTx Error: %s", err)
		return err
	}

	_, err = tx.Exec("DELETE FROM plugin WHERE name = ?", name)
	if err != nil {
		log.Error().Str("module", "db").Msgf("DeletePlugin > tx.Exec Error: %s", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Str("module", "db").Msgf("DeletePlugin > tx.Commit Error: %s", err)
		return err
	}

	return nil
}

func (p *pluginDB) UpdatePlugin(name string, isEnabled bool) error {
	// TODO: 1. Set best identifier for plugins either of Plugin_id or Name
	log.Info().Str("module", "db").Msg("UpdatePlugin")

	response := "disabled"
	if isEnabled {
		response = "enabled"
	}

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	q := `UPDATE plugin SET is_enabled = ? WHERE name = ?`

	isEnabledInt := 0
	if isEnabled {
		isEnabledInt = 1
	}

	_, err = tx.Exec(q, isEnabledInt, name)
	if err != nil {
		log.Error().Str("module", "db").Msgf("UpdatePlugin > tx.Exec Error: %s", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Error().Str("module", "db").Msgf("UpdatePlugin > tx.Commit Error: %s", err)
		return err
	}

	log.Info().Str("module", "db").Msgf("Plugin %s has been updated: %s.", name, response)

	return nil
}

func (p *pluginDB) List() ([]pluginEntry, error) {
	log.Info().Str("module", "db").Msg("List")

	q := `SELECT name, is_enabled, repository, binary_location, version, installed_at FROM plugin`
	rows, err := p.conn.QueryContext(p.ctx, q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return nil, err
	}

	defer rows.Close()
	retPlugins := make([]pluginEntry, 0)

	for rows.Next() {
		e := pluginEntry{}
		err := rows.Scan(&e.Name, &e.IsEnabled, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
		if err != nil {
			log.Error().Str("module", "db").Msgf("List > rows.Scan Error: %s", err)
			return nil, err
		}
		retPlugins = append(retPlugins, e)
	}

	return retPlugins, nil
}

func (p *pluginDB) Get(name string) (*pluginEntry, error) {
	log.Info().Str("module", "db").Msgf("Get %s", name)

	q := `SELECT name, is_enabled, repository, binary_location, version, installed_at FROM plugin WHERE name=?`
	e := pluginEntry{}
	err := p.conn.
		QueryRowContext(p.ctx, q, name).
		Scan(&e.Name, &e.IsEnabled, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
	if err != nil {
		log.Error().Str("module", "db").Msgf("Get > QueryRowContext.Scan Error: %s", err)
		return nil, err
	}

	return &e, err

}

func newWriter(dbfile string) (dbWriter, error) {
	//log.Debug().Str("module", "db").Msgf("newWriter %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			log.Error().Str("module", "db").Msgf("newWriter > getDBConnection Error: %s", err)
			chanErr <- err
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
		log.Info().Str("module", "db").Msg("Wait creation")
		err = <-chanErr
	}

	return db, err
}

func newReader(dbfile string) (dbReader, error) {

	//log.Debug().Str("module", "db").Msgf("newReader %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Read DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			log.Error().Str("module", "db").Msgf("newReader > getDBConnection Error: %s", err)
			chanErr <- err
			return
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
		log.Info().Str("module", "db").Msg("Wait creation")
		err = <-chanErr
	}

	return db, err
}

func getDBConnection(ctx context.Context, dbfile string) (*sql.Conn, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initDB(dbfile string) error {
	log.Info().Str("module", "db").Msgf("initDB %s", dbfile)

	if db != nil {
		db.conn.Close()
		db = nil

		once = sync.Once{}
	}

	err := os.Remove(dbfile)
	if err != nil && !os.IsNotExist(err) {
		log.Info().Err(err)
		return err
	}

	path := filepath.Dir(dbfile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir(path, 0755)
	}

	_, err = newWriter(dbfile)
	if err != nil {
		return err
	}

	err = db.createPluginTable()
	if err != nil {
		return err
	}

	return nil
}
