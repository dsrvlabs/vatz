package plugin

import (
	"database/sql"
	"sync"

	// Load sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var (
	once = sync.Once{}

	db *pluginDB
)

type pluginEntry struct {
	Name       string
	Repository string
}

type dbWriter interface {
	AddPlugin(e pluginEntry) error
	DeletePlugin(name string) error
	UpdatePlugin() error
}

type dbReader interface {
	List() ([]pluginEntry, error)
	Get(name string) (*pluginEntry, error)
}

type pluginDB struct {
	conn *sql.Conn
	ctx  context.Context
}

func (p *pluginDB) AddPlugin(e pluginEntry) error {
	err := p.createPluginTable()
	if err != nil {
		return err
	}

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO plugin(name, repository) VALUES(?, ?)", e.Name, e.Repository)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
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

	_, err = tx.Exec("CREATE TABLE IF NOT EXISTS plugin (name varchar(256) PRIMARY KEY, repository varchar(256))")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *pluginDB) DeletePlugin(name string) error {
	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM plugin WHERE name = ?", name)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *pluginDB) UpdatePlugin() error {
	// TODO:
	return nil
}

func (p *pluginDB) List() ([]pluginEntry, error) {
	// TODO:
	return nil, nil
}

func (p *pluginDB) Get(name string) (*pluginEntry, error) {
	e := pluginEntry{}
	err := p.conn.
		QueryRowContext(p.ctx, "SELECT name, repository FROM plugin WHERE name=?", name).
		Scan(&e.Name, &e.Repository)

	if err != nil {
		return nil, err
	}

	return &e, err
}

func newWriter(dbfile string) (dbWriter, error) {
	log.Info().Str("module", "db").Msg("newWriter")

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			chanErr <- err
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
		err = <-chanErr
	}

	return db, err
}

func newReader(dbfile string) (dbReader, error) {
	log.Info().Str("module", "db").Msg("newReader")
	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			chanErr <- err
		}

		db = &pluginDB{ctx: ctx, conn: conn}
		chanErr <- nil
	})

	var err error
	if db == nil {
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
