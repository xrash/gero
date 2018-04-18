package msg

import (
	"sync"
)

type MessengerBag struct {
	listSize int
	listCap  int
	lock     *sync.Mutex
	bag      map[string][]string
}

func NewMessengerBag(listSize, listCap int) *MessengerBag {
	return &MessengerBag{
		listSize: listSize,
		listCap:  listCap,
		lock:     &sync.Mutex{},
		bag:      make(map[string][]string),
	}
}

func (b *MessengerBag) Capture(cid, mid string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	list := b.bag[cid]
	if list == nil {
		list = make([]string, 0)
	}

	list = append(list, mid)

	if len(list) > b.listSize {
		list = list[len(list)-b.listCap:]
	}

	b.bag[cid] = list
}

func (b *MessengerBag) Len() int {
	return len(b.bag)
}

func (b *MessengerBag) WithBag(fn func(map[string][]string) map[string][]string) {
	b.lock.Lock()
	b.bag = fn(b.bag)
	b.lock.Unlock()
}
