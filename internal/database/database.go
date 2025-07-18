package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	"github.com/gabe-lee/OurSweeper/attempt_group"
	"github.com/gabe-lee/OurSweeper/coord"
	"github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/logger"
)

type (
	Time         = time.Time
	AttemptGroup = attempt_group.AttemptGroup
	Timeout      = attempt_group.Timeout
	Logger       = logger.Logger
	SubLogger    = logger.SubLogger
	ServerWorld  = common.ServerWorld
	Tile         = common.Tile
)

const (
	SQL_DIR    string = "sql"
	SQLITE_DIR string = "sqlite"
	DB_NAME    string = "our_sweeper.db"
	DB_PERMS   int    = 0755
)

const (
	DB_EXEC_CREATE_PROP_TABLE_IF_NEEDED string = `
CREATE TABLE IF NOT EXISTS props (
  id INTEGER PRIMARY KEY, 
  int INTEGER,
  str TEXT
);`

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
	DB_PROP_VER_CURRENT int = 1
)

const (
	DB_VER_1 int64 = iota
	DB_VER_COUNT

	DB_VER_CURRENT = DB_VER_1
)

var PROP_NAMES = [...]string{
	DB_PROP_VER_CURRENT: "version",
}

var UP_MIGRATIONS = [DB_VER_COUNT]string{
	DB_VER_1: `
CREATE TABLE IF NOT EXISTS props (
  id INTEGER PRIMARY KEY, 
  int INTEGER,
  str TEXT
);

CREATE TABLE users (
  id INTEGER PRIMARY KEY,
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
  total_flags INTEGER DEFAULT 0 NOT NULL,
  good_flags INTEGER DEFAULT 0 NOT NULL,
  warnings INTEGER DEFAULT 0 NOT NULL,
  banned INTEGER DEFAULT 0 NOT NULL,
  ban_reason TEXT DEFAULT '(not banned)' NOT NULL
);

CREATE TABLE anon_users (
  uuid_a INTEGER NOT NULL,
  uuid_b INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,
  playtime INTEGER DEFAULT 0 NOT NULL,
  score INTEGER DEFAULT 0 NOT NULL,
  score_sweeps INTEGER DEFAULT 0 NOT NULL,
  score_flags INTEGER DEFAULT 0 NOT NULL,
  sweeps INTEGER DEFAULT 0 NOT NULL,
  deaths INTEGER DEFAULT 0 NOT NULL,
  good_flags INTEGER DEFAULT 0 NOT NULL,
  PRIMARY KEY (uuid_a, uuid_b)
);

CREATE TABLE login_tokens (
  token INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE TABLE participation (
  user_id INTEGER NOT NULL,
  world_id INTEGER NOT NULL,
  score INTEGER DEFAULT 0 NOT NULL,
  sweeps INTEGER DEFAULT 0 NOT NULL,
  good_flags INTEGER DEFAULT 0 NOT NULL,
  deaths INTEGER DEFAULT 0 NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (world_id) REFERENCES worlds (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE TABLE flags (
  user_id INTEGER NOT NULL,
  world_id INTEGER NOT NULL,
  pos INTEGER NOT NULL,
  good INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (world_id) REFERENCES worlds (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION
);

CREATE TABLE worlds (
  id INTEGER PRIMARY KEY,
  seed_a INTEGER NOT NULL,
  seed_b INTEGER NOT NULL,
  difficulty INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  expires_at INTEGER NOT NULL,
  cleared_at INTEGER NOT NULL,
  expired INTEGER DEFAULT 0 NOT NULL,
  cleared INTEGER DEFAULT 0 NOT NULL,
  participants INTEGER DEFAULT 0 NOT NULL,
  total_score INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE chunks (
  world_id INTEGER NOT NULL,
  idx INTEGER NOT NULL,
  data BLOB NOT NULL,
  FOREIGN KEY (world_id) REFERENCES worlds (id) 
    ON DELETE CASCADE ON UPDATE NO ACTION
);`,
}

var DB_DIR = filepath.Join(SQL_DIR, SQLITE_DIR)
var DB_PATH = filepath.Join(SQL_DIR, SQLITE_DIR, DB_NAME)

type SweepDB struct {
	Db        *sql.DB
	WriteLock sync.Mutex
	Ver       int64
	Log       SubLogger
}

func NewSweepDB(masterLogger *Logger) SweepDB {
	return SweepDB{
		Log: masterLogger.NewSubLogger("Database"),
	}
}

