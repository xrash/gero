package stats

import (
	"sync"
)

type Bag struct {
	intValues     map[string]int
	intValuesLock *sync.Mutex
}

func NewBag() *Bag {
	return &Bag{
		intValues:     make(map[string]int),
		intValuesLock: &sync.Mutex{},
	}
}

func (b *Bag) IntSet(key string, value int) {
	b.intValuesLock.Lock()
	b.intValues[key] = value
	b.intValuesLock.Unlock()
}

func (b *Bag) IntAdd(key string, value int) {
	b.intValuesLock.Lock()
	b.intValues[key] = b.intValues[key] + value
	b.intValuesLock.Unlock()
}

func (b *Bag) IntGet(key string) int {
	b.intValuesLock.Lock()
	i := b.intValues[key]
	b.intValuesLock.Unlock()
	return i
}

func (b *Bag) IntDel(key string) {
	b.intValuesLock.Lock()
	delete(b.intValues, key)
	b.intValuesLock.Unlock()
}
