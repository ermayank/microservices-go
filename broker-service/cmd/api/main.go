package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {}

func main() {
	app := Config{}

	log.Printf("Starting Broker Service on PORT %s\n", webPort);

	//Define Http Server
	srv := &http.Server {
		Addr : fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//Start Server

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}