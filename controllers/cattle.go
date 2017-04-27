package controllers

import (
	"encoding/json"
	"game/models"

	"github.com/astaxie/beego"
)

// Operations about Cattle
type CattleController struct {
	beego.Controller
}

// @Title Rancher
// @Description set rancher
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Param	body		body 	models.RedReq	true		"body for rancher"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /create [post]
func (u *CattleController) Post() {
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.Data["json"] = "the room is not exist"
	} else {
		var req models.RedReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			if err := room.SendRedhat(&req); err != nil {
				u.Data["json"] = err.Error()
			} else {
				u.Data["json"] = "ok"
			}
		}
	}

	u.ServeJSON()
}

// @Title Distribute
// @Description Distribute cattle
// @Param	roomid		query 	string	true		"The roomid for cattle"
// @Param	rancher		query 	string	true		"the rancher for room"
// @Success 200 {body} models.Mark
// @Failure 403 query is empty
// @router /send [get]
func (u *CattleController) Distribute() {
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.Data["json"] = "the room is not exist"
	} else {
		rs, err := room.Diver(u.GetString("rancher"))
		if err != nil {
			u.Data["json"] = err.Error()
		} else {

			u.Data["json"] = models.MakeReport(rs)

		}
	}

	u.ServeJSON()
}

// @Title Gain cattle
// @Description gain cattle
// @Param	roomid		query 	string	true		"The roomid for cattle"
// @Param	custom		query 	string	true		"the custom for room"
// @Success 200 {int} score number
// @Failure 403 query is empty
// @router /gain [get]
func (u *CattleController) Gain() {
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.Data["json"] = "the room is not exist"
	} else {
		score, err := room.GetScore(u.GetString("custom"))
		if err != nil {
			u.Data["json"] = err.Error()
		} else {

			u.Data["json"] = score

		}
	}

	u.ServeJSON()
}
