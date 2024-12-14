package tee

type Processor interface {
	ProcessGeneData(fileHash string) (int, error)
}
