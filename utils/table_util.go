package utils

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitTable(db *sql.DB, initFile string) error {
	c, ioErr := ioutil.ReadFile(initFile)
	if ioErr != nil {
		log.Fatal(ioErr)
	}

	sqlStr := string(c)
	_, err := db.Exec(sqlStr)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func DropTable(db *sql.DB, tables []string) error {
	query := "DROP TABLE IF EXISTS $1 CASCADE"

	var err error = nil
	for _, table := range tables {
		_, err = db.Exec(query, table)
	}

	return err
}

func DockerDBUp() (*sql.DB, *dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	_ = resource.Expire(220)

	var db *sql.DB
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db, pool, resource
}