package memory

type MemReader interface {
	Read(offset uint32, length uint32) ([]byte, error)
}
