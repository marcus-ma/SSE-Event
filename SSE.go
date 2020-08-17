package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)
const (
	streamEventBufferSize = 10
)


type StreamHandler struct {
	requests map[*http.Request]chan []byte
	mu       sync.RWMutex
	done     chan struct{}
}

func NewStreamHandler() *StreamHandler {
	return &StreamHandler{}
}

func (sh *StreamHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Make sure that the writer supports flushing.
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	events := sh.register(req)
	defer sh.unregister(req)

	notify := rw.(http.CloseNotifier).CloseNotify()

	rw.Header().Add("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	for {
		select {
		case <-notify:
			// client is gone
			return
		case event := <-events:
			_, err := rw.Write(event)
			if err != nil {
				return
			}
			f.Flush()
		}
	}
}
func (sh *StreamHandler) register(req *http.Request) <-chan []byte {
	sh.mu.RLock()
	events, ok := sh.requests[req]
	sh.mu.RUnlock()
	if ok {
		return events
	}

	events = make(chan []byte, streamEventBufferSize)
	sh.mu.Lock()
	sh.requests[req] = events
	sh.mu.Unlock()
	return events
}
func (sh *StreamHandler) unregister(req *http.Request) {
	sh.mu.Lock()
	delete(sh.requests, req)
	sh.mu.Unlock()
}


func (sh *StreamHandler) writeToRequests(eventBytes []byte) error {
	var b bytes.Buffer
	_, err := b.Write([]byte("data:"))
	if err != nil {
		return err
	}

	_, err = b.Write(eventBytes)
	if err != nil {
		return err
	}
	_, err = b.Write([]byte("\n\n"))
	if err != nil {
		return err
	}
	dataBytes := b.Bytes()
	sh.mu.RLock()

	for _, requestEvents := range sh.requests {
		select {
		case requestEvents <- dataBytes:
		default:
		}
	}
	sh.mu.RUnlock()

	return nil
}

// Start begins watching the in-memory circuit breakers for metrics
func (sh *StreamHandler) Start() {
	sh.requests = make(map[*http.Request]chan []byte)
	sh.done = make(chan struct{})
	go sh.loop()
}

// Stop shuts down the metric collection routine
func (sh *StreamHandler) Stop() {
	close(sh.done)
}



func (sh *StreamHandler) loop() {
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			eventBytes, err := json.Marshal(struct {
				Msg string `json:"msg"`
			}{Msg: fmt.Sprintf("The server time is:%d",time.Now().Unix())})
			if err!=nil {
				fmt.Println(err)
			}
			sh.writeToRequests(eventBytes)
		case <-sh.done:
			return
		}
	}
}




func main() {

	sh := NewStreamHandler()
	sh.Start()

	fmt.Println("Open 8081")
	http.ListenAndServe(":8081",sh)


}
