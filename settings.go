package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// ----------- YOU CAN CHANGE SETTINGS FOR ENV LOCATION --------------
const (
	env_path string = ".env" // path to env file
)

// -----------  GETS ENV FILE AND STORES API KEY IN VARIABLE --------------
var ebay_api_key string

func init() {
	err := godotenv.Load(env_path)
	if err != nil {
		fmt.Println("Error loading env file (.env)")
		os.Exit(1)
	}
	fmt.Println(".env file loaded")
	ebay_api_key = os.Getenv("EBAY_API_KEY")
}
