# Check if node is ready
curl -X GET "http://localhost:9650/ext/health"

curl -X POST "http://localhost:9650/ext/info" -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"info.uptime","params":{"chain":"A"},"id":1}' 

curl -s -X POST http://127.0.0.1:9650/ext/bc/O -H 'content-type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"omega.getCurrentValidators","params":{"subnetID":"11111111111111111111111111111111LpoYY"}}'

curl -s -X POST http://127.0.0.1:9650/ext/bc/O -H 'content-type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"omega.getPendingValidators","params":{"subnetID":"11111111111111111111111111111111LpoYY"}}'
# Check node health status
curl -X POST "http://localhost:9650/ext/health" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"health.health","params":{"tags":[]},"id":1}'

# Check node liveness
curl -X POST "http://localhost:9650/ext/health" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"health.liveness","params":{"tags":[]},"id":1}'


curl -X POST "http://localhost:9650/ext/info" \   
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "info.getBlockchainID",
    "params": {
      "alias": "D"
    },
    "id": 1
  }'

  # Get node version
curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"info.getNodeVersion","params":{},"id":1}'

# Get node ID
curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"info.getNodeID","params":{},"id":1}'

# Get network ID
curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"info.getNetworkID","params":{},"id":1}'

  curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"info.peers","params":{},"id":1}'

# Check if chain is bootstrapped
curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"info.isBootstrapped","params":{"chain":"A"},"id":1}'


# Health check
curl -X GET "http://localhost:9650/ext/health"

# Info API
curl -X POST "http://localhost:9650/ext/info" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "info.getBlockchainID",
    "params": {
      "alias": "D"
    },
    "id": 1
  }'

curl -X POST -H "Content-Type: application/json" -d '{                                         
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["0x1", false],  
    "id": 1
  }' http://localhost:9650/ext/bc/D/rpc


curl -X POST "http://localhost:9650/ext/bc/D/rpc" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_blockNumber",
    "params": [],
    "id": 1
  }'

  curl -X POST "http://localhost:9650/ext/bc/D/rpc" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["latest", true],
    "id": 1
  }'


  curl -X POST "http://localhost:9650/ext/bc/D/rpc" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0x3470a4336413ed7a1f7a5fb06ffb2d2102873ef0", "latest"],
    "id": 1
  }'

curl -s -X POST http://127.0.0.1:9650/ext/bc/O -H 'content-type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"omega.getUTXOs","params":{"addresses":["O-testnet1klqj9yz44y2flc5dyxum23nn5ucr35xycwfqxq"]}}'

curl -s -X POST http://127.0.0.1:9650/ext/bc/O -H 'content-type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"omega.getBalance","params":{"addresses":["O-testnet1klqj9yz44y2flc5dyxum23nn5ucr35xycwfqxq"]}}'

curl -s -X POST http://127.0.0.1:9650/ext/bc/O -H 'content-type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"omega.getTxStatus","params":{"txID":"<importTxID>"}}'

curl -s -X POST -H 'content-type: application/json' --data '{
  "jsonrpc":"2.0","id":1,"method":"omega.getMinStake",
  "params":{"subnetID":"11111111111111111111111111111111LpoYY"}
}' http://localhost:9650/ext/bc/O

curl -s -H 'content-type: application/json' -d '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xe4175f400f5b748943ea9a59ece822f19a1312cca0bbdd9ff66043d7641a5d1a"],"id":1}' http://45.55.93.124:9650/ext/bc/D/rpc | jq .


./build/odysseygo --network-id=testnet --chain-config-dir="scripts/configs/archive" --http-port=9662 --staking-port=9663 --db-dir=./data/node1 --log-dir=./data/node1/logs --data-dir=./data/node1/.odysseygo

./build/odysseygo --network-id=123456 --http-port=9652 --staking-port=9653 --db-dir=./data/node2 --log-dir=./data/node2/logs --data-dir=./data/node2/.odysseygo

./build/odysseygo --network-id=123456 --http-port=9653 --staking-port=9654 --db-dir=./data/node3 --log-dir=./data/node3/logs --data-dir=./data/node3/.odysseygo
--bootstrap-ids string                                            Comma separated list of bootstrap peer ids to connect to. Example: NodeID-JR4dVmy6ffUGAKCBDkyCbeZbyHQBeDsET,NodeID-8CrVPQZ4VSqgL8zTdvL14G8HqAfrBr4z
      --bootstrap-ips string 

--bootstrap-ids=NodeID-DFRpSjRPPByYm2jtX74aCpdrcstAfpAwR --bootstrap-ips=127.0.0.1:9651