package controllers

import (
	"game/auth"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title Join Room
// @Description Join Room
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for user"
// @Success 200 {string} token
// @Failure 403 body is empty
// @router /join [post]
func (u *UserController) Join() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {

		if hxtoken, err := room.AppendUser(mc.Id, mc.Audience); err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			u.Data["json"] = hxtoken
		}

	}

	u.ServeJSON()
}

// @Title Set Assistant
// @Description Set Assistant
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for user"
// @Param	uid		query 	string	true		"The uid for user"
// @Success 200 {string} set success
// @Failure 403 body is empty
// @router /assistant [post]
func (u *UserController) Assistant() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) {
			u.CustomAbort(405, "permission is not allow!")
			return
		}

		if err := room.Assistant(u.GetString("uid")); err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			u.Data["json"] = "ok"
		}

	}

	u.ServeJSON()
}

// @Title 玩家列表
// @Description Get Users
// @Param	token		query 	string	false		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for user"
// @Success 200 {string} set success
// @Failure 403 uid is null
// @router /list [get]
func (u *UserController) List() {
	rid := u.GetString("roomid")
	_, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	}
	u.Data["json"] = ""

	u.ServeJSON()
}

// @Title Self
// @Description Get Self
// @Param	token		query 	string	true		"The token for user"
// @Success 200 {object} auth.MyCustomClaims
// @router /self [get]
func (u *UserController) GetSelf() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	u.Data["json"] = mc
	u.ServeJSON()
}
