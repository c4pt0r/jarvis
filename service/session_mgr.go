package main

import (
	"sync"

	"github.com/ngaut/log"
)

type SessionMgr struct {
	aliveSessions    map[string]*Session
	closingSessionCh chan *Session

	mu sync.Mutex
}

func NewSessionMgr() *SessionMgr {
	return &SessionMgr{
		aliveSessions:    make(map[string]*Session),
		closingSessionCh: make(chan *Session),
	}
}

func (mgr *SessionMgr) AddSession(s *Session) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.aliveSessions[s.ID] = s
	go func(c chan *Session) {
		log.Info("register closing chan")
		session := <-c
		mgr.closingSessionCh <- session
	}(s.DoneChan())
}

func (mgr *SessionMgr) GetSession(ID string) *Session {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if r, exist := mgr.aliveSessions[ID]; exist {
		return r
	}
	return nil
}

func (mgr *SessionMgr) Run() {
	for {
		select {
		case closingSession := <-mgr.closingSessionCh:
			log.Info("Session:", closingSession.ID, "is closing...removing it from SessionMgr")
			mgr.mu.Lock()
			delete(mgr.aliveSessions, closingSession.ID)
			mgr.mu.Unlock()
		}
	}

}
