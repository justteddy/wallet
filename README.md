# Wallet service

## Prerequisites
Docker and docker-compose are required to run locally

## How to start/stop service in docker
Start service `make up`

Stop service `make down`

## How to start in dev mode with database in docker
Build binary `go build -o ./wallet`

Start database in docker `docker-compose up -d db`

Run binary `./wallet`

## Maintenance commands
Run unit tests - `make test`

Run integration tests - `make integration-test`

Linter check (golangci-lint should be installed) - `make lint`


## Usage
Once launched, the service is available on `http://localhost:8080`.

It has the following endpoints:

### 1. Create wallet

`POST /wallet`

Creates new wallet, doesn't have any params

Request example:
```
curl --location --request POST 'http://localhost:8080/wallet'
```
Response example:

`200 OK`
```
{
    "wallet_id": "ab2ee047683d8880849f89a581298f139e90f1668bd5fa67f1f7e593ac64bea9"
}
```

### 2. Deposit

`POST /deposit/:wallet`

Increases wallet balance by amount value

Body payload:
```
{
    "amount": 100 // required, integer, accepts as input the amount specified in cents
}
```

Request example:
```
curl --location --request POST 'http://localhost:8080/deposit/cd56cabcac74c9ea4520b4cfc6ab9ab4089554b40f24735f25b0c518ed5a8164' \
--header 'Content-Type: application/json' \
--data-raw '{
    "amount": 100
}'
```
Response example: `HTTP 200 OK` with empty body

### 3. Transfer

`POST /transfer`

Transfers amount value from one wallet to another

Body payload:
```
{
    "from_wallet": "walet1", // required, string
    "to_wallet": "wallet2",  // required, string
    "amount": 100            // required, integer, accepts as input the amount specified in cents
}
```

Request example:
```
curl --location --request POST 'http://localhost:8080/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from_wallet": "97e7da3986d84a35cbcb6cc2ce8ac3bcc07337ab169435f707245de440d4c297",
    "to_wallet": "107e9e098a3587b18a5d44aca58e25255e2afeb96971f59b346481879863acfe",
    "amount": 100
}'
```
Response example: `HTTP 200 OK` with empty body

### 4. Report

`POST /report/:format/:wallet`

Reports wallet operations in specified format `json|csv` with optional filters

Body payload:
```
{
    "from_date": "2030-12-30",           // optional, string, date in format YYYY-MM-DD
    "to_date": "2030-12-31",             // optional, string, date in format YYYY-MM-DD
    "operation_type": "deposit|withdraw" // optional, string, "deposit" or "withdraw"
}
```

Request example:

JSON

```
curl --location --request POST 'http://localhost:8080/report/json/95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from_date": "",
    "to_date": "2021-11-25",
    "operation_type": ""
}'
```
CSV
```
curl --location --request POST 'http://localhost:8080/report/csv/95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from_date": "",
    "to_date": "2021-11-25",
    "operation_type": ""
}'
```
Response example:

`200 OK`

JSON

```
[
    {
        "wallet_id": "95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4",
        "operation_type": "deposit",
        "amount": "3.00$",
        "date": "2021-11-25"
    },
    {
        "wallet_id": "95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4",
        "operation_type": "deposit",
        "amount": "123.12$",
        "date": "2021-11-25"
    },
    {
        "wallet_id": "95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4",
        "operation_type": "deposit",
        "amount": "22.22$",
        "date": "2021-11-25"
    }
]
```
CSV
```
wallet_id,operation_id,amount,date
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,1.00$,2021-11-25
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,1.00$,2021-11-25
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,1.00$,2021-11-25
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,3.00$,2021-11-25
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,123.12$,2021-11-25
95e0fde5b14cd8f9afc48ec8d87ddd3570dc9ccc0fd32214090ae98c132e5be4,deposit,22.22$,2021-11-25

```
