package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"sermorpheus-engine-test/internal/config"
	"sermorpheus-engine-test/internal/models"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BlockchainService struct {
	db           *gorm.DB
	client       *ethclient.Client
	config       *config.Config
	usdtContract string
	usdtDecimals int
}

func NewBlockchainService(db *gorm.DB, cfg *config.Config) *BlockchainService {
	client, err := ethclient.Dial(cfg.BSCRPCUrl)
	if err != nil {
		log.Printf("Failed to connect to BSC testnet: %v", err)
		return &BlockchainService{
			db:           db,
			config:       cfg,
			usdtContract: cfg.USDTContract,
			usdtDecimals: cfg.USDTDecimals,
		}
	}

	return &BlockchainService{
		db:           db,
		client:       client,
		config:       cfg,
		usdtContract: cfg.USDTContract,
		usdtDecimals: cfg.USDTDecimals,
	}
}

func (bs *BlockchainService) GeneratePaymentAddress() (*models.PaymentAddress, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))

	paymentAddress := &models.PaymentAddress{
		Address:    address.Hex(),
		PrivateKey: privateKeyHex,
		IsUsed:     false,
	}

	if err := bs.db.Create(paymentAddress).Error; err != nil {
		return nil, fmt.Errorf("failed to save payment address: %w", err)
	}

	return paymentAddress, nil
}

func (bs *BlockchainService) GetAvailablePaymentAddress() (*models.PaymentAddress, error) {
	var paymentAddress models.PaymentAddress

	err := bs.db.Where("is_used = ?", false).First(&paymentAddress).Error
	if err != nil {
		return nil, fmt.Errorf("no available payment address: %w", err)
	}

	return &paymentAddress, nil
}

func (bs *BlockchainService) MarkAddressAsUsed(addressID string) error {
	return bs.db.Model(&models.PaymentAddress{}).
		Where("id = ?", addressID).
		Update("is_used", true).Error
}

func (bs *BlockchainService) CheckUSDTBalance(address string) (*big.Int, error) {
	if bs.client == nil {
		return nil, fmt.Errorf("blockchain client not available")
	}

	return big.NewInt(0), nil
}

func (bs *BlockchainService) MonitorPayment(transactionID uuid.UUID, expectedAmount float64, paymentAddress string) error {
	if bs.client == nil {
		log.Printf("Blockchain client not available, skipping monitoring for transaction %s", transactionID)
		return nil
	}

	go func() {
		log.Printf("Starting payment monitoring for transaction %s, expecting %.6f USDT to %s",
			transactionID, expectedAmount, paymentAddress)

		if bs.CheckRecentTransfer(transactionID, expectedAmount, paymentAddress) {
			return
		}

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		timeout := time.After(30 * time.Minute)

		for {
			select {
			case <-timeout:
				log.Printf("Payment monitoring timeout for transaction %s", transactionID)
				return
			case <-ticker.C:
				if bs.CheckRecentTransfer(transactionID, expectedAmount, paymentAddress) {
					return
				}
			}
		}
	}()

	return nil
}

func (bs *BlockchainService) CheckRecentTransfer(transactionID uuid.UUID, expectedAmount float64, paymentAddress string) bool {

	var existingTx models.Transaction
	err := bs.db.Where("id = ? AND status = ?", transactionID, "paid").First(&existingTx).Error
	if err == nil {
		log.Printf("Transaction %s already paid, skipping check", transactionID)
		return true
	}

	return bs.checkRecentTransactionsDirectly(transactionID, expectedAmount, paymentAddress)
}

func (bs *BlockchainService) checkRecentTransactionsDirectly(transactionID uuid.UUID, expectedAmount float64, paymentAddress string) bool {

	latestBlock, err := bs.client.BlockNumber(context.Background())
	if err != nil {
		log.Printf("Failed to get latest block: %v", err)
		return false
	}

	blocksToCheck := uint64(20)
	if latestBlock < blocksToCheck {
		blocksToCheck = latestBlock
	}

	log.Printf("Checking last %d blocks individually for address %s (expecting %.6f USDT)",
		blocksToCheck, paymentAddress, expectedAmount)

	for i := uint64(0); i < blocksToCheck; i++ {
		blockNum := latestBlock - i
		if bs.checkSingleBlockForTransfer(blockNum, transactionID, expectedAmount, paymentAddress) {
			return true
		}

		time.Sleep(100 * time.Millisecond)
	}

	return false
}

func (bs *BlockchainService) checkSingleBlockForTransfer(blockNum uint64, transactionID uuid.UUID, expectedAmount float64, paymentAddress string) bool {
	block, err := bs.client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
	if err != nil {
		log.Printf("Failed to get block %d: %v", blockNum, err)
		return false
	}

	for _, tx := range block.Transactions() {

		if tx.To() == nil || tx.To().Hex() != bs.usdtContract {
			continue
		}

		receipt, err := bs.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			continue
		}

		if bs.checkTransactionForPayment(receipt, transactionID, expectedAmount, paymentAddress) {
			return true
		}
	}

	return false
}

