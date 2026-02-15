package orchestrator

import (
	"log/slog"
	"sync"
)

// EventDispatcher ensures thread-safe, sequential event delivery to the callback.
type EventDispatcher struct {
	callback EventCallback
	eventCh  chan event
	wg       sync.WaitGroup
	closed   bool
	mu       sync.Mutex
	traceID  string
}

type event struct {
	Type string
	Data string
}

// NewEventDispatcher creates a new event dispatcher.
func NewEventDispatcher(traceID string, callback EventCallback) *EventDispatcher {
	if callback == nil {
		return &EventDispatcher{callback: nil, traceID: traceID}
	}

	d := &EventDispatcher{
		callback: callback,
		eventCh:  make(chan event, 100),
		traceID:  traceID,
	}

	d.wg.Add(1)

	d.wg.Add(1)
	go d.dispatchLoop()

	return d
}

func (d *EventDispatcher) dispatchLoop() {
	defer d.wg.Done()
	for e := range d.eventCh {
		// Recover from panic in callback to protect dispatcher loop
		func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("event dispatcher: recovered from panic", "panic", r, "trace_id", d.traceID)
				}
			}()
			d.callback(e.Type, e.Data)
		}()
	}
}

// Send sends an event strictly sequentially.
func (d *EventDispatcher) Send(eventType, eventData string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.callback == nil || d.closed {
		return
	}

	// Non-blocking send with drop policy if full?
	// Or blocking? Blocking is safer for correctness, but can stall execution.
	// Given we are in an agent executor, stalling on slow UI is bad.
	// But dropping events is also bad.
	// Let's use a reasonably large buffer and blocking send, assuming UI consumes fast enough.
	// Or use select with warning.

	d.eventCh <- event{Type: eventType, Data: eventData}
}

// Close closes the dispatcher and waits for all events to be processed.
func (d *EventDispatcher) Close() {
	d.mu.Lock()
	if d.callback == nil || d.closed {
		d.mu.Unlock()
		return
	}
	d.closed = true
	d.mu.Unlock()

	close(d.eventCh)
	d.wg.Wait()
}
