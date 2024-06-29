package core

type Storage interface {
    Put(*Block) error
    Get(hash []byte) (*Block, error)
}

type MemoryStore struct {
    store map[string]*Block
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        store: make(map[string]*Block),
    }
}

func (s *MemoryStore) Put(b *Block) error {
    // Assuming the block has a method Hash() that returns its hash as a string
    hash := string(b.Hash())
    s.store[hash] = b
    return nil
}

func (s *MemoryStore) Get(hash []byte) (*Block, error) {
    b, ok := s.store[string(hash)]
    if !ok {
        return nil, fmt.Errorf("block not found")
    }
    return b, nil
}