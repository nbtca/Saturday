package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

type MockDB struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	db       *sqlx.DB
	path     string
	schema   string
}

func (m MockDB) isImageExist(name string) (bool, error) {
	out, err := exec.Command("docker", "image", "ls").Output()
	if err != nil {
		return false, err
	}
	outStr := string(out)
	searched := strings.Contains(outStr, "test_db")
	return searched, nil
}

func (m *MockDB) Start() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	m.pool, err = dockertest.NewPool("")
	m.pool.MaxWait = time.Minute * 2
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}
	imageExists, _ := m.isImageExist("test_db")
	if imageExists {
		m.resource, err = m.pool.Run("test_db", "latest", []string{})
	} else {
		//TODO should just use mysql image since the schema is reset before each test
		m.resource, err = m.pool.BuildAndRun("test_db", path.Join(m.path, "dockerfile"), []string{})
	}
	if err != nil {
		return nil, fmt.Errorf("could not start resource %s", err)
	}
	if err = m.pool.Retry(func() error {
		db, err = sqlx.Connect("mysql", fmt.Sprintf("root:password@(localhost:%s)/saturday_test?multiStatements=true", m.resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker %s", err)
	}
	db.SetMaxOpenConns(1000) // The default is 0 (unlimited)
	db.SetMaxIdleConns(10)   // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(0) // 0, connections are reused forever.
	m.db = db
	return db, nil
}

func (m *MockDB) SetSchema() error {
	if m.schema == "" {
		b, err := ioutil.ReadFile(path.Join(m.path, "saturday.sql"))
		if err != nil {
			return err
		}
		m.schema = string(b)
	}
	m.db.MustExec(m.schema)
	return nil
}

func (m *MockDB) Close() {
	m.db.Close()
	if err := m.pool.Purge(m.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
func MakeMockDB(assetsPath string) *MockDB {
	return &MockDB{
		path: assetsPath,
	}
}
