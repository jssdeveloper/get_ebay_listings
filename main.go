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
	"sync"
	"time"
)

type Listing struct {
	ActiveList struct {
		ItemArray struct {
			Item []struct {
				ItemID        int `xml:"ItemID"`
				BuyItNowPrice struct {
					Text       string `xml:",chardata"`
					CurrencyID string `xml:"currencyID,attr"`
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

type ItemOut struct {
	ItemId     int
	Sku        string
	Price      float64
	Title      string
	WatchCount int
}

var allItems []ItemOut
var wg sync.WaitGroup

func main() {

	// get page count
	fmt.Println("Gettig page count")
	pages, err := getListings(1)
	if err != nil {
		panic(err)
	}
	pageCount := pages.ActiveList.PaginationResult.TotalNumberOfPages
	fmt.Println("Total number of pages:", pageCount, "Each page contains 200 listings")

	for i := 1; i <= pageCount; i++ {

		data, err := getListings(i)
		if err != nil {
			panic(err)
		}

		fmt.Println("Done Page", i, "Fetched", len(data.ActiveList.ItemArray.Item), "items")
	}
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

func getListings(page int) (Listing, error) {
	xmlPayload := xmlBody(page)

	req, err := http.NewRequest("POST", "https://api.ebay.com/ws/api.dll", bytes.NewReader(xmlPayload))
	if err != nil {
		panic(err)
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
		panic(err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	listing := &Listing{}

	_ = xml.Unmarshal([]byte(body), &listing)

	fmt.Println("Fetched page", page)

	return *listing, nil
}
