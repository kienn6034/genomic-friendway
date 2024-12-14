package tee_test

import (
	"genomic-service/internal/tee"
	teesdk "genomic-service/pkg/tee"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTEE(t *testing.T) {
	// Create new TEE instance
	tee := tee.NewTEE()
	assert.NotNil(t, tee)

	// Create user with TEE's public key
	user := teesdk.NewTeeEncoder(tee.GetPublicKey())
	assert.NotNil(t, user)

	// Test cases for different gene data files
	testCases := []struct {
		filename      string
		expectedScore int
		expectedError bool
	}{
		{"invalid.txt", 0, true},  // "invalid" = 0
		{"alice.txt", 4, false},   // "extremely high risk" = 4
		{"bob.txt", 3, false},     // "high risk" = 3
		{"charlie.txt", 2, false}, // "slightly high risk" = 2
		{"dave.txt", 1, false},    // "low risk" = 1
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			// Get gene data from file
			fileData, err := user.GetFileDataFromFile("../../gene-datas/" + tc.filename)
			assert.NoError(t, err)
			assert.NotNil(t, fileData)

			// Encrypt data using user's SDK
			encrypted, err := user.EncryptGeneData(fileData.Data)
			assert.NoError(t, err)
			assert.NotNil(t, encrypted)

			// Process (decrypt and calculate risk score) the encrypted data using TEE
			riskScore, err := tee.ProcessEncryptedData(encrypted, fileData.FileHash)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedScore, riskScore.RiskScore)
			}
		})
	}
}
