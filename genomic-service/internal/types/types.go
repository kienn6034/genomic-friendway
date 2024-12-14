package types

// Shared types accross internal packages
type GeneData struct {
	ID            string
	EncryptedData []byte
	RiskScore     int
	FileHash      string
}

type ProcessResult struct {
	DocID       string
	RiskScore   int
	ContentHash string
	SessionID   string
}
