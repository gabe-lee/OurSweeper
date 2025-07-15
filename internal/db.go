package internal

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

type (
	Time = time.Time
)

const (
	SQL_DIR    string = "sql"
	SQLITE_DIR string = "sqlite"
	DB_NAME    string = "our_sweeper.db"
	DB_PERMS          = 0755

	Q_READ string = `
SELECT name FROM ` + DB_NAME + ` WHERE type='table' AND name='{table_name}';
	`
)

const (
	DB_EXEC_CREATE_PROP_TABLE_IF_NEEDED string = `
CREATE TABLE IF NOT EXISTS props (
  id INTEGER PRIMARY KEY, 
  int INTEGER,
  str TEXT
);`
	DB_QUER_VER             string = "SELECT int FROM props WHERE id=" + DB_PROP_VER + ";"
	DB_EXEC_SET_VER_CURRENT string = "UPDATE props SET int=" + DB_PROP_VER_CURRENT_STR + " WHERE id=" + DB_PROP_VER + ";"
	DB_EXEC_ADD_VER         string = "INSERT INTO props (id, int) VALUES (" + DB_PROP_VER + ", 0);"

	DB_EXEC_DROP_ALL string = `
DELETE FROM props;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS login_tokens;
DROP TABLE IF EXISTS participation;
DROP TABLE IF EXISTS flags;
DROP TABLE IF EXISTS worlds;
`
)

const (
	DB_PROP_VER             string = "1"
	DB_PROP_VER_CURRENT     int    = 1
	DB_PROP_VER_CURRENT_STR string = "1"
)

const (
	DB_VER_1 int64 = iota
	DB_VER_COUNT

	DB_VER_CURRENT = DB_VER_1
)

var UP_MIGRATIONS = [DB_VER_COUNT]string{
	DB_VER_1: `
CREATE TABLE users (
  id_a INTEGER NOT NULL,
  id_b INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  login TEXT NOT NULL,
  email_hash TEXT NOT NULL,
  pass_hash TEXT NOT NULL,
  use_name INTEGER DEFAULT 0 NOT NULL,
  screen_name TEXT,
  playtime INTEGER DEFAULT 0 NOT NULL,
  score INTEGER DEFAULT 0 NOT NULL,
  score_sweeps INTEGER DEFAULT 0 NOT NULL,
  score_flags INTEGER DEFAULT 0 NOT NULL,
  sweeps INTEGER DEFAULT 0 NOT NULL,
  deaths INTEGER DEFAULT 0 NOT NULL,
  good_flags INTEGER DEFAULT 0 NOT NULL,
  warnings INTEGER DEFAULT 0 NOT NULL,
  banned INTEGER DEFAULT 0 NOT NULL,
  ban_reason TEXT DEFAULT '(not banned)' NOT NULL,
  PRIMARY KEY(id_a, id_b)
);

CREATE TABLE login_tokens (
  token_a INTEGER NOT NULL,
  token_b INTEGER NOT NULL,
  user_id_a INTEGER NOT NULL,
  user_id_b INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,
  PRIMARY KEY(token_a, token_b),
  FOREIGN KEY (user_id_a, user_id_b) REFERENCES users (id_a, id_b) 
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE participation (
  user_id_a INTEGER NOT NULL,
  user_id_b INTEGER NOT NULL,
  world_id_a INTEGER NOT NULL,
  world_id_b INTEGER NOT NULL,
  score INTEGER DEFAULT 0 NOT NULL,
  sweeps INTEGER DEFAULT 0 NOT NULL,
  good_flags INTEGER DEFAULT 0 NOT NULL,
  deaths INTEGER DEFAULT 0 NOT NULL,
  FOREIGN KEY (user_id_a, user_id_b) REFERENCES users (id_a, id_b) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (world_id_a, world_id_b) REFERENCES worlds (id_a, id_b) 
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE flags (
  user_id_a INTEGER NOT NULL,
  user_id_b INTEGER NOT NULL,
  world_id_a INTEGER NOT NULL,
  world_id_b INTEGER NOT NULL,
  pos INTEGER NOT NULL,
  good BOOLEAN NOT NULL,
  FOREIGN KEY (user_id_a, user_id_b) REFERENCES users (id_a, id_b) 
    ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (world_id_a, world_id_b) REFERENCES worlds (id_a, id_b) 
    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE worlds (
  id_a INTEGER NOT NULL,
  id_b INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,
  cleared_at INTEGER NOT NULL,
  expired INTEGER DEFAULT 0 NOT NULL,
  cleared INTEGER DEFAULT 0 NOT NULL,
  participants INTEGER DEFAULT 0 NOT NULL,
  total_score INTEGER DEFAULT 0 NOT NULL,
  tile_state BLOB,
  PRIMARY KEY(id_a, id_b)
);`,
}

var DB_DIR = filepath.Join(SQL_DIR, SQLITE_DIR)
var DB_PATH = filepath.Join(SQL_DIR, SQLITE_DIR, DB_NAME)

type SweepDB struct {
	Db  *sql.DB
	Ver int64
}

func (SweepDB) CheckFile() {
	err := os.MkdirAll(DB_DIR, os.ModeDir|0755)
	if err != nil {
		log.Fatalf("could not find/create database directory '%s': %s", DB_DIR, err)
	}
	_, err = os.Stat(DB_PATH)
	if err != nil {
		file, err := os.Create(DB_PATH)
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatalf("unable to close database file '%s': %s", DB_PATH, err)
			}
		}()
		if err != nil {
			log.Fatalf("unable to create database file '%s': %s", DB_PATH, err)
		}
	}
}

func (s *SweepDB) Open() {
	d, err := sql.Open("sqlite", DB_PATH)
	if err != nil {
		log.Fatalf("could not open database '%s': %s", DB_PATH, err)
	}
	s.Db = d
	_, err = s.Db.Exec(DB_EXEC_CREATE_PROP_TABLE_IF_NEEDED)
	if err != nil {
		log.Fatalf("could not create or verify the 'props' table in database: %s", err)
	}
	row := s.Db.QueryRow(DB_QUER_VER)
	var ver int64
	err = row.Scan(&ver)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.Db.Exec(DB_EXEC_DROP_ALL)
			if err != nil {
				log.Fatalf("unable to drop all database tables: %s", err)
			}
			_, err := s.Db.Exec(DB_EXEC_ADD_VER)
			if err != nil {
				log.Fatalf("unable to set 'props.db_ver' (id "+DB_PROP_VER+") to ver 0: %s", ver, err)
			}
		} else {
			log.Fatalf("unable to read the 'db version' property from the database: %s", err)
		}
	}
	if ver < int64(DB_VER_COUNT) {
		for ver < int64(DB_VER_COUNT) {
			mig := UP_MIGRATIONS[ver]
			_, err = s.Db.Exec(mig)
			ver += 1
			if err != nil {
				log.Fatalf("unable to perform UP migration VER %d -> %d: %s", ver-1, ver, err)
			}
		}
		_, err = s.Db.Exec(DB_EXEC_SET_VER_CURRENT)
		if err != nil {
			log.Fatalf("unable to set 'props.db_ver' (id 1) to ver %d: %s", ver, err)
		}
	}
	s.Ver = ver
}

func (s *SweepDB) Close() {
	err := s.Db.Close()
	if err != nil {
		log.Fatalf("Failed to close database: %s", err)
	}
	s.Db = nil
	s.Ver = 0
}
