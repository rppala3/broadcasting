package broadcasting

import "sync"

// Broadcaster implements a Publisher.
type Broadcaster[T any] struct {
	mutex      sync.Mutex
	bufferSize int
	channels   map[uint]chan<- T
	nextID     uint
	closed     bool
}

// NewBroadcaster returns a new Broadcaster (channelSize = 0 means un-buffered).
func NewBroadcaster[T any](channelSize int) *Broadcaster[T] {
	return &Broadcaster[T]{
		bufferSize: channelSize,
	}
}

// Listen returns a Listener for the broadcast channel.
func (bcast *Broadcaster[T]) Listen(data T) *Listener[T] {
	bcast.mutex.Lock()
	defer bcast.mutex.Unlock()

	if bcast.channels == nil {
		bcast.channels = make(map[uint]chan<- T)
	}

	if bcast.channels[bcast.nextID] != nil { // why? Is it necessary?
		bcast.nextID++
	}

	channel := make(chan T, bcast.bufferSize)
	bcast.channels[bcast.nextID] = channel

	if bcast.closed {
		close(channel)
	}

	return NewListener(bcast.nextID, channel, bcast)
}

// Send broadcasts a message to the channel.
// Sending on a closed channel causes a runtime panic.
func (bcast *Broadcaster[T]) Send(data T) {
	bcast.mutex.Lock()
	defer bcast.mutex.Unlock()

	if bcast.closed {
		panic("broadcast: send after close")
	}

	for _, channel := range bcast.channels {
		channel <- data
	}
}

// Close closes all channels, disabling the sending of further messages.
func (bcast *Broadcaster[T]) Close() {
	bcast.mutex.Lock()
	defer bcast.mutex.Unlock()

	if !bcast.closed {
		bcast.closed = true
		for _, channel := range bcast.channels {
			close(channel)
		}
	}
}
