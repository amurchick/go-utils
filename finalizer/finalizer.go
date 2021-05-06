package finalizer

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/amurchick/go-utils/logger"
)

var log = logger.Log

type Item struct {
	item interface{}
	from string
	mux  sync.Mutex
	c    Channel
}

func (item *Item) fillFileAndLine(args ...int) {
	ofs := 1
	if len(args) > 0 {
		ofs = args[0]
	}
	_, file, line, ok := runtime.Caller(ofs)
	if ok {
		file = fmt.Sprintf("%v:%v", filepath.Base(file), line)
	} else {
		file = "n/a"
	}
	item.mux.Lock()
	item.from = file
	item.mux.Unlock()
}

type Finalizer struct {
	mux      sync.Mutex
	parallel bool
	timeout  time.Duration
	items    []*Item
}

func New() *Finalizer {
	this := &Finalizer{
		timeout: 5 * time.Second,
	}
	return this
}

func (this *Finalizer) SetParallel(parallel bool) *Finalizer {
	this.parallel = parallel
	return this
}

func (this *Finalizer) SetTimeout(timeout time.Duration) *Finalizer {
	this.timeout = timeout
	return this
}

func (this *Finalizer) Add(item interface{}) {
	itemToAdd := &Item{item: item}
	itemToAdd.fillFileAndLine(2)
	this.mux.Lock()
	this.items = append(this.items, itemToAdd)
	this.mux.Unlock()
}

type Channel chan ChannelChan
type ChannelChan chan struct{}

func (c ChannelChan) Done() {
	if len(c) == 0 {
		c <- struct{}{}
	}
}

func (this *Finalizer) AddChannelItem() Channel {
	item := &Item{
		c: make(Channel, 1),
	}
	item.item = func() {
		cc := make(ChannelChan, 1)
		item.c <- (cc)
		<-cc
	}
	item.fillFileAndLine(2)
	this.mux.Lock()
	this.items = append(this.items, item)
	this.mux.Unlock()
	return item.c
}

func (this *Finalizer) Remove(item interface{}) {
	this.mux.Lock()
	newLen := 0
	for idx := range this.items {
		if newLen != idx {
			this.items[newLen] = this.items[idx]
		}
		if this.items[idx].item != item && this.items[idx].c != item {
			newLen++
		}
	}
	this.items = this.items[:newLen]
	this.mux.Unlock()
}

func (this *Finalizer) Run() (ok bool) {
	var wg sync.WaitGroup
	ok = true
	for idx := range this.items {
		fn := func(idx int) {
			done := make(chan bool, 2)
			timer := time.AfterFunc(this.timeout, func() {
				log.Warn("finalizer timeout %v (added at %s)", this.timeout, this.items[idx].from)
				ok = false
				done <- true
			})
			go func() {
				this.run(idx)
				timer.Stop()
				done <- true
			}()
			<-done
		}
		if this.parallel {
			wg.Add(1)
			go func(idx int) {
				fn(idx)
				wg.Done()
			}(idx)
		} else {
			fn(idx)
		}
	}
	if this.parallel {
		wg.Wait()
	}
	return
}

func (this *Finalizer) run(idx int) (ok bool) {
	if this.items[idx].item == nil {
		return
	}
	ok = true
	switch value := this.items[idx].item.(type) {
	case func():
		value()
	default:
		log.Warn("finalizer unknown type %T", this.items[idx].item)
		ok = false
	}
	return
}
