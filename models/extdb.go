package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type DBRoomPost struct {
	//	Id       int64     `db:"id","primarykey" json:"id"`
	UserId   string    `db:"uid","primarykey" json:"uid"`
	Nicname  string    `db:"nicname" json:"nicname"`
	RoomName string    `db:"room_name" json:"room_name"`
	CreateAt time.Time `db:"create_at" json:"create_at"`
	Duration int       `db:"duration" json:"duration"`
	Active   bool      `db:"active" json:"active"`
}

func (post DBRoomPost) TableName() string {
	return "room_post"
}

func (post *DBRoomPost) Insert() error {
	post.CreateAt = time.Now()
	return dBEngine.Insert(post)
}

func (post *DBRoomPost) Delete() error {
	_, err := dBEngine.Delete(post)
	return err
}

func (post *DBRoomPost) Update() error {
	_, err := dBEngine.Update(post)
	return err
}

func (post *DBRoomPost) Fetch() error {
	query := fmt.Sprintf(`select * from %s where uid = "%s"`, post.TableName(), post.UserId)
	//fmt.Println(query)
	return dBEngine.SelectOne(post, query)
}

func ListRoomPosts(page int, size int) (interface{}, error) {
	list := []*DBRoomPost{}
	query := fmt.Sprintf(`select * from %s`, DBRoomPost{}.TableName())
	_, err := dBEngine.Select(&list, query)
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Pagination
		Data []*DBRoomPost `json:"data"`
	}{}
	if size == 0 {
		size = 10
	}
	if page == 0 {
		page = 1
	}
	start, end := PageLocate(len(list), size, page)
	resp.Total = len(list)
	resp.TotalPage = resp.Total / size
	if resp.Total%size != 0 {
		resp.TotalPage += 1
	}
	resp.Data = list[start:end]
	return resp, nil
}

func (post *DBRoomPost) CreateRoom() error {
	req := RoomReq{}
	req.Base = 10
	req.CountDown = 1
	req.CountUp = 3
	req.Duration = post.Duration
	req.RedCountDown = 2
	req.RedCountUp = 5
	req.RedDown = 1.0
	req.RedUp = 9.9
	req.RoomName = post.RoomName
	req.ScoreLimit = -1000
	req.Timeout = 3
	req.UserId = post.UserId
	req.UserLimit = 200
	_, err := CreateRoom(&req)
	if err != nil {
		return err
	}
	post.Fetch()
	post.Active = true
	post.Update()

	return err
}

func GenerateName(uid string) string {
	return RoomNamePrefix + uid + string(Krand(8, 0))
}

type DBRed struct {
	Id       int64     `db:"id","primarykey" json:"id"`
	UserId   string    `db:"uid","primarykey" json:"uid"`
	CreateAt time.Time `db:"create_at" json:"create_at"`
	RedId    string    `db:"red_id" json:"red_id"`
	RoomId   int       `db:"room_id" json:"room_id"`
	Score    string    `db:"score" json:"score"`
	NicName  string    `db:"nick_name" json:"nicname"`
}

func (red DBRed) TableName() string {
	return "db_red"
}

func (red *DBRed) Insert() error {
	red.CreateAt = time.Now()
	return dBEngine.Insert(red)
}

func RedInsert(roomid int, redid string, uid string, nicname string, score string) error {
	dbred := &DBRed{
		UserId:  uid,
		RedId:   redid,
		RoomId:  roomid,
		Score:   score,
		NicName: nicname,
	}
	return dbred.Insert()
}

type redPList []*DBRed

func (s redPList) Less(i, j int) bool {
	return s[i].Id > s[j].Id
}

func (s redPList) Len() int {
	return len(s)
}

func (s redPList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func RedList(roomid int, redid string) ([]*DBRed, error) {
	res := []*DBRed{}
	query := fmt.Sprintf(`select * from db_red where room_id=? and red_id=?`)
	_, err := dBEngine.Select(&res, query, roomid, redid)
	if err != nil {
		return nil, err
	}
	sort.Sort(redPList(res))

	room, ok := RoomList[roomid]
	if ok {
		for i, red := range res {
			u, ok := room.users[red.UserId]
			if ok {
				res[i].NicName = u.NicName
			}
		}
	}

	return res, nil
}

type UpdateInfo struct {
	MSG string
}

var filelock sync.RWMutex

func InfoUpdate(msg string) error {
	filelock.Lock()
	defer filelock.Unlock()
	f, err := os.Create(".amsg")
	if err != nil {
		return err
	}
	_, err = f.WriteString(msg)
	if err != nil {
		return err
	}

	return nil
}

func ReadInfo() string {
	filelock.RLock()
	defer filelock.RUnlock()
	msg, err := ioutil.ReadFile(".amsg")
	if err != nil {
		return ""
	}
	return string(msg)
}

var idfilelock sync.RWMutex

func idUp() (int, error) {
	idfilelock.Lock()
	defer idfilelock.Unlock()
	fid, err := ioutil.ReadFile(".roomid")
	if err != nil {
		return 0, err
	}
	fid = bytes.TrimSpace(fid)
	id, err := strconv.Atoi(string(fid))
	if err != nil {
		return 0, err
	}
	id += 1
	err = ioutil.WriteFile(".roomid", []byte(strconv.Itoa(id)), 0777)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func FilesInit() error {
	fi, err := os.Stat(".roomid")
	if err != nil {
		if err == os.ErrNotExist {
			errs := ioutil.WriteFile(".roomid", []byte(strconv.Itoa(0)), 0777)
			if errs != nil {
				return err
			}
		} else {
			return err
		}

	}
	if fi.IsDir() {
		return fmt.Errorf("not dir")
	}

	mfi, err := os.Stat(".amsg")
	if err != nil {
		if err == os.ErrNotExist {
			errs := ioutil.WriteFile(".amsg", []byte(`管理员信息`), 0777)
			if errs != nil {
				return err
			}
		} else {
			return err
		}

	}
	if mfi.IsDir() {
		return fmt.Errorf("not dir")
	}

	return nil
}
