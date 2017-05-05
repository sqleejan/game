package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about Cattle
type CattleController struct {
	beego.Controller
}

// @Title 起庄
// @Description set rancher
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Param	body		body 	models.RedReq	true		"body for rancher"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /create [post]
func (u *CattleController) Post() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	// if mc.Id != "admin" {
	// 	u.CustomAbort(405, "permission is not allow!")
	// 	return
	// }
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) && !room.IsAssistant(mc.Id) {
			u.CustomAbort(405, "permission is not allow!")
			return
		}
		var req models.RedReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			if err := room.SendRedhat(&req); err != nil {
				u.CustomAbort(500, err.Error())
				return
			} else {
				u.Data["json"] = "ok"
			}
		}
	}

	u.ServeJSON()
}

// @Title 抢庄
// @Description fetch rancher
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /master [post]
func (u *CattleController) Master() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	// if mc.Id != "admin" {
	// 	u.CustomAbort(405, "permission is not allow!")
	// 	return
	// }
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if room.IsAdmin(mc.Id) || room.IsAssistant(mc.Id) {
			u.CustomAbort(405, "permission is not allow!")
			return
		}
		if err := room.MasterRedhat(mc.Id); err != nil {
			u.Data["json"] = "fail"
		} else {
			u.Data["json"] = "ok"
		}

	}

	u.ServeJSON()
}

// @Title Distribute
// @Description Distribute cattle
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for cattle"
// @Success 200 {body} models.Marks
// @Failure 403 query is empty
// @router /send [get]
func (u *CattleController) Distribute() {
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
		rs, err := room.Diver(mc.Id)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {

			u.Data["json"] = rs

		}
	}

	u.ServeJSON()
}

// @Title Gain cattle
// @Description gain cattle
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for cattle"
// @Success 200 {int} score number
// @Failure 403 query is empty
// @router /gain [get]
func (u *CattleController) Gain() {
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
		if !room.IsCustom(mc.Id) {
			u.CustomAbort(405, "permission is not allow!")
			return
		}
		score, err := room.GetScore(mc.Id)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {

			u.Data["json"] = float32(score) / 100

		}
	}

	u.ServeJSON()
}
