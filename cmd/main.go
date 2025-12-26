package main

import (
	"flag"
	"log"

	"UPS2Telegram/internal/config"
)

func main() {
	env := flag.String("env", "", "Environment")
	flag.Parse()
	if *env == "" {
		log.Fatalf("Input parameter -env is nessesary")
	}
	if *env != "local" && *env != "prod" {
		log.Fatalf("invalid env: '%s'", *env)

	}
	
	cfg, err := config.Load(*env)
	if err != nil {
		log.Fatalf("%v", err)
	}
	
	// TODO: remove
	if cfg != nil {
		log.Println("Done!")
	}
}