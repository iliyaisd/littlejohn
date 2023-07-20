package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/iliyaisd/littlejohn"
)

func main() {
	config, err := prepareConfig()
	if err != nil {
		log.Fatalf("Cannot prepare configuration: %s", err)
	}

	app, err := littlejohn.BuildApp()
	if err != nil {
		log.Fatalf("Cannot initialize Portfolio API: %s", err)
	}

	log.Printf("Portfolio API initialized\n")

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), app.MainHandler)
	if err != nil {
		log.Fatalf("Cannot listen and serve: %s\n", err.Error())
	}
}

func prepareConfig() (littlejohn.Config, error) {
	var config littlejohn.Config
	var err error

	config.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return littlejohn.Config{}, fmt.Errorf("cannot parse PORT: %w", err)
	}

	config.DataSource = os.Getenv("DATASOURCE")
	if len(config.DataSource) == 0 {
		config.DataSource = littlejohn.DataSourceLocal
	}

	return config, nil
}
