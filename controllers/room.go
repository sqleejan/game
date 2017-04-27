package controllers

import (
	"encoding/json"
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
// @Success 200 {object} models.RoomRespone
// @Failure 403 body is empty
// @router / [post]
func (u *RoomController) Post() {
	var req models.RoomReq
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		room := models.CreateRoom(&req)
		u.Data["json"] = room.Convert()
	}

	u.ServeJSON()
}

// @Title GetAll
// @Description get all Rooms
// @Success 200 {object} models.RoomRespone
// @router / [get]
func (u *RoomController) GetAll() {
	rooms := []interface{}{}
	for _, r := range models.RoomList {
		rooms = append(rooms, r.Convert())
	}
	u.Data["json"] = rooms
	u.ServeJSON()
}

// @Title Get
// @Description get user by roomid
// @Param	roomid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.RoomRespone
// @Failure 403 :roomid is empty
// @router /:roomid [get]
func (u *RoomController) Get() {
	roomid := u.GetString(":roomid")
	if roomid != "" {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.Data["json"] = "the room in not exsit"
		} else {
			u.Data["json"] = room.Convert()
		}
	}
	u.ServeJSON()
}

//// @Title Update
//// @Description update the user
//// @Param	uid		path 	string	true		"The uid you want to update"
//// @Param	body		body 	models.User	true		"body for user content"
//// @Success 200 {object} models.User
//// @Failure 403 :uid is not int
//// @router /:uid [put]
//func (u *UserController) Put() {
//	uid := u.GetString(":uid")
//	if uid != "" {
//		var user models.User
//		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
//		uu, err := models.UpdateUser(uid, &user)
//		if err != nil {
//			u.Data["json"] = err.Error()
//		} else {
//			u.Data["json"] = uu
//		}
//	}
//	u.ServeJSON()
//}

//// @Title Delete
//// @Description delete the user
//// @Param	uid		path 	string	true		"The uid you want to delete"
//// @Success 200 {string} delete success!
//// @Failure 403 uid is empty
//// @router /:uid [delete]
//func (u *UserController) Delete() {
//	uid := u.GetString(":uid")
//	models.DeleteUser(uid)
//	u.Data["json"] = "delete success!"
//	u.ServeJSON()
//}

//// @Title Login
//// @Description Logs user into the system
//// @Param	username		query 	string	true		"The username for login"
//// @Param	password		query 	string	true		"The password for login"
//// @Success 200 {string} login success
//// @Failure 403 user not exist
//// @router /login [get]
//func (u *UserController) Login() {
//	username := u.GetString("username")
//	password := u.GetString("password")
//	if models.Login(username, password) {
//		u.Data["json"] = "login success"
//	} else {
//		u.Data["json"] = "user not exist"
//	}
//	u.ServeJSON()
//}

//// @Title logout
//// @Description Logs out current logged in user session
//// @Success 200 {string} logout success
//// @router /logout [get]
//func (u *UserController) Logout() {
//	u.Data["json"] = "logout success"
//	u.ServeJSON()
//}
