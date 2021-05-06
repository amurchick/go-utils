package finalizer

import (
	"testing"
	"time"
)

func TestFinalizerAddRemove(t *testing.T) {
	f := New()
	item := struct{}{}
	f.Add(nil)
	f.Add(item)
	f.Add(nil)
	f.Remove(item)
	if len(f.items) != 2 {
		t.Fail()
	}
}

func TestFinalizerAddRemoveChan(t *testing.T) {
	f := New()
	f.Add(nil)
	a := f.AddChannelItem()
	f.Add(nil)
	f.Remove(a)
	if len(f.items) != 2 {
		t.Fail()
	}
}

func TestFinalizerTimeout(t *testing.T) {
	f := New().SetTimeout(100 * time.Millisecond)
	f.Add(func() {
		time.Sleep(200 * time.Millisecond)
	})
	if f.Run() {
		t.Fail()
	}
}

func TestFinalizerSerial(t *testing.T) {
	f := New().SetTimeout(200 * time.Millisecond)
	f.Add(func() {
		time.Sleep(100 * time.Millisecond)
	})
	f.Add(func() {
		time.Sleep(100 * time.Millisecond)
	})
	ts := time.Now()
	if !f.Run() {
		t.Fail()
	}
	if time.Since(ts) < 200*time.Millisecond {
		t.Fail()
	}
}

func TestFinalizerParallel(t *testing.T) {
	f := New().SetTimeout(200 * time.Millisecond).SetParallel(true)
	f.Add(func() {
		time.Sleep(100 * time.Millisecond)
	})
	f.Add(func() {
		time.Sleep(100 * time.Millisecond)
	})
	ts := time.Now()
	if !f.Run() {
		t.Fail()
	}
	if time.Since(ts) > 150*time.Millisecond {
		t.Fail()
	}
}

func TestFinalizerSerialChan(t *testing.T) {
	f := New().SetTimeout(200 * time.Millisecond)
	a := f.AddChannelItem()
	go func() {
		stop := <-a
		time.Sleep(100 * time.Millisecond)
		stop.Done()
	}()
	b := f.AddChannelItem()
	go func() {
		stop := <-b
		time.Sleep(100 * time.Millisecond)
		stop.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	ts := time.Now()
	if !f.Run() {
		t.Fail()
	}
	if time.Since(ts) < 200*time.Millisecond {
		t.Fail()
	}
}

func TestFinalizerParallelChan(t *testing.T) {
	f := New().SetTimeout(200 * time.Millisecond).SetParallel(true)
	a := f.AddChannelItem()
	go func() {
		stop := <-a
		time.Sleep(100 * time.Millisecond)
		stop.Done()
	}()
	b := f.AddChannelItem()
	go func() {
		stop := <-b
		time.Sleep(100 * time.Millisecond)
		stop.Done()
	}()
	time.Sleep(100 * time.Millisecond)
	ts := time.Now()
	if !f.Run() {
		t.Fail()
	}
	if time.Since(ts) > 150*time.Millisecond {
		t.Fail()
	}
}
