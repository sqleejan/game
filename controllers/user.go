package controllers

import (
	"encoding/json"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title Join Room
// @Description Join Room
// @Param	roomid		query 	string	true		"The roomid for user"
// @Param	body		body 	models.UserReq	true		"body for user content"
// @Success 200 {object} models.UserReq
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.Data["json"] = "the room is not exist"
	} else {
		var req models.RoomReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			if err := room.AppendUser(models.UserReq{req.UserId, req.Username}); err != nil {
				u.Data["json"] = err.Error()
			} else {
				u.Data["json"] = "ok"
			}
		}
	}

	u.ServeJSON()
}

// @Title Set Assistant
// @Description Set Assistant
// @Param	roomid		query 	string	true		"The roomid for user"
// @Param	body		body 	models.UserReq	true		"body for user content"
// @Success 200 {string} set success
// @Failure 403 body is empty
// @router /assistant [post]
func (u *UserController) Assistant() {
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.Data["json"] = "the room is not exist"
	} else {
		var req models.UserReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			if err := room.Assistant(req.UserId); err != nil {
				u.Data["json"] = err.Error()
			} else {
				u.Data["json"] = "ok"
			}
		}
	}

	u.ServeJSON()
}

// @Title Get User
// @Description Get User
// @Param	uid		path 	string	true		"The uid for user"
// @Success 200 {string} set success
// @Failure 403 uid is null
// @router /:uid [get]
func (u *UserController) GetUser() {
	uid := u.GetString("uid")
	user, ok := models.UserList[uid]
	if !ok {
		u.Data["json"] = "the user is not exist"
	} else {
		u.Data["json"] = user
	}
	u.ServeJSON()
}
