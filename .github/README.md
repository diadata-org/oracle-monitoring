## Features

- concurrent scraping of multiple oracle contracts
- the data is saved to a local TimescaleDB instance
- packaged in a Docker container

## Todo

- lint go files
- log the scraping events
- docker

## Usage

Copy the file `env.example` to `.env`, in the root directory, and then fill it.

```json
[
    {
        "contract-address": "0xa93546947f3015c986695750b8bbea8e26d65856",
        "contract-abi": "oracle-v2",
        "chain-id": "1",
        "node-url": "https://eth-mainnet.g.alchemy.com/v2/",
        "creation-block": 0,
        "latest-scraped-block": 0
    }
]
```

```shell
./build/oracle-monitoring --targets oracles.json
```

## Compile

```shell
go build -o build/
```

docker for both:

- the daemon
- timescaledb

## Deployment

```shell
DB_CONNECTION_STRING='postgres://user:password@localhost:5432/dbname'
psql -h <database_host> -U <database_user> -f scripts/create_database.sql
```

##

for now, it uses the same ABI for all the oracles

might add a property in the config file with the path to the matching ABI