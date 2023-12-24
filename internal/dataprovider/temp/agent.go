package temp

import (
	"context"
	"time"
)

func (s *Storage) agentTTLExpire(t time.Time) bool {
	return time.Now().After(t.Add(agentHandleTTL))
}

func (s *Storage) TryLockBookHandle(ctx context.Context, bookID int) bool {
	s.lockBookHandleMutex.Lock()
	defer s.lockBookHandleMutex.Unlock()

	t, ok := s.lockBookHandle[bookID]
	if ok && !s.agentTTLExpire(t) {
		return false
	}

	s.lockBookHandle[bookID] = time.Now()

	return true
}

func (s *Storage) UnLockBookHandle(ctx context.Context, bookID int) {
	s.lockBookHandleMutex.Lock()
	defer s.lockBookHandleMutex.Unlock()

	delete(s.lockBookHandle, bookID)
}

func (s *Storage) HasLockBookHandle(ctx context.Context, bookID int) bool {
	s.lockBookHandleMutex.RLock()
	defer s.lockBookHandleMutex.RUnlock()

	t, ok := s.lockBookHandle[bookID]
	if ok && s.agentTTLExpire(t) {
		ok = false
	}

	return ok
}

func (s *Storage) TryLockPageHandle(ctx context.Context, bookID int, pageNumber int) bool {
	s.lockPageHandleMutex.Lock()
	defer s.lockPageHandleMutex.Unlock()

	id := pageSimple{BookID: bookID, PageNumber: pageNumber}

	t, ok := s.lockPageHandle[id]
	if ok && !s.agentTTLExpire(t) {
		return false
	}

	s.lockPageHandle[id] = time.Now()

	return true
}

func (s *Storage) UnLockPageHandle(ctx context.Context, bookID int, pageNumber int) {
	s.lockPageHandleMutex.Lock()
	defer s.lockPageHandleMutex.Unlock()

	id := pageSimple{BookID: bookID, PageNumber: pageNumber}

	delete(s.lockPageHandle, id)
}

func (s *Storage) HasLockPageHandle(ctx context.Context, bookID int, pageNumber int) bool {
	s.lockPageHandleMutex.RLock()
	defer s.lockPageHandleMutex.RUnlock()

	id := pageSimple{BookID: bookID, PageNumber: pageNumber}

	t, ok := s.lockPageHandle[id]
	if ok && s.agentTTLExpire(t) {
		ok = false
	}

	return ok
}
