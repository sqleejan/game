package models

import (
	"sync"
	"time"
)

const (
	Role_Custom = iota
	Role_Assistant
	Role_Admin
)

var (
	UserList = make(map[string]*User)
	userLock sync.Mutex
)

//func init() {
//	UserList = make(map[string]*User)
//	u := User{"user_11111", "leejan"}
//	UserList["user_11111"] = &u
//}

type User struct {
	Id       string
	Username string

	Rooms  map[*Room]*Player
	locker sync.Mutex
}

type Player struct {
	Role    int
	Score   int
	NicName string
}

type UserReq struct {
	UserId   string
	Username string
}

type TmpRespone struct {
	RoomId    string
	RoomName  string
	Admin     string
	StartTime time.Time
	EndTime   time.Time
	Active    bool
	LenUser   int
}

func NewUser(u UserReq) *User {
	userLock.Lock()
	user, ok := UserList[u.UserId]
	if !ok {
		user = &User{Id: u.UserId, Username: u.Username, Rooms: make(map[*Room]*Player)}
		UserList[user.Id] = user
	}
	userLock.Unlock()
	return user
}

func (u *User) AppendRoom(room *Room, role int) {
	u.locker.Lock()
	if _, ok := u.Rooms[room]; !ok {
		u.Rooms[room] = &Player{
			Role: role,
		}
	}
	u.locker.Unlock()
}

//func AddUser(u User) string {
//	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
//	UserList[u.Id] = &u
//	return u.Id
//}

//func GetUser(uid string) (u *User, err error) {
//	if u, ok := UserList[uid]; ok {
//		return u, nil
//	}
//	return nil, errors.New("User not exists")
//}

//func GetAllUsers() map[string]*User {
//	return UserList
//}

//func UpdateUser(uid string, uu *User) (a *User, err error) {
//	if u, ok := UserList[uid]; ok {
//		if uu.Username != "" {
//			u.Username = uu.Username
//		}
//		if uu.Password != "" {
//			u.Password = uu.Password
//		}
//		if uu.Profile.Age != 0 {
//			u.Profile.Age = uu.Profile.Age
//		}
//		if uu.Profile.Address != "" {
//			u.Profile.Address = uu.Profile.Address
//		}
//		if uu.Profile.Gender != "" {
//			u.Profile.Gender = uu.Profile.Gender
//		}
//		if uu.Profile.Email != "" {
//			u.Profile.Email = uu.Profile.Email
//		}
//		return u, nil
//	}
//	return nil, errors.New("User Not Exist")
//}

//func Login(username, password string) bool {
//	for _, u := range UserList {
//		if u.Username == username && u.Password == password {
//			return true
//		}
//	}
//	return false
//}

//func DeleteUser(uid string) {
//	delete(UserList, uid)
//}
