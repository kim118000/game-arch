package session

import (
	"fmt"
	"github.com/kim118000/core/define"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/network"
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/core/pkg/scheduler"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/gate/internal/conf"
	"github.com/kim118000/gate/internal/constant"
	"github.com/kim118000/protocol/proto/login"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"strings"
	"sync/atomic"
	"time"
)

type Session struct {
	conn       network.IConnection
	roleId     uint64
	onlineTime time.Time

	authRedisKey string
	authSuccess  int32
	loginInfo    *login.LoginInfo
	sessionMgr   *SessionManager
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

func (s *Session) GetRoleId() uint64 {
	return s.roleId
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
	logger.Log.Infof("kick session offline %s", s)
}

func (s *Session) RemoveSession() {
	ok := s.sessionMgr.Remove(s)
	if ok {
		s.OffLine()
	}
}

//通知game下线
func (s *Session) OffLine() {
	s.loginInfo = nil
}

func (s *Session) Auth(roleId uint64, tokenTs uint32, sign string) bool {
	now := toolkit.TimeUtils.TimeToSecond(time.Now())
	diff := int(now - tokenTs)
	if diff < -constant.AuthenticateTokenDiff || diff > constant.AuthenticateTokenDiff {
		return false
	}

	s.authRedisKey = fmt.Sprintf("%s%d", define.LoginInfoKey, roleId)

	//TODO
	ls := redis.Client.Get(s.authRedisKey)
	if ls == "" {
		return false
	}

	var loginInfo = &login.LoginInfo{}
	_ = protojson.Unmarshal(toolkit.StringToBytes(ls), loginInfo)

	if loginInfo.GateServerId == uint32(conf.Config.ServerId) {

		var str strings.Builder
		str.WriteString(fmt.Sprintf("%d|%d|%d|%s", conf.Config.ServerId, roleId, tokenTs, loginInfo.Secret))
		md5 := toolkit.Crypto.Md5(str.String())
		if sign == md5 {
			s.AuthSuccess(roleId, loginInfo)
			return true
		}
	}

	return false
}

func (s *Session) AuthSuccess(roleId uint64, loginInfo *login.LoginInfo) {
	s.roleId = roleId
	s.onlineTime = time.Now()

	s.sessionMgr.Add(s)
	s.SetAuthSuccess()
	s.LoadLoginInfo(loginInfo)

	scheduler.NewAfterTimerBySecondOnce(constant.AuthenticateTokenRefreshInterval, func() {
		s.RefreshLoginInfo()
	})
}

func (s *Session) LoadLoginInfo(loginInfo *login.LoginInfo) {
	s.loginInfo = loginInfo
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
func (s *Session) RefreshLoginInfo() {
	loginKey := fmt.Sprintf("%s%d", define.LoginInfoKey, s.roleId)
	redis.Client.Set(loginKey, s.loginInfo, 30*time.Minute)
}

func (s *Session) GetGameNode() network.IConnection {
	return nil
}

func (s *Session) String() string {
	return fmt.Sprintf("[roleId=%d]", s.roleId)
}
