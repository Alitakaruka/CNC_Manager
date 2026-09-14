package CNCService

type BitSet struct {
	data []byte
}

func NewBitSet(size int) *BitSet {
	return &BitSet{
		data: make([]byte, (size+7)/8),
	}
}
func (b *BitSet) Set(i int) {
	byteIndex := i / 8
	bit := i % 8

	b.data[byteIndex] |= 1 << bit
}

func (b *BitSet) Clear(i int) {
	byteIndex := i / 8
	bit := i % 8

	b.data[byteIndex] &^= 1 << bit
}

func (b *BitSet) Has(i int) bool {
	byteIndex := i / 8
	bit := i % 8

	return b.data[byteIndex]&(1<<bit) != 0
}

func (b *BitSet) GetData() []byte {
	copy := make([]byte, len(b.data))
	copy = append(copy, b.data...)
	return copy
}