func (bs *BlockchainService) checkTransactionForPayment(receipt *types.Receipt, transactionID uuid.UUID, expectedAmount float64, paymentAddress string) bool {
	contractAddress := common.HexToAddress(bs.usdtContract)
	toAddress := common.HexToAddress(paymentAddress)
	transferSig := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

	for _, vLog := range receipt.Logs {

		if vLog.Address != contractAddress {
			continue
		}

		if len(vLog.Topics) < 3 || vLog.Topics[0] != transferSig {
			continue
		}

		logToAddress := common.BytesToAddress(vLog.Topics[2].Bytes())
		if logToAddress != toAddress {
			continue
		}

		if bs.processTransferLog(*vLog, transactionID, expectedAmount, paymentAddress) {
			return true
		}
	}

	return false
}

func (bs *BlockchainService) processTransferLog(vLog types.Log, transactionID uuid.UUID, expectedAmount float64, paymentAddress string) bool {

	if len(vLog.Data) < 32 {
		return false
	}

	amountBig := new(big.Int).SetBytes(vLog.Data[len(vLog.Data)-32:])
	amountFloat := new(big.Float).SetInt(amountBig)

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(bs.usdtDecimals)), nil))
	amountUSDT, _ := new(big.Float).Quo(amountFloat, divisor).Float64()

	log.Printf("Found transfer: %.6f USDT (tx: %s)", amountUSDT, vLog.TxHash.Hex())

	tolerance := math.Max(expectedAmount*0.001, 0.000001)
	if amountUSDT >= expectedAmount-tolerance && amountUSDT <= expectedAmount+tolerance {
		log.Printf("Payment match found! Amount: %.6f USDT, Expected: %.6f USDT, Processing transaction %s",
			amountUSDT, expectedAmount, transactionID)

		err := bs.confirmTransactionPayment(transactionID, vLog.TxHash.Hex(), amountUSDT)
		if err != nil {
			log.Printf("Failed to confirm transaction: %v", err)
			return false
		}
		return true
	}

	return false
}

func (bs *BlockchainService) handleTransferEvent(vLog types.Log, transactionID uuid.UUID, expectedAmount float64, paymentAddress string) {
	log.Printf("Transfer event detected for transaction %s", transactionID)

	if len(vLog.Data) < 32 {
		log.Printf("Invalid log data length: %d", len(vLog.Data))
		return
	}

	amountBig := new(big.Int).SetBytes(vLog.Data[len(vLog.Data)-32:])
	amountFloat := new(big.Float).SetInt(amountBig)

	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(bs.usdtDecimals)), nil))
	amountUSDT, _ := new(big.Float).Quo(amountFloat, divisor).Float64()

	log.Printf("Received %.6f USDT, expected %.6f USDT", amountUSDT, expectedAmount)

	tolerance := 0.000001
	if amountUSDT >= expectedAmount-tolerance {
		log.Printf("Payment confirmed! Processing transaction %s", transactionID)

		err := bs.confirmTransactionPayment(transactionID, vLog.TxHash.Hex(), amountUSDT)
		if err != nil {
			log.Printf("Failed to confirm transaction: %v", err)
		}
	} else {
		log.Printf("Amount mismatch: received %.6f, expected %.6f", amountUSDT, expectedAmount)
	}
}

func (bs *BlockchainService) confirmTransactionPayment(transactionID uuid.UUID, txHash string, amount float64) error {
	return bs.db.Transaction(func(tx *gorm.DB) error {

		var existingTx models.Transaction
		err := tx.Where("id = ?", transactionID).First(&existingTx).Error
		if err != nil {
			return fmt.Errorf("transaction not found: %w", err)
		}

		if existingTx.Status == "paid" {
			log.Printf("Transaction %s already paid, skipping confirmation", transactionID)
			return nil
		}

		now := time.Now()
		err = tx.Model(&models.Transaction{}).
			Where("id = ? AND status = ?", transactionID, "pending").
			Updates(map[string]interface{}{
				"status":               "paid",
				"payment_confirmed_at": &now,
				"updated_at":           now,
			}).Error
		if err != nil {
			return fmt.Errorf("failed to update transaction status: %w", err)
		}

		var existingBlockchainTx models.BlockchainTransaction
		err = tx.Where("transaction_id = ? AND tx_hash = ?", transactionID, txHash).First(&existingBlockchainTx).Error
		if err == nil {
			log.Printf("Blockchain transaction record already exists for tx %s", txHash)
			return nil
		}

		blockchainTx := &models.BlockchainTransaction{
			TransactionID: transactionID,
			TxHash:        txHash,
			ToAddress:     existingTx.PaymentAddress,
			Amount:        amount,
			Status:        "confirmed",
		}
		if err := tx.Create(blockchainTx).Error; err != nil {
			return fmt.Errorf("failed to create blockchain transaction: %w", err)
		}

		log.Printf("Transaction %s confirmed with tx hash %s, amount: %.6f USDT", transactionID, txHash, amount)
		return nil
	})
}
