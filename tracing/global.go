package tracing

import (
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

const Version = "v0.1.0"

type Tracer struct {
}

func connect() (*nats.Conn, error) {
	server := viper.GetString("NATS")
	if server == "" {
		server = "nats://localhost:4222"
	}
	return nats.Connect(server)

}
