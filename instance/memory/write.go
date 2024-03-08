package memory

type MemWriter interface {
	Write(offset uint32, data []byte) error
}
