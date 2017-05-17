package controllers

import (
	"game/auth"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about super
type FlyController struct {
	beego.Controller
}

// @Title 启动
// @Description 启动
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The id for room"
// @Success 200 {string} ok
// @router /start [post]
func (f *FlyController) Start() {
	token := f.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		f.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		f.CustomAbort(408, "permission is not allow!")
		return
	}
	rid := f.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		f.CustomAbort(500, "the room is not exist")
		return
	} else {

		err = room.Super(true)
		if err != nil {
			f.CustomAbort(500, err.Error())
			return
		}

	}
	f.Data["json"] = "ok"
	f.ServeJSON()
}

// @Title 停止
// @Description 停止
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The id for room"
// @Success 200 {string} ok
// @router /stop [post]
func (f *FlyController) Stop() {
	token := f.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		f.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		f.CustomAbort(408, "permission is not allow!")
		return
	}
	rid := f.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		f.CustomAbort(500, "the room is not exist")
		return
	} else {

		err = room.Super(false)
		if err != nil {
			f.CustomAbort(500, err.Error())
			return
		}

	}
	f.Data["json"] = "ok"
	f.ServeJSON()
}