func (*SweepDB) CheckFile() {
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

const DB_EXEC_OPEN_SETTINGS = `
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache = shared;
PRAGMA temp_store = memory;`

func (s *SweepDB) Open() {
	d, err := sql.Open("sqlite", DB_PATH)
	s.Log.FatalIfErr(err, "could not open database '%s'", DB_PATH)
	s.Db = d
	s.Db.SetMaxOpenConns(256)
	s.Db.SetMaxIdleConns(64)
	s.Db.SetConnMaxIdleTime(time.Second * 60)
	_, err = s.Db.Exec(DB_EXEC_OPEN_SETTINGS)
	s.Log.WarnIfErr(err, "could not set database settings: %s", DB_EXEC_OPEN_SETTINGS)
	_, err = s.Db.Exec(DB_EXEC_CREATE_PROP_TABLE_IF_NEEDED)
	s.Log.FatalIfErr(err, "could not create or verify the 'props' table in database")
	ver, err := s.GetPropNum(DB_PROP_VER_CURRENT)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = s.Db.Exec(DB_EXEC_DROP_ALL)
			s.Log.FatalIfErr(err, "unable to drop all database tables")
			err = s.SetPropNum(DB_PROP_VER_CURRENT, 0)
			s.Log.FatalIfErr(err, "unable to initialize database 'version' property to 0")
		} else {
			s.Log.FatalIfErr(err, "unable to read the 'db version' property from the database")
		}
	}
	if int64(ver) < int64(DB_VER_COUNT) {
		for int64(ver) < int64(DB_VER_COUNT) {
			mig := UP_MIGRATIONS[ver]
			_, err = s.Db.Exec(mig)
			ver += 1
			s.Log.FatalIfErr(err, "unable to perform UP migration VER %d -> %d", ver-1, ver)
		}
		err = s.SetPropNum(DB_PROP_VER_CURRENT, ver)
		s.Log.FatalIfErr(err, "unable to set database version property")
	}
	s.Ver = int64(ver)
}

func (s *SweepDB) Close() {
	err := s.Db.Close()
	s.Log.WarnIfErr(err, "failed to close database")
	s.Db = nil
	s.Ver = 0
}

const DB_QUERY_GET_PROP_NUM string = `
SELECT int
FROM props
WHERE id=$1;`

func (s *SweepDB) GetPropNum(code int) (int, error) {
	timeout := attempt_group.NewTimeout(time.Second * 5)
	row := s.Db.QueryRowContext(timeout, DB_QUERY_GET_PROP_NUM, code)
	var val int
	err := row.Scan(&val)
	s.Log.WarnIfErr(err, "could not read property '%s' integer value from database", PROP_NAMES[code])
	return val, err
}

const DB_QUERY_GET_PROP_STR string = `
SELECT str
FROM props
WHERE id=$1;`

func (s *SweepDB) GetPropStr(code int) (string, error) {
	timeout := attempt_group.NewTimeout(time.Second * 5)
	row := s.Db.QueryRowContext(timeout, DB_QUERY_GET_PROP_NUM, code)
	var val string
	err := row.Scan(&val)
	s.Log.WarnIfErr(err, "could not read property '%s' string value from database", PROP_NAMES[code])
	return val, err
}

const DB_QUERY_SET_PROP_NUM string = `
INSERT INTO props (id, int)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET int=$2;`

func (s *SweepDB) SetPropNum(code int, num int) error {
	timeout := attempt_group.NewTimeout(time.Second * 5)
	_, err := s.Db.ExecContext(timeout, DB_QUERY_SET_PROP_NUM, code, num)
	s.Log.WarnIfErr(err, "could not set property '%s' val %d to database", PROP_NAMES[code], num)
	return err
}

const DB_QUERY_SET_PROP_STR string = `
INSERT INTO props (id, str)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET str=$2;`

func (s *SweepDB) SetPropStr(code int, str string) error {
	timeout := attempt_group.NewTimeout(time.Second * 5)
	_, err := s.Db.ExecContext(timeout, DB_QUERY_SET_PROP_STR, code, str)
	s.Log.WarnIfErr(err, "could not set prop '%s' val '%s' to database", PROP_NAMES[code], str)
	return err
}

const DB_QUERY_CREATE_WORLD string = `
INSERT INTO worlds (seed_a, seed_b, difficulty, created_at, expires_at, cleared_at)
VALUES ($1, $2, $3, $4, $5, $5) RETURNING id;`

func (s *SweepDB) CreateNewWorld(world *ServerWorld, difficulty byte) {
	timeGo := time.Now()
	created := timeGo.Unix()
	expires := timeGo.Add(time.Hour * 24).Unix()
	var seed_a uint64 = *(*uint64)(unsafe.Pointer(&created))
	var seed_b uint64 = *(*uint64)(unsafe.Pointer(&expires))
	worldRow := s.Db.QueryRow(DB_QUERY_CREATE_WORLD, seed_a, seed_b, difficulty, created, expires, expires)
	var worldId uint32
	err := worldRow.Scan(&worldId)
	s.Log.FatalIfErr(err, "unable to create new world")
	world.InitNew(worldId, difficulty, expires, seed_a, seed_b)
	result := attempt_group.NewWithTimeout("create chunks for new world", time.Second*10, int64(common.WORLD_CHUNK_COUNT))
	for idx := range common.WORLD_CHUNK_COUNT {
		go s.CreateChunk(&result, world, idx)
	}
	s.Log.WarnIfErr(result.Wait(), "failed to create all new chunks in database")
}

