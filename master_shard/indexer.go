package master_shard


type FuncIndexer[T any] struct {
	ver uint64
	indexer func(item T) []string
}

func NewFuncIndexer[T any](ver uint64, indexer func(item T) []string) Indexer[T] {
	if ver == 0 {
		panic("NewFuncIndexer ver > 0")
	}
	return &FuncIndexer[T]{ver: ver, indexer: indexer}
}

func (it *FuncIndexer[T]) Index(item T) []string {
	return it.indexer(item)
}

func (it *FuncIndexer[T]) Ver() uint64 {
	return it.ver
}
