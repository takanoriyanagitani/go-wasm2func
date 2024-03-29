package memory

import (
	wa "github.com/tetratelabs/wazero/api"

	wzim "github.com/takanoriyanagitani/go-wasm2func/instance/memory"
	util "github.com/takanoriyanagitani/go-wasm2func/util"
)

type RwMem struct {
	wa.Memory
}

func (m RwMem) Read(offset uint32, length uint32) ([]byte, error) {
	data, ok := m.Memory.Read(offset, length)
	return util.Select(
		func() ([]byte, error) { return nil, wzim.ErrOutOfRange },
		func() ([]byte, error) { return data, nil },
		ok,
	)()
}

func (m RwMem) WriteDirect(offset uint32, data []byte) error {
	var ok bool = m.Memory.Write(offset, data)
	return util.Select(
		wzim.ErrOutOfRange,
		nil,
		ok,
	)
}

func (m RwMem) WriteByRead(offset uint32, data []byte) error {
	target, ok := m.Memory.Read(offset, uint32(len(data)))
	if !ok {
		return wzim.ErrOutOfRange
	}
	copy(target, data)
	return nil
}

func (m RwMem) Write(offset uint32, data []byte) error {
	return m.WriteDirect(offset, data)
}

func (m RwMem) AsIf() wzim.ReadWriteMemory { return m }
