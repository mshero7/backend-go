package main

import (
	"backend-go/db/ent"
	"context"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"

	"log"
)

func serverLoad() {
	web.Run()
}

func dbLoad() {
	client, err := ent.Open("postgres", "host=<host> port=<port> user=<user> dbname=<database> password=<pass>")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

func main() {
	serverLoad()
	dbLoad()
}
