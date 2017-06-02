package models

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

const (
	defautAdminID        = "admin"
	defaultAdminPassword = "caozihua"
)

func GetPassword(username string) (string, error) {
	u := &DBUser{
		Id: username,
	}
	err := u.Fetch(dBEngine)
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

func ModifyPassword(username, password string) error {
	u := &DBUser{
		Id: username,
	}
	err := u.Fetch(dBEngine)
	if err != nil {
		return err
	}
	u.Password = password
	err = u.Update(dBEngine)
	if err != nil {
		return err
	}
	return err
}

func adminInsert() {
	u := &DBUser{
		Id:       defautAdminID,
		Password: defaultAdminPassword,
	}
	u.Insert(dBEngine)
	ufly := &DBUser{
		Id:       "admin_fly",
		Password: "aniewuli",
	}
	ufly.Insert(dBEngine)
}

type RLConvert map[int]*Room

func (room *Room) StatByte() byte {
	stat := byte(0)
	if room.reqStatus {
		stat = stat | 1<<4
	}
	now := time.Now()
	if room.active {
		stat = stat | 1<<3
	} else {
		stat = stat | 1
	}

	isused := room.endTime.After(now)
	expired := false
	if isused {
		stat = stat | 1<<2
		expired = now.Add(time.Minute * 30).After(room.endTime)
		if expired {
			stat = stat | 1<<1
		}
	}
	return stat

}

func (lreq *ListReq) StatByte() byte {
	stat := byte(0)
	if lreq.Req {
		stat = stat | 1<<4
	}
	if lreq.Active {
		stat = stat | 1<<3
	}
	if lreq.IsUsed {
		stat = stat | 1<<2
	}
	if lreq.Expire {
		stat = stat | 1<<1
	}
	if lreq.Cancle {
		stat = stat | 1
	}
	return stat
}

type intKeyPair struct {
	key   int
	value int
}

type sortSlice []*intKeyPair

func (s sortSlice) Less(i, j int) bool {
	return s[i].value < s[j].value
}

func (s sortSlice) Len() int {
	return len(s)
}

func (s sortSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (rl RLConvert) Convert(listreq *ListReq, page int, size int) interface{} {
	fmt.Println(len(rl))
	leftList := map[int]*Room{}
	query := ".*"
	if listreq.Like != "" {
		query = listreq.Like + query
	}
	for id, room := range rl {

		matched, _ := regexp.MatchString(listreq.Like, room.name)
		fmt.Println(room.name, listreq.Like, matched)
		if matched {
			reqstat := listreq.StatByte()
			fmt.Println(reqstat, room.StatByte())
			if (reqstat & room.StatByte()) == reqstat {
				//if room.reqStatus
				leftList[id] = rl[id]
			}

		}
	}

	list := sortSlice{}
	for k := range leftList {
		//ros.CreateAt.Unix()
		index := k
		list = append(list, &intKeyPair{index, index % 8000})
	}
	sort.Sort(list)
	resp := &struct {
		Pagination
		Data []*RoomResponeNoUsers `json:"data"`
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
	data := []*RoomResponeNoUsers{}
	for _, v := range list[start:end] {
		data = append(data, leftList[v.key].ConvertNoUsers())
	}
	resp.Data = data
	return resp
}
