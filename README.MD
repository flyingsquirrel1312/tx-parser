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

### Configuration

The configuration file is a JSON file that contains the following fields as shown below:

```json
{
  "rpc_url": "https://cloudflare-eth.com",
  "port": 8080, // HTTP Port
  "buffer_size": 1000, // Maximum number of transactions can be stored by the parser
  "ticker_period": 10, // Time in seconds to check for new blocks
  "start_block": 19552784 // Block number to start parsing from
}
```

### Example:

Refer to TxParser.postman_collection.json for an example of how to use the APIs.

### TODO:
- Add retry mechanism for RPC calls
- Use a proper database to store transactions. My recommendation would be to use a time-series database like ClickHouse.