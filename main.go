// Ebay Quantites sync
// By Janis Stals aka jssdeveloper 2023 under MIT license
// For further information please read LICENSE file

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
)

// For ease of use I dont follow do not follow DRY principle in this program all the time.

var (
	pwd         string     // print working directory
	maxQuantity int    = 5 // Change this to adjust maximum quantity items of single SKU to be active after update
)

type EbayItem struct {
	Sku    string
	EbayId string
}

func init() {
	// Gets active directory of the exec file
	path, err := filepath.Abs("./main.go")
	if err != nil {
		panic(err)
	}
	pwd = filepath.Dir(path)

	loadEnv()

}

func main() {

	// Use this if you would like to read quantities from csv file
	// csv()

	// Use this if you would like to read quantities from sqlite database
	// sqlite()

	// Use this if you would like to read quantities from postgres database
	// postgres()

}

func loadEnv() {
	err := godotenv.Load(path.Join(pwd, "env_app", ".env"))
	if err != nil {
		fmt.Println("Error loading env file (env_app/.env)")
		os.Exit(1)
	}
	fmt.Println(".env file loaded")
}

// change the required data
func csv() {

}

func sqlite() {

}

func postgres() {

}
