package profile

import (
	"context"
	"time"
)

// WatchEvent describes a change detected on a profile.
type WatchEvent struct {
	ProfileName string
	Kind        WatchEventKind
	At          time.Time
}

// WatchEventKind categorises a profile change.
type WatchEventKind string

const (
	WatchEventModified WatchEventKind = "modified"
	WatchEventDeleted  WatchEventKind = "deleted"
)

// Watcher polls a Store for changes to a set of profiles and emits
// WatchEvents on a channel.
type Watcher struct {
	store    Store
	names    []string
	interval time.Duration
	snapshots map[string]map[string]string
}

// NewWatcher creates a Watcher for the given profile names.
// interval controls how often the store is polled.
func NewWatcher(store Store, names []string, interval time.Duration) *Watcher {
	return &Watcher{
		store:     store,
		names:     names,
		interval:  interval,
		snapshots: make(map[string]map[string]string),
	}
}

// Watch starts polling and sends events to the returned channel.
// It stops when ctx is cancelled, closing the channel.
func (w *Watcher) Watch(ctx context.Context) <-chan WatchEvent {
	ch := make(chan WatchEvent, 16)
	go func() {
		defer close(ch)
		// capture initial state silently
		w.captureAll()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for _, ev := range w.poll() {
					select {
					case ch <- ev:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return ch
}

func (w *Watcher) captureAll() {
	for _, name := range w.names {
		p, err := w.store.Load(name)
		if err != nil {
			continue
		}
		w.snapshots[name] = p.ToEnvMap()
	}
}

func (w *Watcher) poll() []WatchEvent {
	var events []WatchEvent
	for _, name := range w.names {
		p, err := w.store.Load(name)
		if err != nil {
			if _, had := w.snapshots[name]; had {
				events = append(events, WatchEvent{ProfileName: name, Kind: WatchEventDeleted, At: time.Now()})
				delete(w.snapshots, name)
			}
			continue
		}
		current := p.ToEnvMap()
		if !mapsEqual(w.snapshots[name], current) {
			events = append(events, WatchEvent{ProfileName: name, Kind: WatchEventModified, At: time.Now()})
			w.snapshots[name] = current
		}
	}
	return events
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
