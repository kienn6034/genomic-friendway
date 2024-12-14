package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"genomic-service/contracts"
	"genomic-service/internal/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainService struct {
	client     *ethclient.Client
	wallet     *Wallet
	controller *contracts.Controller
	nft        *contracts.GeneNFT
	token      *contracts.PCSPToken
}

func NewBlockchainService(rpcURL, privateKey, controllerAddr string) (*BlockchainService, error) {
	// Connect to network
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to network: %v", err)
	}

	// Setup wallet
	wallet, err := NewWallet(privateKey)
	if err != nil {
		return nil, err
	}

	// Load contracts
	controller, err := contracts.NewController(common.HexToAddress(controllerAddr), client)
	if err != nil {
		return nil, fmt.Errorf("failed to load controller: %v", err)
	}

	// Get NFT and Token addresses from controller
	nftAddr, err := controller.GeneNFT(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT address: %v", err)
	}

	tokenAddr, err := controller.PcspToken(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get token address: %v", err)
	}

	// Load NFT and Token contracts
	nft, err := contracts.NewGeneNFT(nftAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load NFT contract: %v", err)
	}

	token, err := contracts.NewPCSPToken(tokenAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load token contract: %v", err)
	}

	return &BlockchainService{
		client:     client,
		wallet:     wallet,
		controller: controller,
		nft:        nft,
		token:      token,
	}, nil
}

// InitiateDataUpload starts the upload session on blockchain
func (s *BlockchainService) InitiateDataUpload(docID string) (string, error) {
	opts, err := s.wallet.GetTransactOpts()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction opts: %v", err)
	}

	// Call uploadData on controller contract
	tx, err := s.controller.UploadData(opts, docID)
	if err != nil {
		return "", fmt.Errorf("failed to initiate upload: %v", err)
	}

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return "", fmt.Errorf("failed to wait for transaction: %v", err)
	}

	// Get session ID from event
	for _, log := range receipt.Logs {
		event, err := s.controller.ParseUploadData(*log)
		if err == nil && event != nil {
			return event.SessionId.String(), nil
		}
	}

	return "", fmt.Errorf("failed to get session ID from event")
}

// ProcessAndMint handles the confirmation, NFT minting, and token rewards
func (s *BlockchainService) ProcessAndMint(result *types.ProcessResult) error {
	opts, err := s.wallet.GetTransactOpts()
	if err != nil {
		return fmt.Errorf("failed to get transaction opts: %v", err)
	}

	// Convert session ID to big.Int
	sessionID := new(big.Int)
	sessionID.SetString(result.SessionID, 10)

	// Call confirm on controller contract
	tx, err := s.controller.Confirm(
		opts,
		result.DocID,
		result.ContentHash,
		"0x1234", // Simplified proof for development
		sessionID,
		big.NewInt(int64(result.RiskScore)),
	)
	if err != nil {
		return fmt.Errorf("failed to confirm upload: %v", err)
	}

	// Wait for transaction and process events
	receipt, err := bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction: %v", err)
	}

	// Process events
	for _, log := range receipt.Logs {
		// Check for NFT minted event
		if nftEvent, err := s.controller.ParseGeneNFTMinted(*log); err == nil && nftEvent != nil {
			fmt.Printf("NFT minted with token ID: %s\n", nftEvent.TokenId)
		}
		// Check for PCSP rewarded event
		if pcspEvent, err := s.controller.ParsePCSPRewarded(*log); err == nil && pcspEvent != nil {
			fmt.Printf("PCSP tokens rewarded: %s\n", pcspEvent.Amount)
		}
	}

	return nil
}
