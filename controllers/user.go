package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"

	"fmt"

	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title Join Room
// @Description Join Room
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	int		true		"The roomid for user"
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
	rid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
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

// @Title Active User
// @Description Active User
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	int		true		"The roomid for user"
// @Param	uid		query 	string	true		"The uid for user"
// @Param	body		body 	models.UserActiveReq	true		"body for DiverReq"
// @Success 200 {string} token
// @Failure 403 body is empty
// @router /active [post]
func (u *UserController) Active() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}

	rid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}
		var req models.UserActiveReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			err = room.ActiveUser(u.GetString("uid"), &req)
			if err != nil {
				u.CustomAbort(500, err.Error())
				return
			} else {

				u.Data["json"] = "ok"

			}
		}

	}

	u.ServeJSON()
}

// @Title 修改积分
// @Description 修改积分
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	int		true		"The roomid for user"
// @Param	score		query 	int		true		"the score of user"
// @Param	uid			query 	string	true		"the id of user"
// @Success 200 {string} token
// @router /score [post]
func (u *UserController) Score() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}

	score, err := u.GetInt("score", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	rid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) && !room.IsFinace(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}
		err = room.ModifyScore(u.GetString("uid"), score)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
		u.Data["json"] = score
	}

	u.ServeJSON()
}

// @Title Set Assistant
// @Description Set Assistant
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	int		true		"The roomid for user"
// @Param	uid		query 	string	true		"The uid for user"
// @Param	role	query 	int	true		"The role for user"
// @Success 200 {string} set success
// @Failure 403 body is empty
// @router /assistant [post]
func (u *UserController) Assistant() {
	token := u.GetString("token")
	role, err := u.GetInt("role")
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	rid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}

		if err := room.Assistant(u.GetString("uid"), role); err != nil {
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
// @Param	roomid		query 	int		true		"The roomid for user"
// @Success 200 {string} set success
// @Failure 403 uid is null
// @router /list [get]
func (u *UserController) List() {
	rid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
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

	res := struct {
		*auth.MyCustomClaims
		Roles []*models.SelfObj `json:roles`
	}{
		MyCustomClaims: mc,
		Roles:          []*models.SelfObj{},
	}

	nicname, js, err := models.FetchUserInfo(mc.Id)
	if err != nil {
		//u.CustomAbort(500, err.Error())
		fmt.Println(err)

	} else {

		mc.Audience = nicname

		for _, k := range js {
			res.Roles = append(res.Roles, k)
			// room, ok := models.RoomList[k]
			// if !ok {
			// 	delete(js, k)
			// 	continue
			// }
			// if !room.Active() {
			// 	delete(js, k)
			// }
		}

	}

	u.Data["json"] = res
	u.ServeJSON()
}
