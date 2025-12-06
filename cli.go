package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const (
	dataDir        = ".bloxer"
	blockchainFile = "blockchain.json"
	walletFile     = "wallet.json"
)

// Persistence types
type WalletData struct {
	PrivateKey []byte `json:"private_key"`
	Address    string `json:"address"`
}

type TransactionData struct {
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Amount      float64 `json:"amount"`
	Signature   []byte  `json:"signature"`
}

type BlockData struct {
	Data      map[string]interface{} `json:"data"`
	PrevHash  string                 `json:"prev_hash"`
	TimeStamp int64                  `json:"timestamp"`
	Hash      string                 `json:"hash"`
	Nonce     int                    `json:"nonce"`
}

type BlockchainData struct {
	Chain               []BlockData       `json:"chain"`
	Difficulty          int               `json:"difficulty"`
	PendingTransactions []TransactionData `json:"pending_transactions"`
	MiningReward        float64           `json:"mining_reward"`
}

// CLI colors and formatting
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func getDataDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, dataDir)
}

func ensureDataDir() error {
	return os.MkdirAll(getDataDir(), 0755)
}

// Wallet persistence
func saveWallet(privateKey *ecdsa.PrivateKey, address string) error {
	if err := ensureDataDir(); err != nil {
		return err
	}
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	wallet := WalletData{PrivateKey: keyBytes, Address: address}
	data, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(getDataDir(), walletFile), data, 0600)
}

func loadWallet() (*ecdsa.PrivateKey, string, error) {
	data, err := os.ReadFile(filepath.Join(getDataDir(), walletFile))
	if err != nil {
		return nil, "", err
	}
	var wallet WalletData
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, "", err
	}
	privateKey, err := x509.ParseECPrivateKey(wallet.PrivateKey)
	if err != nil {
		return nil, "", err
	}
	return privateKey, wallet.Address, nil
}

func walletExists() bool {
	_, err := os.Stat(filepath.Join(getDataDir(), walletFile))
	return err == nil
}

// Blockchain persistence helpers
func transactionsToData(txs []Transaction) []TransactionData {
	result := make([]TransactionData, len(txs))
	for i, tx := range txs {
		result[i] = TransactionData{
			FromAddress: tx.FromAddress,
			ToAddress:   tx.ToAddress,
			Amount:      tx.Amount,
			Signature:   tx.Signature,
		}
	}
	return result
}

func dataToTransactions(data []TransactionData) []Transaction {
	result := make([]Transaction, len(data))
	for i, td := range data {
		result[i] = Transaction{
			FromAddress: td.FromAddress,
			ToAddress:   td.ToAddress,
			Amount:      td.Amount,
			Signature:   td.Signature,
		}
	}
	return result
}

func saveBlockchain(bc *Blockchain) error {
	if err := ensureDataDir(); err != nil {
		return err
	}

	chainData := make([]BlockData, len(bc.Chain))
	for i, block := range bc.Chain {
		blockDataMap := make(map[string]interface{})
		for k, v := range block.Data {
			if k == "transactions" {
				if txs, ok := v.([]Transaction); ok {
					blockDataMap[k] = transactionsToData(txs)
				} else {
					blockDataMap[k] = v
				}
			} else {
				blockDataMap[k] = v
			}
		}
		chainData[i] = BlockData{
			Data:      blockDataMap,
			PrevHash:  block.PrevHash,
			TimeStamp: block.TimeStamp,
			Hash:      block.Hash,
			Nonce:     block.Nonce,
		}
	}

	bcData := BlockchainData{
		Chain:               chainData,
		Difficulty:          bc.Difficulty,
		PendingTransactions: transactionsToData(bc.PendingTransactions),
		MiningReward:        bc.MiningReward,
	}

	data, err := json.MarshalIndent(bcData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(getDataDir(), blockchainFile), data, 0644)
}

