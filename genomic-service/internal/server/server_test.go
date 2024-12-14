package server

import (
	"bytes"
	"encoding/json"
	"genomic-service/internal/config"
	"genomic-service/internal/types"
	teesdk "genomic-service/pkg/tee"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) *Server {
	config.LoadEnv("../../.env")
	cfg := config.NewConfig("../config/app.ini")
	server, err := NewServer(cfg)
	assert.NoError(t, err)
	return server
}

func TestServerEndpoints(t *testing.T) {
	testCases := []struct {
		name           string
		geneDataFile   string
		expectedScore  int
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Process Extremely High Risk Data",
			geneDataFile:   "alice.txt",
			expectedScore:  4,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Process High Risk Data",
			geneDataFile:   "bob.txt",
			expectedScore:  3,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Process Slightly High Risk Data",
			geneDataFile:   "charlie.txt",
			expectedScore:  2,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Process Low Risk Data",
			geneDataFile:   "dave.txt",
			expectedScore:  1,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := setupTestServer(t)

			// 1. Get TEE public key
			pubKey := getTEEPublicKey(t, server)

			// 2. Encrypt gene data
			encryptedData := encryptGeneData(t, pubKey, tc.geneDataFile)

			// 3. Upload encrypted data
			uploadResp := uploadData(t, server, encryptedData)

			if tc.expectError {
				assert.Empty(t, uploadResp)
				return
			}

			// 4. Confirm and process data
			result := confirmData(t, server, uploadResp["fileHash"], uploadResp["sessionId"])
			assert.Equal(t, tc.expectedScore, result.RiskScore)
		})
	}
}

func getTEEPublicKey(t *testing.T, server *Server) string {
	req, _ := http.NewRequest("GET", "/api/tee/public-key", nil)
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	publicKey := response["publicKey"]
	assert.NotEmpty(t, publicKey)
	return publicKey
}

func encryptGeneData(t *testing.T, publicKey, filename string) []byte {
	encoder := teesdk.NewTeeEncoder(publicKey)
	fileData, err := encoder.GetFileDataFromFile("../../gene-datas/" + filename)
	assert.NoError(t, err)

	encryptedData, err := encoder.EncryptGeneData(fileData.Data)
	assert.NoError(t, err)
	return encryptedData
}

func uploadData(t *testing.T, server *Server, encryptedData []byte) map[string]string {
	req, _ := http.NewRequest("POST", "/api/upload", bytes.NewBuffer(encryptedData))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		return nil
	}

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response
}

func confirmData(t *testing.T, server *Server, fileHash, sessionID string) *types.ProcessResult {
	reqBody := map[string]string{
		"fileHash":  fileHash,
		"sessionId": sessionID,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/api/confirm", bytes.NewBuffer(reqBodyJSON))
	resp := httptest.NewRecorder()
	server.router.ServeHTTP(resp, req)

	var response struct {
		Message string               `json:"message"`
		Result  *types.ProcessResult `json:"result"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	return response.Result
}
