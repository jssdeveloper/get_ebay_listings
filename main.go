// Ebay Quantites sync
// By Janis Stals aka jssdeveloper 2023 under MIT license
// For further information please read LICENSE file

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	pwd         string     // print working directory
	maxQuantity int    = 5 // Change this to adjust maximum quantity items of single SKU to be active after update
)

type Listing struct {
	ActiveList struct {
		ItemArray struct {
			Item []struct {
				ItemID            string `xml:"ItemID"`
				Title             string `xml:"Title"`
				WatchCount        string `xml:"WatchCount"`
				QuantityAvailable string `xml:"QuantityAvailable"`
				SKU               string `xml:"SKU"`
				PictureDetails    string `xml:"PictureDetails"`
			} `xml:"Item"`
		} `xml:"ItemArray"`
		PaginationResult struct {
			TotalNumberOfPages   string `xml:"TotalNumberOfPages"`
			TotalNumberOfEntries string `xml:"TotalNumberOfEntries"`
		} `xml:"PaginationResult"`
	} `xml:"ActiveList"`
}

type AllListings struct {
	Id     int
	ItemId string
	Sku    string
}

func init() {
	getPath()
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

func getPath() {
	// Gets active directory of the exec file
	path, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to load .env from env_app/.env")
		panic(err)
	}
	pwd = filepath.Dir(path)
}

func loadEnv() {
	err := godotenv.Load(path.Join(pwd, ".env"))
	if err != nil {
		fmt.Println("Error loading env file (env_app/.env)")
		// os.Exit(1)
	}
	fmt.Println(".env file loaded")
}

// Function to get particular page
func xmlBody(page int) []byte {
	body := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<GetMyeBaySellingRequest xmlns="urn:ebay:apis:eBLBaseComponents">   
		<ErrorLanguage>en_US</ErrorLanguage>
		<WarningLevel>High</WarningLevel>
	<ActiveList>
		<Sort>TimeLeft</Sort>
		<Pagination>
		<EntriesPerPage>200</EntriesPerPage>
		<PageNumber>%v</PageNumber>
		</Pagination>
	</ActiveList>
	</GetMyeBaySellingRequest>`, page))
	return body
}

func GetAllListings(apiKey string) ([]AllListings, error) {
	var api_key string = os.Getenv("EBAY_API_KEY")
	var page int = 1
	var id int = 1

	allListings := []AllListings{}

	for {
		// Define the XML payload
		xmlPayload := xmlBody(page)

		// Create a new HTTP request with the XML payload
		req, err := http.NewRequest("POST", "https://api.ebay.com/ws/api.dll", bytes.NewReader(xmlPayload))
		if err != nil {
			return nil, err
		}

		// Set the request headers
		req.Header.Set("X-EBAY-API-SITEID", "77")
		req.Header.Set("X-EBAY-API-COMPATIBILITY-LEVEL", "967")
		req.Header.Set("X-EBAY-API-CALL-NAME", "GetMyeBaySelling")
		req.Header.Set("X-EBAY-API-IAF-TOKEN", api_key)

		// Create a new HTTP client with a timeout
		client := &http.Client{
			Timeout: time.Second * 10,
		}

		// Send the request and get the response
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		listing := &Listing{}

		_ = xml.Unmarshal([]byte(body), &listing)

		getitem := listing.ActiveList.ItemArray.Item
		TotalPages := listing.ActiveList.PaginationResult.TotalNumberOfPages
		totalPagesInt, err := strconv.Atoi(TotalPages)
		if err != nil {
			return nil, err
		}
		fmt.Println(totalPagesInt)

		for _, v := range getitem {
			itemId := v.ItemID
			sku := v.SKU

			currentListing := AllListings{id, itemId, sku}
			allListings = append(allListings, currentListing)
			id++
		}
		fmt.Println("Page", page, "of", totalPagesInt)
		fmt.Println("-------------")

		if page == totalPagesInt {
			break
		}
		page++
	}

	return allListings, nil
}
