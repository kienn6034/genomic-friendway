package tee

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

// Client sdk for user to encrypt their data
type TeeEncoder struct {
	teePublicKeyHex string
}

func NewTeeEncoder(teePublicKeyHex string) *TeeEncoder {
	return &TeeEncoder{
		teePublicKeyHex: teePublicKeyHex,
	}
}

func (u *TeeEncoder) EncryptGeneData(data []byte) ([]byte, error) {
	// Decode hex string to bytes
	pubKeyBytes, err := hex.DecodeString(u.teePublicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key hex: %v", err)
	}

	// Parse public key from bytes
	pubKey, err := crypto.DecompressPubkey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress public key: %v", err)
	}

	// Convert to ECIES public key
	eciesPubKey := ecies.ImportECDSAPublic(pubKey)

	// Encrypt the data using crypto/rand.Reader directly
	encrypted, err := ecies.Encrypt(
		rand.Reader, // Use crypto/rand.Reader instead of bytes.NewReader
		eciesPubKey,
		data,
		nil,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	return encrypted, nil
}

func (u *TeeEncoder) GetFileDataFromFile(filePath string) (*FileData, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("gene data file not found: %v", err)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gene data file: %v", err)
	}

	// Check if data is empty
	if len(data) == 0 {
		return nil, fmt.Errorf("gene data file is empty")
	}

	return &FileData{
		FileHash: filePath,
		Data:     data,
	}, nil
}
