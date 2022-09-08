package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
)

/*
MockDB is a mock database for testing.
Use MakeMockDB to create a MockDB, and call MockDB.Start() to init resource.
Under the hood, it use a docker container that runs mysql as data source.
So you need to keep docker running for the test to work.
You should call SetSchema() before each test to reset database schema.
*/
type MockDB struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	db       *sqlx.DB
	path     string
	schema   string
}

const SQL_FILE string = "saturday.sql"
const IMAGE_NAME string = "test_db"
const RESOURCE_NAME string = "test_db"

func (m MockDB) isImageExist(name string) (bool, error) {
	out, err := exec.Command("docker", "image", "ls").Output()
	if err != nil {
		return false, err
	}
	outStr := string(out)
	searched := strings.Contains(outStr, name)
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
	imageExists, _ := m.isImageExist(IMAGE_NAME)
	option := &dockertest.RunOptions{
		Name:       RESOURCE_NAME,
		Repository: IMAGE_NAME,
		Tag:        "latest",
	}
	if imageExists {
		res, ok := m.pool.ContainerByName(RESOURCE_NAME)
		if ok {
			m.resource = res
		} else {
			m.resource, err = m.pool.RunWithOptions(option)
		}
	} else {
		//TODO should just use mysql image since the schema is reset before each test
		m.resource, err = m.pool.BuildAndRunWithOptions(path.Join(m.path, "dockerfile"), option)
	}
	if err != nil {
		return nil, fmt.Errorf("could not start resource %s", err)
	}
	log.Printf("A docker container is created for database testing. For the convince running latter tests, it will not be stopped after the test, please close it manually.")
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
		b, err := os.ReadFile(path.Join(m.path, SQL_FILE))
		if err != nil {
			return err
		}
		m.schema = string(b)
	}
	m.db.MustExec(m.schema)
	return nil
}

func (m *MockDB) CloseDb() {
	m.db.Close()
}

func (m *MockDB) Close() {
	m.db.Close()
	if err := m.pool.Purge(m.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

/*
assetsPath should be the relative path to the assets folder,
MockDB needs to read the dockerfile located in the assets folder.
*/
func MakeMockDB(assetsPath string) *MockDB {
	return &MockDB{
		path: assetsPath,
	}
}
