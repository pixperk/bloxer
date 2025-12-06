# Bloxer

An educational blockchain implementation in Go. Learn the fundamentals of blockchain technology through a simple, interactive CLI.

## What You'll Learn

- **Key Generation**: ECDSA (Elliptic Curve Digital Signature Algorithm) for secure key pairs
- **Digital Signatures**: How transactions are signed and verified
- **Proof of Work**: Mining blocks with adjustable difficulty
- **Chain Validation**: Ensuring blockchain integrity

## Installation

```bash
git clone https://github.com/pixperk/bloxer.git
cd bloxer
go build -o bloxer .
```

## Quick Start

```bash
# 1. Create a wallet
./bloxer wallet create

# 2. Mine some blocks to earn coins (mine twice to collect reward)
./bloxer mine
./bloxer mine

# 3. Check your balance
./bloxer balance

# 4. Send coins to someone
./bloxer send --to <recipient-address> --amount 10

# 5. Mine to confirm the transaction
./bloxer mine

# 6. View the blockchain
./bloxer chain
```

## CLI Reference

### Wallet Management

```bash
bloxer wallet create    # Create a new wallet (generates ECDSA key pair)
bloxer wallet show      # Display your wallet address
bloxer wallet delete    # Delete your wallet (irreversible)
```

### Transactions

```bash
bloxer send --to <address> --amount <coins>   # Create a transaction
bloxer send -t <address> -a <coins>           # Short form
```

### Mining

```bash
bloxer mine             # Mine pending transactions into a new block
```

Mining does two things:
1. Packages pending transactions into a block
2. Awards you a mining reward (added to next block's pending transactions)

**Note**: Mining rewards are collected in the *next* block you mine, not immediately.

### Viewing Data

```bash
bloxer balance              # Check your wallet balance
bloxer balance <address>    # Check any address balance
bloxer chain                # View all blocks in the chain
bloxer validate             # Verify blockchain integrity
```

### Reset

```bash
bloxer reset          # Reset blockchain only (keeps wallet)
bloxer reset --all    # Reset blockchain and delete wallet
```

## How It Works

### Architecture

```
~/.bloxer/
  ├── wallet.json       # Your private key and address
  └── blockchain.json   # The entire blockchain state
```

### Key Generation

Bloxer uses ECDSA with the P-256 curve for cryptographic operations:

```
Private Key (random 256-bit number)
       │
       ▼
Public Key (point on elliptic curve)
       │
       ▼
Address (hex-encoded public key bytes)
```

### Transaction Flow

```
1. Create Transaction
   ┌─────────────────────────────────┐
   │ From: your_address              │
   │ To: recipient_address           │
   │ Amount: 10.00                   │
   │ Signature: (empty)              │
   └─────────────────────────────────┘
                 │
                 ▼
2. Sign Transaction (with your private key)
   ┌─────────────────────────────────┐
   │ From: your_address              │
   │ To: recipient_address           │
   │ Amount: 10.00                   │
   │ Signature: 0x3045...            │
   └─────────────────────────────────┘
                 │
                 ▼
3. Add to Pending Transactions
                 │
                 ▼
4. Mine Block (transactions included)
```

### Block Structure

```
┌─────────────────────────────────────┐
│ Block                               │
├─────────────────────────────────────┤
│ Hash:      00a3f2...  (starts with  │
│                        leading 0s)  │
│ PrevHash:  7b2c91...                │
│ Timestamp: 1701892345               │
│ Nonce:     42851                    │
│ Data:                               │
│   └── Transactions: [...]           │
└─────────────────────────────────────┘
```

### Proof of Work

Mining requires finding a `nonce` such that the block's hash starts with N zeros (where N = difficulty):

```
Difficulty: 2
Valid hash must start with: "00..."

Attempt 1: nonce=0     hash=a3f2b1... (invalid)
Attempt 2: nonce=1     hash=8c4e22... (invalid)
...
Attempt N: nonce=42851 hash=00a3f2... (valid!)
```

### Chain Validation

The blockchain is valid if:
1. Each block's hash matches its calculated hash
2. Each block's `prevHash` matches the previous block's hash
3. All transactions have valid signatures

```
┌─────────┐    ┌─────────┐    ┌─────────┐
│ Block 0 │───▶│ Block 1 │───▶│ Block 2 │
│ Genesis │    │         │    │         │
│ prev: 0 │    │ prev:   │    │ prev:   │
│         │    │ hash(0) │    │ hash(1) │
└─────────┘    └─────────┘    └─────────┘
```

## Configuration

Default settings (hardcoded for simplicity):
- **Difficulty**: 2 (hash must start with "00")
- **Mining Reward**: 100 coins
- **Data Directory**: `~/.bloxer/`

## Example Session

```bash
$ ./bloxer wallet create

[OK] Wallet created successfully!

  Your address:
  04a1b2c3d4e5f6...

  Keep your wallet file safe!
  Location: /home/user/.bloxer/wallet.json

$ ./bloxer mine

Mining block...

  Difficulty: 2
  Pending transactions: 0

Block mined: 00f3a2b1c4...
Block successfully mined!

[OK] Block mined successfully!

  Time taken: 12ms
  Reward: 100.00 coins
  New balance: 0.00 coins

$ ./bloxer mine

Mining block...

  Difficulty: 2
  Pending transactions: 1

Block mined: 00b2c3d4e5...
Block successfully mined!

[OK] Block mined successfully!

  Time taken: 8ms
  Reward: 100.00 coins
  New balance: 100.00 coins

$ ./bloxer balance

Balance

  Address: 04a1b2c3d4...c3d4e5f6
  Balance: 100.00 coins
```

## Limitations

This is an educational implementation. It does not include:
- Networking/P2P communication
- Merkle trees
- UTXO model
- Consensus mechanisms beyond PoW
- Balance validation (you can spend coins you don't have)
- Persistent mempool

## License

MIT
