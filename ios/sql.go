package ios

import (
	"time"

	"github.com/u03013112/ss-ios-purchase/sql"
)

// User : ios 用户表，用于存储UUID和有效期
type User struct {
	sql.BaseModel
	UUID       string    `json:"uuid,omitempty"`
	ExpireDate time.Time `json:"expireDate,omitempty"`
	Sj         string    `json:"sj,omitempty"`
	Token      string    `json:"token,omitempty"`
	Online     bool      `json:"online,omitempty"`
}

// TableName :
func (User) TableName() string {
	return "ios_user"
}

// InitDB : 初始化表格，建议在整个数据库初始化之后调用
func InitDB() {
	sql.GetInstance().AutoMigrate(&User{}, &Bills{})
}

func getOrCreateUserByUUID(uuid string) User {
	var user User
	db := sql.GetInstance().First(&user, "uuid = ?", uuid)
	if db.RowsAffected == 0 {
		user.UUID = uuid
		user.ExpireDate = time.Unix(0, 0)
		user.Sj = ""
		user.Token = ""
		user.Online = false
		sql.GetInstance().Create(&user)
	}
	return user
}

func (u *User) updateToken(token string) {
	sql.GetInstance().Model(u).Update(User{
		Token: token,
	})
}

func getUserByToken(token string) (*User, error) {
	var user User
	db := sql.GetInstance().Model(&User{}).Where("token = ?", token).First(&user)
	if db.RowsAffected == 1 {
		return &user, nil
	}
	return nil, db.Error
}

func (u *User) updateExpireDate(t time.Time, sj string) error {
	db := sql.GetInstance().Model(u).Update(User{
		ExpireDate: t,
		Sj:         sj,
	})
	return db.Error
}
