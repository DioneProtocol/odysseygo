# Contract Value Reader

A Go script to read values from the BasicContract deployed at address `0x9` in the Odyssey genesis.

## 📋 Complete Process

### Step 1: Create Contract ABI

1. **Compile the Solidity contract:**
```bash
cd contracts
solc --abi BasicContract.sol
```

2. **Save the ABI to a file:**
```bash
solc --abi BasicContract.sol > BasicContract.abi
```

3. **Generate bytecode:**
```bash
solc --bin BasicContract.sol > BasicContract.bin
```

### Step 2: Add Contract to Custom Genesis

1. **Open the genesis file:**
```bash
nano genesis/genesis_custom.json
```

2. **Add contract to dChainGenesis alloc section:**
```json
{
  "dChainGenesis": "{\"config\":{...},\"alloc\":{\"5b9a1ef78f359e73b14142878704ac831e9500a8\":{\"balance\":\"0x295BE96E64066972000000\"},\"0000000000000000000000000000000000000009\":{\"code\":\"0x6080604052613e805f553480156013575f5ffd5b5061012b806100215f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c806360fe47b11460345780636d4ce63c14604c575b5f5ffd5b604a60048036038101906046919060ab565b6066565b005b60526070565b604051605d919060de565b60405180910390f35b606f5f8190555050565b5f6103e8905090565b5f5ffd5b5f819050919050565b608d81607d565b81146096575f5ffd5b50565b5f8135905060a5816086565b92915050565b5f6020828403121560bd5760bc6079565b5b5f60c8848285016099565b91505092915050565b60d881607d565b82525050565b5f60208201905060ef5f83018460d1565b9291505056fea26469706673582212202a534f8c5f270dcade1d0da1ced7544f9d453afcf985c9d8b121940f390ff3c764736f6c634300081e0033\",\"balance\":\"0x0\",\"storage\":{\"0x0000000000000000000000000000000000000000000000000000000000000000\":\"0x000000000000000000000000000000000000000000000000000000000001d97c\"}}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"
}
```

3. **Key components added:**
   - **Contract address:** `0000000000000000000000000000000000000009` (0x9)
   - **Bytecode:** The compiled contract code
   - **Storage:** Initial value `0x1d97c` (121212 in decimal)

### Step 3: Start Odyssey Node with Custom Genesis

1. **Build the project:**
```bash
go build -o build/odysseygo
```

2. **Start node with custom genesis:**
```bash
./build/odysseygo --network-id=123456 --genesis-file=genesis/genesis_custom.json --http-port=9650 --staking-port=9651 --db-dir=./data/node1 --log-dir=./data/node1/logs --data-dir=./data/node1/.odysseygo
```

3. **Wait for node to sync and start RPC services**

### Step 4: Run the Value Reader Script

1. **Navigate to contracts directory:**
```bash
cd contracts
```

2. **Run the script:**
```bash
go run get_value_simple.go
```

3. **Expected output:**
```
Contract value: 121212
```

## 🔧 Script Details

### get_value_simple.go

**Purpose:** Reads the current value from contract storage slot 0.

**How it works:**
1. Connects to Odyssey node RPC endpoint
2. Reads directly from storage slot 0 of contract `0x9`
3. Converts hex value to decimal
4. Displays the result

**Key features:**
- **Direct storage reading** (most reliable method)
- **No authentication required**
- **Works even when contract functions fail**
- **Fast and simple**

## 📋 Prerequisites

### Required Files
- `contracts/BasicContract.sol` - Solidity contract source
- `genesis/genesis_custom.json` - Custom genesis with contract
- `contracts/get_value_simple.go` - Value reader script

### Required Tools
- **Go 1.20+** - For running the script
- **Solidity compiler** - For generating ABI/bytecode
- **Odyssey node** - Must be running with custom genesis

### Network Configuration
- **RPC Endpoint:** `http://localhost:9650/ext/bc/m3Fo7ibkiC3311FGwE8qKZpv6EmA4KmwWCLb4etPEtPhRcpNN/rpc`
- **Contract Address:** `0x0000000000000000000000000000000000000009`
- **Storage Slot:** `0x0000000000000000000000000000000000000000000000000000000000000000`

## 🎯 Contract Information

### BasicContract.sol
```solidity
pragma solidity ^0.5.0;

contract BasicContract {
    uint256 public value;
    
    constructor() public {
        value = 1000;
    }
    
    function get() public view returns (uint256) {
        return value;
    }
    
    function set(uint256 _value) public {
        value = _value;
    }
}
```

### Genesis Configuration
- **Contract deployed at:** `0x9`
- **Initial value:** `121212` (overrides constructor value)
- **Storage slot 0:** Contains the `value` variable

## 🔍 Troubleshooting

### Script Not Working
1. **Check if node is running:**
```bash
curl -X POST "http://localhost:9650/ext/bc/m3Fo7ibkiC3311FGwE8qKZpv6EmA4KmwWCLb4etPEtPhRcpNN/rpc" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'
```

2. **Verify contract exists:**
```bash
curl -X POST "http://localhost:9650/ext/bc/m3Fo7ibkiC3311FGwE8qKZpv6EmA4KmwWCLb4etPEtPhRcpNN/rpc" \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_getStorageAt","params":["0x0000000000000000000000000000000000000009","0x0000000000000000000000000000000000000000000000000000000000000000","latest"],"id":1}'
```

### Node Issues
- **Ensure custom genesis is loaded**
- **Check node logs for errors**
- **Verify RPC endpoint is enabled**

### Contract Issues
- **Verify bytecode is correct**
- **Check storage initialization**
- **Ensure contract address matches**

## 📊 Expected Results

### Successful Run
```
Contract value: 121212
```

### Value Breakdown
- **Hex:** `0x000000000000000000000000000000000000000000000000000000000001d97c`
- **Decimal:** `121212`
- **Purpose:** Governance min validator stake parameter

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Go Script      │───▶│  Odyssey Node    │───▶│  Contract 0x9   │
│  (Reader)       │    │  (RPC Endpoint)  │    │  (Storage)      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

The script reads directly from contract storage, bypassing function calls for maximum reliability.

## 📝 Notes

- **Storage reading is the primary method** for this contract
- **Function calls may not work** due to VM limitations
- **The contract serves as a governance configuration store**
- **Value 121212 represents min validator stake parameter**
- **Direct storage access is more reliable than function calls**

---

**Last Updated:** September 5, 2025  
**Compatible with:** Odyssey v1.10.10+ with lucanali/coreth integration