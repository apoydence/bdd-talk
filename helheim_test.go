package bdd_talk_test

import "net/http"

type mockResponseStorer struct {
	StoreCalled chan bool
	StoreInput  struct {
		resp  chan *http.Response
		index chan int
	}
	StoreOutput struct {
		ret0 chan error
	}
}

func newMockResponseStorer() *mockResponseStorer {
	m := &mockResponseStorer{}
	m.StoreCalled = make(chan bool, 100)
	m.StoreInput.resp = make(chan *http.Response, 100)
	m.StoreInput.index = make(chan int, 100)
	m.StoreOutput.ret0 = make(chan error, 100)
	return m
}
func (m *mockResponseStorer) Store(resp *http.Response, index int) error {
	m.StoreCalled <- true
	m.StoreInput.resp <- resp
	m.StoreInput.index <- index
	return <-m.StoreOutput.ret0
}
