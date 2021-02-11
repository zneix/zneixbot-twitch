package main

import (
	//"flag"
	"log"

	"github.com/zneix/zniksbot/mongo"
)

func main() {
	log.Println("Starting zniksbot!")

	//configFile := flag.String("config", "config.json", "json config file")

	mongo.Connect()
}
