package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about Rooms
type RoomController struct {
	beego.Controller
}

// @Title CreateRoom
// @Description create room
// @Param	body		body 	models.RoomReq	true		"body for room content"
// @Param	token		query 	string	true		"The token for user"
// @Success 200 {object} models.TmpRespone
// @Failure 403 body is empty
// @router / [post]
func (u *RoomController) Post() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(405, "permission is not allow!")
		return
	}

	var req models.RoomReq
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &req)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
		//u.Data["json"] = err.Error()
	} else {
		room, err := models.CreateRoom(&req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
		u.Data["json"] = room.Convert()
	}

	u.ServeJSON()
}

// @Title GetAll
// @Description get all Rooms
// @Param	token		query 	string	true		"The token for user"
// @Success 200 {object} models.TmpRespone
// @router /list [get]
func (u *RoomController) GetAll() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(405, "permission is not allow!")
		return
	}
	rooms := []interface{}{}
	for _, r := range models.RoomList {
		rooms = append(rooms, r.Convert())
	}
	u.Data["json"] = rooms
	u.ServeJSON()
}

// @Title Get
// @Description get user by roomid
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.TmpRespone
// @Failure 403 :roomid is empty
// @router /:roomid [get]
func (u *RoomController) Get() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	roomid := u.GetString(":roomid")
	if roomid != "" {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			if !room.IsAnyone(mc.Id) && mc.Id != "admin" {
				u.CustomAbort(405, "permission is not allow!")
				return
			}
			u.Data["json"] = room.Convert()
		}
	}
	u.ServeJSON()
}

// @Title Bill
// @Description get bill for roomid
// @Param	token		query 	string	true		"The token for user"
// @Param	limit		query 	int		false		"The default is 20"
// @Param	page		query 	int		false		"The default is 1"
// @Param	roomid		path 	string	true		"The key for staticblock"
// @Success 200 {string} success
// @Failure 403 :roomid is empty
// @router /:roomid/bill [get]
func (u *RoomController) Bill() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	roomid := u.GetString(":roomid")
	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 20)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if roomid != "" {
		room, ok := models.RoomList[roomid]
		if !ok && mc.Id != "admin" {
			u.CustomAbort(500, "the room is not exist")
			return
		}
		if !(room != nil && room.IsAnyone(mc.Id)) && mc.Id != "admin" {
			u.CustomAbort(405, "permission is not allow!")
			return
		}
		b, err := models.BillList(roomid, page, limit)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
		u.Data["json"] = b

	}
	u.ServeJSON()
}

// @Title Listdb
// @Description get room Listdb
// @Param	token		query 	string	true		"The token for user"
// @Param	limit		query 	int		false		"The default is 20"
// @Param	page		query 	int		false		"The default is 1"
// @Success 200 {string} success
// @router /listdb [get]
func (u *RoomController) ListDB() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(405, "permission is not allow!")
		return
	}
	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 20)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}

	ob := &models.DBRecord{}
	list, err := ob.ListDBRoom(page, limit)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
	}
	u.Data["json"] = list

	u.ServeJSON()
}
