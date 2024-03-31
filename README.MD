### tx-parser

This is a simple Ethereum transaction parser that can be used to extract information from Ethereum transactions. 

### Installation

```bash
go build -o tx-parser main.go
```

### Running

```bash
./tx-parser -config=config.json
```

### Usage

There are 3 provided APIs:

* `POST /api/v1/address` - Subscribe to an address
* `GET /api/v1/transactions` - Get all transactions for an address since subscription
* `GET /api/v1/transactions/current_block` - Get the current parsed block number