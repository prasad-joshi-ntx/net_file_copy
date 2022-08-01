package common

type TransferHeader struct {
	BlockSize      int64
	CompressedSize int64
	NumBlocks      int64
	Compression    bool
}
