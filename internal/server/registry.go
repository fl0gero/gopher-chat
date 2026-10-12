package server

import "sync"

// The Registry stores all active rooms by name and protects access to this map using a mutex
type Registry struct {
	mu    sync.Mutex
	Rooms map[string]*room
}

// NewRegistry creates an empty registry with a pre-initialised room map.

func NewRegistry() *Registry {
	return &Registry{
		Rooms: make(map[string]*room),
	}
}

// getOrCreate returns an existing room with the name "name", if one
// exists, or creates a new one, starts its event processing loop (run) in
// a separate goroutine, and stores it in the registry for later reuse.
// The mutex ensures that, in the event of concurrent requests for the same
// room, it will be created only once and not duplicated.
func (reg *Registry) getOrCreate(name string) *room {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	value, ok := reg.Rooms[name]
	if ok {
		return value
	}

	r := &room{
		clients: make(map[*client]bool),
		join:    make(chan *client),
		leave:   make(chan *client),
		forward: make(chan Message),
	}
	reg.Rooms[name] = r
	go r.run()
	return r

}
