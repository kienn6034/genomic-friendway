package tee

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"genomic-service/internal/types"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

type TEE struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewTEE() *TEE {
	// Generate ECDSA key pair (same curve as Ethereum)
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	return &TEE{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
}

// GetPublicKey returns hex-encoded public key that users will use
func (t *TEE) GetPublicKey() string {
	pubKeyBytes := crypto.CompressPubkey(t.publicKey)
	return hex.EncodeToString(pubKeyBytes)
}

// ProcessEncryptedData decrypts and processes the data
func (t *TEE) ProcessEncryptedData(encryptedData []byte, fileHash string) (types.GeneData, error) {
	// Convert ECDSA private key to ECIES private key
	eciesPrivKey := ecies.ImportECDSA(t.privateKey)

	// Decrypt the data
	decrypted, err := eciesPrivKey.Decrypt(encryptedData, nil, nil)
	if err != nil {
		return types.GeneData{}, err
	}

	riskScore := t.calculateRiskScore(string(decrypted))
	if riskScore == 0 {
		return types.GeneData{}, fmt.Errorf("invalid risk score")
	}

	return types.GeneData{
		ID:            fileHash,
		FileHash:      fileHash,
		EncryptedData: encryptedData,
		RiskScore:     riskScore,
	}, nil
}

func (t *TEE) calculateRiskScore(data string) int {
	switch strings.ToLower(data) {
	case "extremely high risk":
		return 4 // 15,000 PCSP
	case "high risk":
		return 3 // 3,000 PCSP
	case "slightly high risk":
		return 2 // 225 PCSP
	case "low risk":
		return 1 // 30 PCSP
	default:
		return 0 // invalid
	}
}
