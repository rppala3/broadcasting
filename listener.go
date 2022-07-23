package broadcasting

// Listener implements a Subscriber that listen a broadcast channel.
type Listener[T any] struct {
	listenerId  uint
	channel     <-chan T
	broadcaster *Broadcaster[T]
}

// NewListener
func NewListener[T any](
	listenerId uint,
	channel <-chan T,
	broadcaster *Broadcaster[T],
) *Listener[T] {
	return &Listener[T]{
		listenerId,
		channel,
		broadcaster,
	}
}

// Channel returns the channel that receives broadcast messages.
func (lstn *Listener[T]) Channel() <-chan T {
	return lstn.channel
}

// Discard closes the Listener, disabling the reception of further messages.
func (lstn *Listener[T]) Discard() {
	lstn.broadcaster.mutex.Lock()
	defer lstn.broadcaster.mutex.Unlock()
	delete(lstn.broadcaster.channels, lstn.listenerId)
}
