package models

import (
	"sort"
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
}

type RLConvert map[string]*Room

func (rl RLConvert) Convert(page, size int) interface{} {
	list := []string{}
	for k := range rl {
		list = append(list, k)
	}
	sort.Strings(list)
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
		data = append(data, rl[v].ConvertNoUsers())
	}
	resp.Data = data
	return resp
}
