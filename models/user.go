package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	db "supermarket-go/local/leveldb"
	"supermarket-go/log"
	"time"
)

// 数据库中存放的UserList
var (
	UserList map[string]*User
)

func init() {
	data, err := db.Database.Get([]byte("UserList"), nil)
	if err != nil {
		log.ErrorLog("获取数据错误,key:UserList", err)
	}
	if len(data) > 0 {
		// var users []User
		// err := json.Unmarshal(data, &users)
		err := json.Unmarshal(data, &UserList)
		if err != nil {
			db.Database.Delete([]byte("UserList"), nil)
			log.ErrorLog("反格式化数据错误,value：", UserList, ".开始执行删除")
		}
		// initUserList(users)
	} else {
		// u := []User{User{"user_11111", "astaxie", "11111", Profile{"male", 20, "Singapore", "astaxie@gmail.com"}},
		// 	User{"user_11112", "astaxie", "11112", Profile{"male", 22, "China", "wangfeifan@gmail.com"}},
		// }
		UserList = make(map[string]*User)
		UserList["user_11111"] = &User{"user_11111", "astaxie", "11111", Profile{"male", 20, "Singapore", "astaxie@gmail.com"}}
		UserList["user_11112"] = &User{"user_11112", "astaxie", "11112", Profile{"male", 22, "China", "wangfeifan@gmail.com"}}
		value, err := json.Marshal(UserList)
		if err != nil {
			log.ErrorLog("格式化数据错误,value：", UserList)
		}
		if err := db.Database.Put([]byte("UserList"), value, nil); err != nil {
			log.ErrorLog("存入数据库失败，key：UserList，value：", UserList, err)
		}
		// UserList = initUserList(u)
	}
}

// initUserList 赋值initUserList
func initUserList(users []User) map[string]*User {
	UserList = make(map[string]*User)
	fmt.Println("users:", users)
	for _, u := range users {
		UserList[u.ID] = &u
	}
	return UserList
}

// User 用户对象
type User struct {
	ID       string
	Username string
	Password string
	Profile  Profile
}

// Profile 用户简况
type Profile struct {
	Gender  string
	Age     int
	Address string
	Email   string
}

// AddUser 增加一个用户
func AddUser(u User) string {
	u.ID = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	UserList[u.ID] = &u
	return u.ID
}

// GetUser 获得一个用户
func GetUser(uID string) (u *User, err error) {
	if u, ok := UserList[uID]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

// GetAllUsers 获得所有用户
func GetAllUsers() map[string]*User {
	return UserList
}

// UpdateUser 根据uID更新一个用户
func UpdateUser(uID string, uu *User) (a *User, err error) {
	if u, ok := UserList[uID]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}
		if uu.Profile.Age != 0 {
			u.Profile.Age = uu.Profile.Age
		}
		if uu.Profile.Address != "" {
			u.Profile.Address = uu.Profile.Address
		}
		if uu.Profile.Gender != "" {
			u.Profile.Gender = uu.Profile.Gender
		}
		if uu.Profile.Email != "" {
			u.Profile.Email = uu.Profile.Email
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

// Login 用户登录
func Login(username, password string) bool {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

// DeleteUser 根据uID删除用户
func DeleteUser(uID string) {
	delete(UserList, uID)
}
