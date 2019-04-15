package members

import "math/rand"

type Watcher struct {
	name string
	ch   chan []string
	// TODO: refactor
	close  func()
	update func()
}

func (w *Watcher) ForceUpdate() {
	w.update()
}

func (w *Watcher) Watch() <-chan []string {
	return w.ch
}

func (w *Watcher) Close() {
	w.close()
}

func (m *Manager) Watch(serviceName string) *Watcher {
	var wtch *Watcher
	ch := make(chan []string, 1)
	wtch = &Watcher{
		ch:   ch,
		name: serviceName,
		close: func() {
			m.lock.Lock()
			watchers := m.watchers[serviceName]
			for i, w := range watchers {
				if w == wtch {
					watchers = append(watchers[:i], watchers[i+1:]...)
					break
				}
			}
			m.watchers[serviceName] = watchers
			m.lock.Unlock()
			close(ch)
		},
		update: func() {
			m.updateWatchers()
		},
	}
	m.lock.Lock()
	if m.watchers == nil {
		m.watchers = make(map[string][]*Watcher)
	}
	m.watchers[serviceName] = append(m.watchers[serviceName], wtch)
	m.lock.Unlock()
	return wtch
}

func (m *Manager) updateWatchers() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.unsafeUpdateWatchers()
}

func (m *Manager) unsafeUpdateWatchers() {
	// TODO: notify only changed
	for serviceName, watchers := range m.watchers {
		addresses := m.findUnsafe(serviceName)
		rand.Shuffle(len(addresses), func(i, j int) {
			addresses[i], addresses[j] = addresses[j], addresses[i]
		})
		for _, w := range watchers {
			select {
			case w.ch <- addresses:
			default:
			}
		}
	}
}
