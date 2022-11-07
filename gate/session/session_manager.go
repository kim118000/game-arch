package session

import (
	"errors"
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/gate/constant"
	"github.com/kim118000/gate/logger"
	"sync"
	"sync/atomic"
)

var (
	SessionMgr *SessionManager
)

func init() {
	SessionMgr = NewSessionManager(10, 3000)
}

type SessionManager struct {
	sync.RWMutex
	sessions map[uint64]*Session

	syncClientNumInterval int
	onlineNumber          int64
}

func NewSessionManager(interval int, maxSessions int) *SessionManager {
	return &SessionManager{
		sessions:              make(map[uint64]*Session, maxSessions),
		syncClientNumInterval: interval,
		onlineNumber:          0,
	}
}

func (sm *SessionManager) Add(session *Session) {
	sm.Lock()
	defer sm.Unlock()

	var counter = true
	if sess, ok := sm.sessions[session.GetUserId()]; ok {
		sess.GetConn().RemoveProperty(constant.SessionAttrKey)
		sess.Kick()
		counter = false
	}

	sm.sessions[session.GetUserId()] = session
	if counter {
		sm.SyncOnlineNumber(1)
	}
}

func (sm *SessionManager) Remove(session *Session) bool {
	sm.Lock()
	defer sm.Unlock()

	_, ok := sm.sessions[session.GetUserId()]
	if ok {
		delete(sm.sessions, session.GetUserId())
		sm.SyncOnlineNumber(-1)
		return true
	}
	return false
}

func (sm *SessionManager) Get(userId uint64) (*Session, error) {
	sm.RLock()
	defer sm.RUnlock()

	if sess, ok := sm.sessions[userId]; ok {
		return sess, nil
	}
	return nil, errors.New("session not found")
}

func (sm *SessionManager) LoadOnlineNumber() int64 {
	return atomic.LoadInt64(&sm.onlineNumber)
}

func (sm *SessionManager) SyncOnlineNumber(delta int64) {
	atomic.AddInt64(&sm.onlineNumber, delta)
}

func (sm *SessionManager) SyncClientNumber() {
	logger.Info.Infof("current online number %d", sm.LoadOnlineNumber())
}

func (sm *SessionManager) Init() {
	scheduler.NewAfterTimerBySecondForever(sm.syncClientNumInterval,
		func() {
			sm.SyncClientNumber()
		})
}

func (sm *SessionManager) Start() {

}

func (sm *SessionManager) Stop() {
	sm.RLock()
	var arr = make([]*Session, 0, sm.LoadOnlineNumber())
	for _, sess := range sm.sessions {
		arr = append(arr, sess)
	}
	sm.RUnlock()

	for _, sess := range arr {
		sess.Kick()
	}
}
