package server

import (
	"genomic-service/internal/blockchain"
	"genomic-service/internal/config"
	"genomic-service/internal/storage"
	"genomic-service/internal/tee"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	storage    storage.Storage
	tee        *tee.TEEService
	blockchain *blockchain.BlockchainService
}

func NewServer(cfg *config.Config) (*Server, error) {
	router := gin.Default()

	// Initialize storage
	storage := storage.NewMemoryStorage()

	// Initialize TEE service
	teeService := tee.NewTEEService(storage)

	// Initialize blockchain service
	blockchainService, err := blockchain.NewBlockchainService(
		cfg.BlockchainSettings.RPCURL,
		cfg.WalletSettings.PrivateKey,
		cfg.BlockchainSettings.ControllerAddress,
	)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		router:     router,
		storage:    storage,
		tee:        teeService,
		blockchain: blockchainService,
	}

	srv.setupRoutes()
	return srv, nil
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		api.POST("/upload", s.handleUploadDoc)
		api.POST("/confirm", s.handleConfirmDoc)

		api.GET("/tee/public-key", s.handleGetTEEPublicKey)
	}
}

func (s *Server) handleUploadDoc(c *gin.Context) {
	// Read file data
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Store data
	fileHash, err := s.storage.Store(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store data"})
		return
	}

	// Initiate blockchain upload
	sessionID, err := s.blockchain.InitiateDataUpload(fileHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate blockchain upload"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fileHash":  fileHash,
		"sessionId": sessionID,
	})
}

func (s *Server) handleConfirmDoc(c *gin.Context) {
	var req struct {
		FileHash  string `json:"fileHash"`
		SessionID string `json:"sessionId"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Process in TEE
	result, err := s.tee.ProcessGeneData(req.FileHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process in TEE"})
		return
	}

	// Update result with session ID
	result.SessionID = req.SessionID

	// Confirm on blockchain and mint NFT
	err = s.blockchain.ProcessAndMint(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process blockchain operations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document confirmed and processed successfully",
		"result":  result,
	})
}

func (s *Server) handleGetTEEPublicKey(c *gin.Context) {
	publicKey := s.tee.GetTEEPublicKey()
	c.JSON(http.StatusOK, gin.H{"publicKey": publicKey})
}

func (s *Server) Run() error {
	return s.router.Run(":8080")
}
