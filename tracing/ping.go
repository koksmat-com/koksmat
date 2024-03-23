package tracing

import (
	"log"
	"time"
)

func Ping() {
	nc, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Println("Pinging NATS")
	msg, err := nc.Request("ping", []byte("ping"), time.Duration(5)*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	if string(msg.Data) != "pong" {
		log.Fatal("Expected pong")
	}
	log.Println("Pong")
}

func Print(args ...interface{}) {
	log.Println(args...)

}
