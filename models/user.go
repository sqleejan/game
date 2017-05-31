package models

import "sync"

const (
	Role_Custom = iota
	Role_Assistant
	Role_Admin
	Role_Finace
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
	Role    int    `db:"role"`
	Score   int    `db:"score"`
	NicName string `db:"nickname"`
	Head    string `db:"head"`
	Active  bool   `db:"active"`
}

type UserReq struct {
	UserId   string
	Username string
}

type UserActiveReq struct {
	Nicname string
	Active  bool
}

// type TmpRespone struct {
// 	RoomId    string
// 	RoomName  string
// 	Admin     string
// 	StartTime time.Time
// 	EndTime   time.Time
// 	Active    bool
// 	LenUser   int
// }

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

type SelfObj struct {
	RoomId   int    `json:"room_id"`
	RoomName string `json:"room_name"`
	Role     int    `json:"role"`
	Score    int    `json:"score"`
}

func FetchUserInfo(uid string) (string, map[int]*SelfObj, error) {
	res := map[int]*SelfObj{}

	u, err := cemsdk.GetUser(uid)
	if err != nil {
		return "", res, err
	}

	for rid, room := range RoomList {
		u, ok := room.users[uid]
		if ok {
			res[rid] = &SelfObj{
				RoomId:   rid,
				RoomName: room.name,
				Role:     u.Role,
				Score:    u.Score,
			}
		}

	}

	return u.Nicname, res, nil
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
