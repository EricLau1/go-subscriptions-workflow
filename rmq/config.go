package rmq

import (
	"flag"
	"fmt"
	"go-subscriptions-workflow/util"
	"os"
	"strconv"
)

var (
	username string
	password string
	hostname string
	port     int
)

type Config interface {
	URL() string
}

type config struct {
	username string
	password string
	hostname string
	port     int
}

func (c *config) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.username, c.password, c.hostname, c.port)
}

func LoadConfigFromEnv() {
	username = os.Getenv("RABBITMQ_USERNAME")
	password = os.Getenv("RABBITMQ_PASSWORD")
	hostname = os.Getenv("RABBITMQ_HOSTNAME")
	var err error
	port, err = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	util.PanicOnError(err)
}

func LoadConfigFromFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&username, "rmq_username", "guest", "set rabbitmq password")
	flagSet.StringVar(&password, "rmq_password", "guest", "set rabbitmq user password")
	flagSet.StringVar(&hostname, "rmq_hostname", "localhost", "set rabbitmq hostname")
	flagSet.IntVar(&port, "rmq_port", 5672, "set amqp port")
}

func NewConfig() Config {
	return &config{
		username: username,
		password: password,
		hostname: hostname,
		port:     port,
	}
}
