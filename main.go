// Ebay Quantites sync
// By Janis Stals 2023 under MIT license
// For further information please read LICENSE file

package main

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// eBay response struct
type Listing struct {
	ActiveList struct {
		ItemArray struct {
			Item []struct {
				ItemID        int `xml:"ItemID"`
				BuyItNowPrice struct {
					Text       float64 `xml:",chardata"`
					CurrencyID string  `xml:"currencyID,attr"`
				} `xml:"BuyItNowPrice"`
				Title             string `xml:"Title"`
				WatchCount        int    `xml:"WatchCount"`
				QuantityAvailable int    `xml:"QuantityAvailable"`
				SKU               string `xml:"SKU"`
				PictureDetails    string `xml:"PictureDetails"`
			} `xml:"Item"`
		} `xml:"ItemArray"`
		PaginationResult struct {
			TotalNumberOfPages   int `xml:"TotalNumberOfPages"`
			TotalNumberOfEntries int `xml:"TotalNumberOfEntries"`
		} `xml:"PaginationResult"`
	} `xml:"ActiveList"`
}

// Struct to save data to csv
type ItemOut struct {
	ItemId            int
	Sku               string
	Price             float64
	Title             string
	WatchCount        int
	QuantityAvailable int
}

var allItems []ItemOut

func main() {

	// Get page count
	fmt.Println("Gettig page count")
	pages, err := getListings(1)
	if err != nil {
		panic(err)
	}
	pageCount := pages.ActiveList.PaginationResult.TotalNumberOfPages
	if pageCount < 1 {
		fmt.Println("Check your API key and make sure ebay account has active listings!")
		os.Exit(1)
	} else {
		fmt.Println("Total number of pages:", pageCount, "Each page contains 200 listings")
	}

	// Runs main script to download data from eBay
	fmt.Println("Starting to download eBay listings..")
	for i := 1; i <= 2; i++ {

		data, err := getListings(i)
		if err != nil {
			panic(err)
		}

		for _, v := range data.ActiveList.ItemArray.Item {
			item := ItemOut{ItemId: v.ItemID, Sku: v.SKU, Price: v.BuyItNowPrice.Text, Title: v.Title, WatchCount: v.WatchCount, QuantityAvailable: v.QuantityAvailable}
			allItems = append(allItems, item)
		}

		fmt.Println("Done Page", i, "Fetched", len(data.ActiveList.ItemArray.Item), "items")
	}

	// Creates output csv. To change file location, please go to settings.go and change csv_path
	err = createCsv()
	if err != nil {
		panic(err)
	}

	fmt.Println("Ebay import finished!")
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

// Main function to download eBay data
func getListings(page int) (Listing, error) {
	fmt.Println("Fetching page", page)
	xmlPayload := xmlBody(page)

	req, err := http.NewRequest("POST", "https://api.ebay.com/ws/api.dll", bytes.NewReader(xmlPayload))
	if err != nil {
		fmt.Println("Error connecting to ebay")
		os.Exit(1)
	}

	// Set the request headers
	req.Header.Set("X-EBAY-API-SITEID", "77")
	req.Header.Set("X-EBAY-API-COMPATIBILITY-LEVEL", "967")
	req.Header.Set("X-EBAY-API-CALL-NAME", "GetMyeBaySelling")
	req.Header.Set("X-EBAY-API-IAF-TOKEN", ebay_api_key)

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return Listing{}, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Listing{}, err
	}

	listing := &Listing{}

	_ = xml.Unmarshal([]byte(body), &listing)

	return *listing, nil
}

// Function to create output csv file. To change csv file path, go to settings.go and change csv_path
func createCsv() error {
	file, err := os.Create(csv_path)
	if err != nil {
		fmt.Println("Error creating csv file")
		return (err)
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	csvWriter.Write([]string{"ItemId", "Sku", "Price", "Title", "Watch Count", "Quantity Available"})
	for _, v := range allItems {
		itemId := fmt.Sprintf("%v", v.ItemId)
		price := fmt.Sprintf("%v", v.Price)
		wathcCount := fmt.Sprintf("%v", v.WatchCount)
		quantityAvailable := fmt.Sprintf("%v", v.QuantityAvailable)
		csvWriter.Write([]string{itemId, v.Sku, price, v.Title, wathcCount, quantityAvailable})
	}
	fmt.Println("Writing CSV file completed!")
	return nil
}
