# oracle-monitoring
Monitoring tool for DIA Oracles.

## Testing

There is a `cmd/oracle-monitoring/main_test.go` test suite that runs `main.go` file in background and tests the functionality by sending `setValue` transactions to the oracle.

Before running the test: 
- You must have correctly configured .env file.
- The wallet that you configured in .env file must have enough balance to pay for 2x `setValue` transactions.
- The wallet must have permission to call `setValue` function on the DIA oracle smart contract.

## Running

This project uses Docker, so it's pretty easy to run it:

(make sure commands are ran from root directory of the project)

1. First you need to prepare docker image:
```
docker build -t dia-oracles-monitoring-tool .
```

2. Now you can run the docker image:
```
docker run dia-oracles-monitoring-tool
```


