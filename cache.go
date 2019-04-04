package raggort

import (
	"sync"
	"time"
)

const (
	NoResponse     = -1
	DefaultTimeout = 5 * time.Second
)

// Cache defines a rapport request cache
type Cache struct {
	data              map[string]chan *HTTPResponse
	timers            map[string]*time.Timer
	timerStopChannels map[string]chan struct{}
	lock              *sync.Mutex
}

// NewCache creates a new empty cache
func NewCache() *Cache {
	return &Cache{
		data:              map[string]chan *HTTPResponse{},
		timers:            map[string]*time.Timer{},
		timerStopChannels: map[string]chan struct{}{},
		lock:              &sync.Mutex{},
	}
}

// AddRequest adds a request to the cache
func (c *Cache) AddRequest(request *Request, timeout time.Duration) chan *HTTPResponse {

	// NoResponse means we don't care about the response
	if timeout == NoResponse {
		return nil
	}

	// Default the timeout if 0 was supplied
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	// Create and start a timeout timer
	timer := time.NewTimer(timeout)
	stopChannel := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopChannel:
				return

			case <-timer.C:
				c.AddResponse(newTimeoutResponse(request.ID))
				return
			}
		}
	}()

	// Set up a channel for the response
	responseChannel := make(chan *HTTPResponse)
	c.setResponseData(request.ID, responseChannel, timer, stopChannel)
	return responseChannel
}

// AddResponse adds a response to the cache
func (c *Cache) AddResponse(response *Response) {

	// Grab the data we need in order to put the response on the right channel
	responseChannel, timer, stopChannel := c.getResponseData(response.ID)
	if responseChannel == nil || timer == nil || stopChannel == nil {
		return
	}

	// Stop the timer and kill the goroutine waiting on it
	close(stopChannel)
	timer.Stop()

	// And write the response to the response channel
	if response.IsError() {
		responseChannel <- response.Error
	} else {
		responseChannel <- response.Body
	}
}

// setResponseData sets the data we need for a response and stores it in our maps
func (c *Cache) setResponseData(id string, responseChannel chan *HTTPResponse, timer *time.Timer, stopChannel chan struct{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[id] = responseChannel
	c.timers[id] = timer
	c.timerStopChannels[id] = stopChannel
}

// getResponseData gets the data from the 3 maps we need for a response, and clears it out of the maps
func (c *Cache) getResponseData(id string) (chan *HTTPResponse, *time.Timer, chan struct{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Grab the data channel for the request
	dataChannel, ok := c.data[id]
	if !ok {
		return nil, nil, nil
	}

	// Grab the timeout timer for the request
	timer, ok := c.timers[id]
	if !ok {
		return nil, nil, nil
	}

	stopChannel, ok := c.timerStopChannels[id]
	if !ok {
		return nil, nil, nil
	}

	delete(c.data, id)
	delete(c.timers, id)
	delete(c.timerStopChannels, id)
	return dataChannel, timer, stopChannel
}
