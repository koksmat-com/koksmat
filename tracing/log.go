package tracing

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func ShipNats(args []string) error {
	subject := viper.GetString("NATS_SUBJECT")
	if subject == "" {
		return nil
	}
	data := strings.Join(args, " ")

	nc, err := connect()

	if err != nil {
		return err
	}
	defer nc.Close()
	return nc.Publish(subject, []byte(data))
}
func Log(args []string) {
	log.Println(args)
	err := ShipNats(args)
	if err != nil {
		log.Println("Cannot ship data to NATS", err)
	}

}
