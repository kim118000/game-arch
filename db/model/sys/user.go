package sys

import (
	"github.com/kim118000/db/model"
	"github.com/kim118000/db/persistence"
)

type User struct {
	model.Base
	UserName string `json:"userName,omitempty"          gorm:"type:VARCHAR(50);column:user_name;not null;index:idx_user_name"`
}

func (u *User) Insert() {
	_, _ = persistence.SysDao.DoSave(u)
}

func (u *User) Update() {
	_, _ = persistence.SysDao.DoUpdate(u)
}

func (u *User) QueryUserByUserName(userName string) {
	_ = persistence.SysDao.DoQuery(u, "user_name = ?", userName)
}

func (u *User) QueryUserByRoleId(roleId uint64) {
	_ = persistence.SysDao.DoQuery(u, "role_id = ?", roleId)
}