const DB_QUERY_CREATE_CHUNK string = `
INSERT INTO chunks (world_id, idx, data)
VALUES ($1, $2, $3);`

func (s *SweepDB) CreateChunk(result *AttemptGroup, world *ServerWorld, idx int) {
	data := world.CopyChunk(idx)
	id := world.Id.Load()
	s.WriteLock.Lock()
	defer s.WriteLock.Unlock()
	_, err := s.Db.ExecContext(result.Timeout, DB_QUERY_CREATE_CHUNK, id, idx, data[:])
	if err != nil {
		s.Log.WarnIfErr(err, "failed to create world %d chunk %d in database", id, idx)
		result.Failure()
	} else {
		result.Success()
	}
}

const DB_QUERY_UPDATE_CHUNK string = `
UPDATE chunks
SET data=$3
WHERE world_id=$1 AND idx=$2;`

func (s *SweepDB) UpdateChunk(result *AttemptGroup, world *ServerWorld, idx int) {
	data := world.CopyChunk(idx)
	id := world.Id.Load()
	s.WriteLock.Lock()
	defer s.WriteLock.Unlock()
	_, err := s.Db.ExecContext(result.Timeout, DB_QUERY_UPDATE_CHUNK, id, idx, data[:])
	if err != nil {
		s.Log.WarnIfErr(err, "failed to create world %d chunk %d in database", id, idx)
		result.Failure()
	} else {
		result.Success()
	}
}

func (s *SweepDB) UpdateAllChunks(world *ServerWorld) {
	result := attempt_group.NewWithTimeout("UpdateAllChunks", time.Second*10, int64(common.WORLD_CHUNK_COUNT))
	for idx := range common.WORLD_CHUNK_COUNT {
		go s.UpdateChunk(&result, world, idx)
	}
	s.Log.WarnIfErr(result.Wait(), "failed to update all chunks in database")
}

const DB_QUERY_LOAD_CHUNKS string = `
SELECT idx, data
FROM chunks
WHERE world_id=$1;`

func (s *SweepDB) LoadAllChunks(world *ServerWorld) {
	timeout := attempt_group.NewTimeout(time.Second * 10)
	id := world.Id.Load()
	rows, err := s.Db.QueryContext(timeout, DB_QUERY_LOAD_CHUNKS, id)
	s.Log.WarnIfErr(err, "failed to load any chunks for world %d from database", id)
	defer rows.Close()
	var idx int
	var data []byte
	var cnt int
	for rows.Next() {
		err = rows.Scan(&idx, &data)
		s.Log.WarnIfErr(err, "failed to scan chunk data")
		cPos := coord.CoordFromIndex(idx, common.CY_SHIFT, common.CX_MASK)
		tTopLeft := cPos.ShiftUpScalar(common.TILE_TO_LOCK_SHIFT)
		bIdx := 0
		for y := range common.WORLD_TILES_PER_CHUNK_AXIS {
			for x := range common.WORLD_TILES_PER_CHUNK_AXIS {
				tPos := tTopLeft.AddXY(x, y)
				tIdx := tPos.ToIndex(common.TY_SHIFT)
				world.Tiles[tIdx] = Tile(data[bIdx])
				bIdx += 1
			}
		}
		cnt += 1
	}
	if cnt != common.WORLD_CHUNK_COUNT {
		s.Log.Warn("only loaded %d/%d chunks from database for world %d", cnt, common.WORLD_CHUNK_COUNT, id)
	}
}

const DB_QUERY_GET_ACTIVE_WORLD string = `
SELECT id, expires_at, participants, total_score
FROM worlds
WHERE expired=0 AND cleared=0 AND difficulty=$1;`

func (s *SweepDB) GetActiveWorld(world *ServerWorld, difficulty byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	row := s.Db.QueryRowContext(ctx, DB_QUERY_GET_ACTIVE_WORLD, difficulty)
	var id, part, score uint32
	err := row.Scan(
		&id,
		&world.Expires,
		&part,
		&score,
	)
	if err != nil {
		return false
	}
	world.Id.Store(id)
	world.Participants.Store(part)
	world.Score.Store(score)
	return err == nil
}

const DB_QUERY_CHECK_LOGIN_TOKEN string = `
SELECT 
`
