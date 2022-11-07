package session

import (
	"fmt"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/gate/constant"
	"github.com/kim118000/gate/logger"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"time"
)

type Session struct {
	conn       network.IConnection
	userId     uint64
	onlineTime time.Time

	authRedisKey string
	authSuccess  int32

	sessionMgr *SessionManager
}

func NewSession(conn network.IConnection, mgr *SessionManager) *Session {
	s := &Session{
		conn:        conn,
		sessionMgr:  mgr,
		authSuccess: 0,
	}

	s.conn.SetProperty(constant.SessionAttrKey, s)
	return s
}

func (s *Session) GetUserId() uint64 {
	return s.userId
}

func (s *Session) GetConn() network.IConnection {
	return s.conn
}

func (s *Session) Send(msg interface{}) {

}

func (s *Session) SendProtobuf(msgId uint32, protobuf proto.Message) {
	msg := network.GetMessage()
	msg.Init(msgId, protobuf)
	s.conn.SendMsg(msg)
}

func (s *Session) Kick() {
	if s.conn.IsClose() {
		return
	}
	s.conn.Stop()
	logger.Info.Infof("kick session offline %s", s)
}

func (s *Session) RemoveSession() {
	ok := s.sessionMgr.Remove(s)
	if ok {
		s.OffLine()
	}
}

//通知game下线
func (s *Session) OffLine() {

}

func (s *Session) Auth(userId uint64, tokenTs int, sign string) bool {
	now := int(time.Now().UnixMilli() / 1000)
	diff := now - tokenTs
	if diff < -constant.AuthenticateTokenDiff || diff > constant.AuthenticateTokenDiff {
		//return false
	}

	s.authRedisKey = fmt.Sprintf("AUTH:%d", userId)

	//TODO
	//从redis加载AuthContract
	var authContract *AuthContract

	if authContract != nil && authContract.GateId == 1 {
		md5 := toolkit.Crypto.Md5(fmt.Sprintf("%d%d%s", userId, tokenTs, authContract.Secret))
		if sign == md5 {
			s.AuthSuccess(userId)
			return true
		}
	}

	//todo 临时通过
	s.AuthSuccess(userId)
	return true
}

func (s *Session) AuthSuccess(userId uint64) {
	s.userId = userId
	s.onlineTime = time.Now()

	s.sessionMgr.Add(s)
	s.SetAuthSuccess()
	s.LoadLoginInfo()

	scheduler.NewAfterTimerBySecondOnce(constant.AuthenticateTokenRefreshInterval, func() {
		s.RefreshAuthToken()
	})
}

func (s *Session) LoadLoginInfo() {

}

func (s *Session) IsAuthSuccess() bool {
	if atomic.LoadInt32(&s.authSuccess) == 1 {
		return true
	}
	return false
}

func (s *Session) SetAuthSuccess() {
	atomic.StoreInt32(&s.authSuccess, 1)
}

//定时刷新token过期时间
func (s *Session) RefreshAuthToken() {
	logger.Info.Infof("refresh auth token expire time user=", s)
}

func (s *Session) GetGameNode() network.IConnection {
	return nil
}

func (s *Session) String() string {
	return fmt.Sprintf("[userId=%d]", s.userId)
}
