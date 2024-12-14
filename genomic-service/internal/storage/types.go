package storage

type Storage interface {
	Store(data []byte) (string, error)
	Retrieve(fileHash string) ([]byte, error)
	Delete(fileHash string) error
}
