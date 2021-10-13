package db

import (
	"flag"
	"fmt"
	"go-subscriptions-workflow/util"
	"os"
	"strconv"
)

var (
	database string
	username string
	password string
	hostname string
	port     int
)

type Config interface {
	URI() string
}

type config struct {
	database string
	username string
	password string
	hostname string
	port     int
}

func(c *config) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?w=majority", c.username, c.password, c.hostname, c.port, c.database)
}

func LoadConfigFromEnv() {
	username = os.Getenv("MONGODB_USERNAME")
	password = os.Getenv("MONGODB_PASSWORD")
	database = os.Getenv("MONGODB_DATABASE")
	hostname = os.Getenv("MONGODB_HOSTNAME")
	var err error
	port, err = strconv.Atoi(os.Getenv("MONGODB_PORT"))
	util.PanicOnError(err)
}

func LoadConfigFromFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&username, "mongodb_username", "myuser", "set mongodb password")
	flagSet.StringVar(&password, "mongodb_password", "mypass", "set mongodb user password")
	flagSet.StringVar(&database, "mongodb_database", "mydb", "set mongodb database name")
	flagSet.StringVar(&hostname, "mongodb_hostname", "localhost", "set mongodb hostname")
	flagSet.IntVar(&port, "mongodb_port", 21017, "set mongodb port")
}

func NewConfig() Config {
	return &config{
		database: database,
		username: username,
		password: password,
		hostname: hostname,
		port:     port,
	}
}
