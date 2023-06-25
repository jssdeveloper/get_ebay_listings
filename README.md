# Go eBay listing downloader
Program written in Go to download active ebay items to csv file without using eBay web page. You can use Cron job to create daily imports.

Used external packages:
Godotenv to load environment variables https://github.com/joho/godotenv

The main goal is to read active eBay listings using eBay API.

Before you start:
Make sure you have ebay API key. If you dont have it, please visit https://developer.ebay.com/ and register account.
After the account has been verified you can generate new production keyset in developer portal or can use oauth2 according to eBay documentation.

Make sure you have go version 1.20 or later installed on your machine
Clone the repository to your local machine
Paste the API key into .env file - replace your_api_key with your API key
In console while in active folder type and run go mod tidy
type and execute go run .
