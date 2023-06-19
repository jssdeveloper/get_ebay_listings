# Sync Ebay Quantities with csv file or databse
Program written in Go to sync eBay quantities with database or csv file (csv, sqlite, postgres, gorm, godotenv)

1. Download the reository with command: git clone https://github.com/jssdeveloper/go_ebay_sync
If you use windows, please make sure git is installed.

2. To use this program you need ebay developer account. You can register here https://developer.ebay.com/.
After your account has been submited you have to create API keyset and generate API token.

3. Copy ebay developer credentials in env_github\.env and replace neccesary fields

4. You can input item stock data from csv file, postgres database or sqlite database. The program outputs csv file to be uploaded to ebay.


Used packages:
Godotenv to load environment variables https://github.com/joho/godotenv
Gorm ORM to interact with database gorm.io/gorm
