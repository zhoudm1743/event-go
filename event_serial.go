package event_go

import (
	"go.uber.org/zap"
	"sync"
)

func (s *EventService) serialWorker() {
	for e := range s.serialQueue {
		s.mu.RLock()
		subs, ok := s.subs[e.Name]
		s.mu.RUnlock()

		if ok {
			for _, entry := range subs {
				func() {
					defer func() {
						if err := recover(); err != nil {
							s.logger.Error("Serial subscription panic",
								zap.String("event", e.Name),
								zap.Any("error", err),
								zap.Stack("stack"),
							)
						}
					}()
					entry.fn(*e)
				}()
			}
		}

		e.Name = ""
		e.Data = nil
		s.pool.Put(e)
	}
}