package plugin

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	// Load sqlite3 driver
	"github.com/dsrvlabs/vatz/utils"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

var (
	once = sync.Once{}

	db *pluginDB
)

type pluginEntry struct {
	PluginID       string
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
		log.Info().Str("module", "db").Err(err)
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
		log.Error().Str("module", "db").Msgf("createPluginTable Error: %s", fieldError)
		return fieldError
	}

	if !isExist {
		opts := &sql.TxOptions{Isolation: sql.LevelDefault}
		tx, err := p.conn.BeginTx(p.ctx, opts)
		if err != nil {
			return err
		}

		q := `CREATE TABLE IF NOT EXISTS plugin_re (
		plugin_id varchar(256) PRIMARY KEY,
    	name varchar(256),
    	is_enabled int,
		repository varchar(256),
		binary_location varchar(256),
		version varchar(256),
		installed_at DATE)`

		_, err = tx.Exec(q)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return err
		}

		if err = tx.Commit(); err != nil {
			log.Info().Str("module", "db").Err(err)
			return err
		}

		log.Info().Str("module", "db").Msg("Update PluginTable")

		q = `SELECT name, repository, binary_location, version, installed_at FROM plugin`
		rows, err := p.conn.QueryContext(p.ctx, q)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
		}

		defer rows.Close()

		retPlugins := make([]pluginEntry, 0)

		for rows.Next() {
			e := pluginEntry{}
			err := rows.Scan(&e.Name, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
			if err != nil {
				log.Info().Str("module", "db").Err(err)
				return nil
			}
			retPlugins = append(retPlugins, e)
		}

		for _, v := range retPlugins {
			tx, err = p.conn.BeginTx(p.ctx, opts)
			if err != nil {
				log.Info().Str("module", "db").Err(err)
				return err
			}

			q = `INSERT INTO plugin_re(plugin_id, name, is_enabled, repository, binary_location, version, installed_at) VALUES(?, ?, ?, ?, ?, ?, ?)`

			hasValue := utils.UniqueHashValue(fmt.Sprintf("%s%s", v.Repository, v.Version))
			_, err = tx.Exec(q, hasValue, v.Name, 1, v.Repository, v.BinaryLocation, v.Version, v.InstalledAt)
			if err != nil {
				log.Error().Str("module", "A").Err(err)
				return err
			}

			if err = tx.Commit(); err != nil {
				log.Info().Str("module", "B").Err(err)
				return err
			}
		}

		tx, err = p.conn.BeginTx(p.ctx, opts)
		if err != nil {
			return err
		}

		q = `DROP TABLE IF EXISTS plugin`

		_, err = tx.Exec(q)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return err
		}
		q = `ALTER TABLE plugin_re RENAME TO plugin`
		_, err = tx.Exec(q)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return err
		}

		if err = tx.Commit(); err != nil {
			log.Info().Str("module", "db").Err(err)
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

	q := `INSERT INTO plugin(plugin_id, name, is_enabled, repository, binary_location, version, installed_at) VALUES(?, ?, ?, ?, ?, ?, ?)
`

	_, err = tx.Exec(q, e.PluginID, e.Name, e.IsEnabled, e.Repository, e.BinaryLocation, e.Version, e.InstalledAt)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
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
    plugin_id varchar(256) PRIMARY KEY,
	name varchar(256),
    is_enabled int,
	repository varchar(256),
	binary_location varchar(256),
	version varchar(256),
	installed_at DATE)
`

	_, err = tx.Exec(q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) DeletePlugin(name string) error {
	log.Info().Str("module", "db").Msg("DeletePlugin")

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	_, err = tx.Exec("DELETE FROM plugin WHERE name = ?", name)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) UpdatePlugin(pluginID string, isEnabled bool) error {
	// TODO: 1. Set best identifier for plugins either of Plugin_id or Name
	log.Info().Str("module", "db").Msgf("Plugin where plugin_id is %s", pluginID)

	opts := &sql.TxOptions{Isolation: sql.LevelDefault}
	tx, err := p.conn.BeginTx(p.ctx, opts)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	q := `UPDATE plugin SET is_enabled = ? WHERE plugin_id = ?`

	isEnabledInt := 0
	if isEnabled {
		isEnabledInt = 1
	}

	result, err := tx.Exec(q, isEnabledInt, pluginID)
	if err != nil {
		log.Error().Str("module", "db").Err(err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Error().Str("module", "db").Msgf("There's is no Plugin with plugin_id: %s. Please, check your plugins_id.", pluginID)
		return nil
	} else {
		response := "disabled"
		if isEnabled {
			response = "enabled"
		}

		log.Info().Str("module", "db").Msgf("Plugin with plugin_id %s has been %s.", pluginID, response)
	}

	if err = tx.Commit(); err != nil {
		log.Info().Str("module", "db").Err(err)
		return err
	}

	return nil
}

func (p *pluginDB) List() ([]pluginEntry, error) {
	log.Info().Str("module", "db").Msg("List Plugin")

	q := `SELECT plugin_id, name, is_enabled, repository, binary_location, version, installed_at FROM plugin`
	rows, err := p.conn.QueryContext(p.ctx, q)
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return nil, err
	}

	defer rows.Close()

	retPlugins := make([]pluginEntry, 0)

	for rows.Next() {
		e := pluginEntry{}
		err := rows.Scan(&e.PluginID, &e.Name, &e.IsEnabled, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return nil, err
		}

		retPlugins = append(retPlugins, e)
	}

	return retPlugins, nil
}

func (p *pluginDB) Get(identifier string) (*pluginEntry, error) {

	log.Info().Str("module", "db").Msgf("Get %s", identifier)

	q := `SELECT plugin_id, name,  is_enabled, repository, binary_location, version, installed_at FROM plugin WHERE plugin_id=? OR name=?`

	rows, err := p.conn.QueryContext(p.ctx, q, identifier, identifier)
	defer rows.Close()
	if err != nil {
		log.Info().Str("module", "db").Err(err)
		return nil, err
	}

	retPlugins := make([]pluginEntry, 0)
	for rows.Next() {
		e := pluginEntry{}
		err := rows.Scan(&e.PluginID, &e.Name, &e.IsEnabled, &e.Repository, &e.BinaryLocation, &e.Version, &e.InstalledAt)
		if err != nil {
			log.Info().Str("module", "db").Err(err)
			return nil, err
		}

		retPlugins = append(retPlugins, e)
	}

	if len(retPlugins) > 1 {
		log.Error().Str("module", "db").Msg("Please, start a plugin with plugin_id.")
		return nil, fmt.Errorf("There's more than one item to start with: %s", identifier)
	}

	getPlugin := retPlugins[0]

	return &getPlugin, nil

}

func newWriter(dbfile string) (dbWriter, error) {
	log.Info().Str("module", "db").Msgf("newWriter %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Create DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			log.Info().Str("module", "db").Msgf("Get conn Err %+v", err)
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
	log.Info().Str("module", "db").Msgf("newReader %s", dbfile)

	chanErr := make(chan error, 1)

	once.Do(func() {
		log.Info().Str("module", "db").Msg("Read DB Instance")

		ctx := context.Background()
		conn, err := getDBConnection(ctx, dbfile)
		if err != nil {
			log.Error().Str("module", "db").Err(err).Msg("")
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
