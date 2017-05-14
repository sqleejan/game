package models

import (
	"fmt"
	"time"
)

type DBRoomPost struct {
	//	Id       int64     `db:"id","primarykey" json:"id"`
	UserId   string    `db:"uid","primarykey" json:"uid"`
	Nicname  string    `db:"nicname" json:"nicname"`
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
	req.RoomName = RoomNamePrefix + post.UserId + string(Krand(8, 0))
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
