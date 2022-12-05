package sub_shard

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Core struct {
	Count uint64
}

func NewCore() *Core {
	return &Core{Count: 0}
}

func (it *Core) Get(req *GetRequest) (*GetResponse, error) {
	return &GetResponse{Value: it.Count}, nil
}

func (it *Core) Incr(req *IncrRequest) (*IncrResponse, error) {
	it.Count++
	return &IncrResponse{Value: it.Count}, nil
}

func (it *Core) save(w io.Writer) error {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, it.Count)
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("writer.Write err: %w", err)
	}
	return nil
}

func (it *Core) load(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("io.ReadAll err: %w", err)
	}
	v := binary.LittleEndian.Uint64(data)
	it.Count = v
	return nil
}
