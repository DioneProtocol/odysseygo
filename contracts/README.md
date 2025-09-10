# Contract Value Reader

A Go script to read values from the BasicContract deployed at address `0x000000000000000000000000000000000000012a` in the Odyssey genesis.

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
   "dChainGenesis": "{\"config\":{\"chainId\":43112,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"5b9a1ef78f359e73b14142878704ac831e9500a8\":{\"balance\":\"0x295BE96E64066972000000\"},\"070a16d8c598fe8767b2c00ec7bb30dab52dd691\": {\"balance\": \"0x295BE96E64066972000000\"},\"000000000000000000000000000000000000012a\":{\"code\":\"0x6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fa4f24514604157806360fe47b114605b5780636d4ce63c14606c575b600080fd5b604960005481565b60405190815260200160405180910390f35b606a606636600460a3565b6073565b005b6000546049565b600081815560405182917fca9076777bea707d7602f398fc73932162e8a91d50202b88929c4cd710b6d6cd91a250565b60006020828403121560b457600080fd5b503591905056fea2646970667358221220839347288ae927517f23fab7d640a3e5a0352115658f078553db4ccee124399464736f6c63430008110033\",\"balance\":\"0x0\",\"storage\":{\"0x0000000000000000000000000000000000000000000000000000000000000000\":\"0x000000000000000000000000000000000000000000000000000000000001e240\"}}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}"

}
```

3. **Key components added:**
   - **Contract address:** `000000000000000000000000000000000000012a` 
   - **Bytecode:** The compiled contract code
   - **Storage:** Initial value `0x1e240` (123456 in decimal)

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
2. Reads directly from storage slot 0 of contract `0x000000000000000000000000000000000000012a`
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
- **Contract Address:** `0x000000000000000000000000000000000000012a`
- **Storage Slot:** `0x0000000000000000000000000000000000000000000000000000000000000000`
