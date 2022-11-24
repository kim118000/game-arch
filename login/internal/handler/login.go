package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kim118000/core/define"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/core/pkg/redis"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/db/model/sys"
	"github.com/kim118000/login/internal/service"
	"github.com/kim118000/protocol/proto/login"
	"net/http"
	"strings"
	"time"
)

var LoginHandler = new(Login)

type Login struct{}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) () {
	username := r.URL.Query().Get("u")
	token := r.URL.Query().Get("token")

	//验证token的正确性
	if token == "" {

	}

	user := new(sys.User)
	user.QueryUserByUserName(username)

	if !user.IsExist() {
		now := time.Now()
		user.CreateTime = now
		user.UpdateTime = now
		user.UserName = username
		user.RoleId = uint64(service.IdSnow.Generate().Int64())
		user.Insert()
	}

	//选择网关
	var gateId uint32 = 1
	var gateUrl = "127.0.0.1:8999"
	//分配游戏服务器
	var gameServerId uint32 = 1

	secret := toolkit.RandString(10, toolkit.LOWERCASE)
	//登录信息
	loginInfo := login.LoginInfo{
		GameServerId: gameServerId,
		GateServerId: gateId,
		RoleId:       user.RoleId,
		Secret:       secret,
	}

	tokenTs := toolkit.TimeUtils.TimeToSecond(time.Now())

	var str strings.Builder
	str.WriteString(fmt.Sprintf("%d|%d|%d|%s", gateId, user.RoleId, tokenTs, secret))
	sign := toolkit.Crypto.Md5(str.String())

	auth := login.LoginAuthInfo{
		RoleId:  user.RoleId,
		TokenTs: tokenTs,
		Gate:    gateUrl,
		Sign:    sign,
	}

	loginKey := fmt.Sprintf("%s%d", define.LoginInfoKey, user.RoleId)
	redis.Client.Set(loginKey, loginInfo, 30*time.Minute)

	data, _ := json.Marshal(auth)
	_, err := w.Write(data)
	if err != nil {
		logger.Log.Errorf("login %s", err)
	}
}