func loadBlockchain() (*Blockchain, error) {
	data, err := os.ReadFile(filepath.Join(getDataDir(), blockchainFile))
	if err != nil {
		return nil, err
	}

	var bcData BlockchainData
	if err := json.Unmarshal(data, &bcData); err != nil {
		return nil, err
	}

	chain := make([]Block, len(bcData.Chain))
	for i, bd := range bcData.Chain {
		blockDataMap := make(map[string]interface{})
		for k, v := range bd.Data {
			if k == "transactions" {
				if txDataRaw, ok := v.([]interface{}); ok {
					txs := make([]Transaction, len(txDataRaw))
					for j, txRaw := range txDataRaw {
						if txMap, ok := txRaw.(map[string]interface{}); ok {
							txs[j] = Transaction{
								FromAddress: getJSONString(txMap, "from_address"),
								ToAddress:   getJSONString(txMap, "to_address"),
								Amount:      getJSONFloat(txMap, "amount"),
								Signature:   getJSONBytes(txMap, "signature"),
							}
						}
					}
					blockDataMap[k] = txs
				}
			} else {
				blockDataMap[k] = v
			}
		}
		chain[i] = Block{
			Data:      blockDataMap,
			PrevHash:  bd.PrevHash,
			TimeStamp: bd.TimeStamp,
			Hash:      bd.Hash,
			Nonce:     bd.Nonce,
		}
	}

	return &Blockchain{
		Chain:               chain,
		Difficulty:          bcData.Difficulty,
		PendingTransactions: dataToTransactions(bcData.PendingTransactions),
		MiningReward:        bcData.MiningReward,
	}, nil
}

func blockchainExists() bool {
	_, err := os.Stat(filepath.Join(getDataDir(), blockchainFile))
	return err == nil
}

func getOrCreateBlockchain() *Blockchain {
	if blockchainExists() {
		bc, err := loadBlockchain()
		if err == nil {
			return bc
		}
	}
	bc := NewBlockchain(2, 100.0)
	saveBlockchain(bc)
	return bc
}

// JSON helpers
func getJSONString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getJSONFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func getJSONBytes(m map[string]interface{}, key string) []byte {
	if v, ok := m[key]; ok {
		// JSON encodes []byte as base64 string
		if s, ok := v.(string); ok {
			decoded, err := base64.StdEncoding.DecodeString(s)
			if err == nil {
				return decoded
			}
		}
		// Fallback for array of numbers
		if arr, ok := v.([]interface{}); ok {
			bytes := make([]byte, len(arr))
			for i, b := range arr {
				if f, ok := b.(float64); ok {
					bytes[i] = byte(f)
				}
			}
			return bytes
		}
	}
	return nil
}

// Formatting helpers
func formatAddress(addr string) string {
	if len(addr) > 20 {
		return addr[:10] + "..." + addr[len(addr)-10:]
	}
	return addr
}

func publicKeyToAddress(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)
	return fmt.Sprintf("%x", pubKeyBytes)
}

// Root command
var rootCmd = &cobra.Command{
	Use:   "bloxer",
	Short: "Bloxer - An educational blockchain CLI",
	Long: colorCyan + colorBold + `
  ╔══════════════════════════════════════════════════════════════╗
  ║                                                              ║
  ║   ██████╗ ██╗      ██████╗ ██╗  ██╗███████╗██████╗           ║
  ║   ██╔══██╗██║     ██╔═══██╗╚██╗██╔╝██╔════╝██╔══██╗          ║
  ║   ██████╔╝██║     ██║   ██║ ╚███╔╝ █████╗  ██████╔╝          ║
  ║   ██╔══██╗██║     ██║   ██║ ██╔██╗ ██╔══╝  ██╔══██╗          ║
  ║   ██████╔╝███████╗╚██████╔╝██╔╝ ██╗███████╗██║  ██║          ║
  ║   ╚═════╝ ╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝          ║
  ║                                                              ║
  ╚══════════════════════════════════════════════════════════════╝` + colorReset + `

  An educational blockchain implementation in Go.

  ` + colorYellow + `Learn about:` + colorReset + `
    - Key generation and digital signatures (ECDSA)
    - Transaction creation and validation
    - Block mining with proof-of-work
    - Chain validation and integrity

  ` + colorGreen + `Get started:` + colorReset + `
    bloxer wallet create    Create a new wallet
    bloxer mine             Mine pending transactions
    bloxer send             Send coins to an address
    bloxer balance          Check your balance
    bloxer chain            View the blockchain
`,
}

// Wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manage your wallet",
	Long:  "Create and manage your blockchain wallet",
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		if walletExists() {
			fmt.Printf("%s%s[ERROR] Wallet already exists!%s\n", colorRed, colorBold, colorReset)
			fmt.Printf("  Use %sbloxer wallet show%s to see your address\n", colorCyan, colorReset)
			return
		}

		privateKey, publicKey, err := GenerateKeyPair()
		if err != nil {
			fmt.Printf("%s[ERROR] Error generating key pair: %v%s\n", colorRed, err, colorReset)
			return
		}

		address := publicKeyToAddress(publicKey)
		if err := saveWallet(privateKey, address); err != nil {
			fmt.Printf("%s[ERROR] Error saving wallet: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%s[OK] Wallet created successfully!%s\n\n", colorGreen, colorBold, colorReset)
		fmt.Printf("  %sYour address:%s\n", colorYellow, colorReset)
		fmt.Printf("  %s%s%s\n\n", colorCyan, address, colorReset)
		fmt.Printf("  %sKeep your wallet file safe!%s\n", colorYellow, colorReset)
		fmt.Printf("  Location: %s\n\n", filepath.Join(getDataDir(), walletFile))
	},
}

var walletShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show wallet address",
	Run: func(cmd *cobra.Command, args []string) {
		if !walletExists() {
			fmt.Printf("%s[ERROR] No wallet found. Create one with: bloxer wallet create%s\n", colorRed, colorReset)
			return
		}

		_, address, err := loadWallet()
		if err != nil {
			fmt.Printf("%s[ERROR] Error loading wallet: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%sYour Wallet%s\n\n", colorCyan, colorBold, colorReset)
		fmt.Printf("  %sAddress:%s\n", colorYellow, colorReset)
		fmt.Printf("  %s\n\n", address)
	},
}

// Balance command
var balanceCmd = &cobra.Command{
	Use:   "balance [address]",
	Short: "Check balance of an address",
	Long:  "Check the balance of your wallet or any address",
	Run: func(cmd *cobra.Command, args []string) {
		bc := getOrCreateBlockchain()

		var address string
		if len(args) > 0 {
			address = args[0]
		} else {
			if !walletExists() {
				fmt.Printf("%s[ERROR] No wallet found. Create one with: bloxer wallet create%s\n", colorRed, colorReset)
				return
			}
			_, address, _ = loadWallet()
		}

		balance := bc.GetBalanceOfAddress(address)

		fmt.Printf("\n%s%sBalance%s\n\n", colorCyan, colorBold, colorReset)
		fmt.Printf("  %sAddress:%s %s\n", colorYellow, colorReset, formatAddress(address))
		fmt.Printf("  %sBalance:%s %s%.2f coins%s\n\n", colorYellow, colorReset, colorGreen, balance, colorReset)
	},
}

