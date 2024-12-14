package tee

import (
	"fmt"
	"genomic-service/internal/storage"
	"genomic-service/internal/types"
)

type TEEService struct {
	tee     *TEE
	storage storage.Storage
}

func NewTEEService(storage storage.Storage) *TEEService {
	return &TEEService{
		tee:     NewTEE(),
		storage: storage,
	}
}

func (s *TEEService) ProcessGeneData(fileHash string) (*types.ProcessResult, error) {
	encryptedData, err := s.storage.Retrieve(fileHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data: %v", err)
	}

	// Process in TEE
	geneData, err := s.tee.ProcessEncryptedData(encryptedData, fileHash)
	if err != nil {
		return nil, fmt.Errorf("failed to process data: %v", err)
	}

	return &types.ProcessResult{
		DocID:     fileHash,
		RiskScore: geneData.RiskScore,
	}, nil
}

func (s *TEEService) GetTEEPublicKey() string {
	return s.tee.GetPublicKey()
}