// Send command
var sendAmount float64
var sendTo string

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send coins to an address",
	Long:  "Create and sign a transaction to send coins",
	Run: func(cmd *cobra.Command, args []string) {
		if !walletExists() {
			fmt.Printf("%s[ERROR] No wallet found. Create one with: bloxer wallet create%s\n", colorRed, colorReset)
			return
		}

		if sendTo == "" {
			fmt.Printf("%s[ERROR] Please specify recipient with --to flag%s\n", colorRed, colorReset)
			return
		}

		if sendAmount <= 0 {
			fmt.Printf("%s[ERROR] Please specify a positive amount with --amount flag%s\n", colorRed, colorReset)
			return
		}

		privateKey, address, err := loadWallet()
		if err != nil {
			fmt.Printf("%s[ERROR] Error loading wallet: %v%s\n", colorRed, err, colorReset)
			return
		}

		bc := getOrCreateBlockchain()

		tx := NewTransaction(address, sendTo, sendAmount)
		tx.signTransaction(privateKey)

		if err := bc.AddTransaction(tx); err != nil {
			fmt.Printf("%s[ERROR] Transaction failed: %v%s\n", colorRed, err, colorReset)
			return
		}

		if err := saveBlockchain(bc); err != nil {
			fmt.Printf("%s[ERROR] Error saving blockchain: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%s[OK] Transaction created!%s\n\n", colorGreen, colorBold, colorReset)
		fmt.Printf("  %sFrom:%s    %s\n", colorYellow, colorReset, formatAddress(address))
		fmt.Printf("  %sTo:%s      %s\n", colorYellow, colorReset, formatAddress(sendTo))
		fmt.Printf("  %sAmount:%s  %.2f coins\n\n", colorYellow, colorReset, sendAmount)
		fmt.Printf("  %sTransaction is pending. Run %sbloxer mine%s to include it in a block.%s\n\n", colorPurple, colorCyan, colorPurple, colorReset)
	},
}

// Mine command
var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "Mine pending transactions",
	Long:  "Mine a new block with pending transactions and receive a reward",
	Run: func(cmd *cobra.Command, args []string) {
		if !walletExists() {
			fmt.Printf("%s[ERROR] No wallet found. Create one with: bloxer wallet create%s\n", colorRed, colorReset)
			return
		}

		_, address, err := loadWallet()
		if err != nil {
			fmt.Printf("%s[ERROR] Error loading wallet: %v%s\n", colorRed, err, colorReset)
			return
		}

		bc := getOrCreateBlockchain()

		fmt.Printf("\n%s%sMining block...%s\n\n", colorYellow, colorBold, colorReset)
		fmt.Printf("  Difficulty: %d\n", bc.Difficulty)
		fmt.Printf("  Pending transactions: %d\n\n", len(bc.PendingTransactions))

		startTime := time.Now()
		bc.MinePendingTransactions(address)
		duration := time.Since(startTime)

		if err := saveBlockchain(bc); err != nil {
			fmt.Printf("%s[ERROR] Error saving blockchain: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%s[OK] Block mined successfully!%s\n\n", colorGreen, colorBold, colorReset)
		fmt.Printf("  %sTime taken:%s %v\n", colorYellow, colorReset, duration.Round(time.Millisecond))
		fmt.Printf("  %sReward:%s %.2f coins\n", colorYellow, colorReset, bc.MiningReward)
		fmt.Printf("  %sNew balance:%s %.2f coins\n\n", colorYellow, colorReset, bc.GetBalanceOfAddress(address))
	},
}

// Chain command
var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "View the blockchain",
	Long:  "Display all blocks in the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := getOrCreateBlockchain()

		fmt.Printf("\n%s%sBlockchain%s\n", colorCyan, colorBold, colorReset)
		fmt.Printf("  Total blocks: %d\n\n", len(bc.Chain))

		for i, block := range bc.Chain {
			fmt.Printf("  %s┌─ Block #%d ─────────────────────────────────────┐%s\n", colorBlue, i, colorReset)
			fmt.Printf("  │ %sHash:%s      %s\n", colorYellow, colorReset, formatAddress(block.Hash))
			fmt.Printf("  │ %sPrev:%s      %s\n", colorYellow, colorReset, formatAddress(block.PrevHash))
			fmt.Printf("  │ %sTimestamp:%s %s\n", colorYellow, colorReset, time.Unix(block.TimeStamp, 0).Format("2006-01-02 15:04:05"))
			fmt.Printf("  │ %sNonce:%s     %d\n", colorYellow, colorReset, block.Nonce)

			if txs, ok := block.Data["transactions"].([]Transaction); ok && len(txs) > 0 {
				fmt.Printf("  │ %sTransactions:%s\n", colorYellow, colorReset)
				for _, tx := range txs {
					from := formatAddress(tx.FromAddress)
					if tx.FromAddress == "" {
						from = colorGreen + "MINING REWARD" + colorReset
					}
					fmt.Printf("  │   %s → %s: %.2f\n", from, formatAddress(tx.ToAddress), tx.Amount)
				}
			}
			fmt.Printf("  %s└────────────────────────────────────────────────┘%s\n\n", colorBlue, colorReset)
		}

		if len(bc.PendingTransactions) > 0 {
			fmt.Printf("  %s%sPending Transactions: %d%s\n\n", colorYellow, colorBold, len(bc.PendingTransactions), colorReset)
		}
	},
}

// Validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the blockchain",
	Long:  "Check if the blockchain is valid and hasn't been tampered with",
	Run: func(cmd *cobra.Command, args []string) {
		bc := getOrCreateBlockchain()

		fmt.Printf("\n%s%sValidating blockchain...%s\n\n", colorCyan, colorBold, colorReset)

		if bc.IsChainValid() {
			fmt.Printf("  %s%s[OK] Blockchain is valid!%s\n\n", colorGreen, colorBold, colorReset)
		} else {
			fmt.Printf("  %s%s[ERROR] Blockchain is INVALID!%s\n\n", colorRed, colorBold, colorReset)
			fmt.Printf("  The chain may have been tampered with.\n\n")
		}
	},
}

// Reset command
var resetAll bool

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the blockchain",
	Long:  "Delete blockchain data and start fresh. Use --all to also delete wallet.",
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := getDataDir()
		bcPath := filepath.Join(dataPath, blockchainFile)

		if err := os.Remove(bcPath); err != nil && !os.IsNotExist(err) {
			fmt.Printf("%s[ERROR] Error resetting blockchain: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%s[OK] Blockchain reset successfully!%s\n\n", colorGreen, colorBold, colorReset)

		if resetAll {
			wPath := filepath.Join(dataPath, walletFile)
			if err := os.Remove(wPath); err != nil && !os.IsNotExist(err) {
				fmt.Printf("%s[ERROR] Error deleting wallet: %v%s\n", colorRed, err, colorReset)
				return
			}
			fmt.Printf("%s%s[OK] Wallet deleted!%s\n\n", colorGreen, colorBold, colorReset)
		}

		fmt.Printf("  A new genesis block will be created on next operation.\n\n")
	},
}

// Wallet delete command
var walletDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete your wallet",
	Long:  "Permanently delete your wallet. This action cannot be undone!",
	Run: func(cmd *cobra.Command, args []string) {
		if !walletExists() {
			fmt.Printf("%s[ERROR] No wallet found.%s\n", colorRed, colorReset)
			return
		}

		wPath := filepath.Join(getDataDir(), walletFile)
		if err := os.Remove(wPath); err != nil {
			fmt.Printf("%s[ERROR] Error deleting wallet: %v%s\n", colorRed, err, colorReset)
			return
		}

		fmt.Printf("\n%s%s[OK] Wallet deleted!%s\n\n", colorGreen, colorBold, colorReset)
	},
}

func initCLI() {
	// Wallet subcommands
	walletCmd.AddCommand(walletCreateCmd)
	walletCmd.AddCommand(walletShowCmd)
	walletCmd.AddCommand(walletDeleteCmd)

	// Send flags
	sendCmd.Flags().Float64VarP(&sendAmount, "amount", "a", 0, "Amount to send")
	sendCmd.Flags().StringVarP(&sendTo, "to", "t", "", "Recipient address")

	// Reset flags
	resetCmd.Flags().BoolVarP(&resetAll, "all", "a", false, "Also delete wallet")

	// Add all commands to root
	rootCmd.AddCommand(walletCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(chainCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(resetCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func runCLI() {
	initCLI()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
