package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"

	"fmt"

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
// @Success 200 {string} success
// @router /create [post]
func (u *CattleController) Post() {
	print("action:qizhuang")
	token := u.GetString("token")
	rid := u.GetString("roomid")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	// if mc.Id != "admin" {
	// 	u.CustomAbort(405, "permission is not allow!")
	// 	return
	// }

	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAdmin(mc.Id) && !room.IsAssistant(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}

		if err := room.SendRedhat(); err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			u.Data["json"] = "ok"
		}

	}

	u.ServeJSON()
}

// @Title 占庄
// @Description 占庄
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Success 200 {string} success
// @router /keepz [get]
func (u *CattleController) Keep() {
	print("action:keepzhuang")
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
		if err := room.KeepZhuang(mc.Id); err != nil {
			u.CustomAbort(409, err.Error())
			return
		}
	}
	u.Data["json"] = "ok"
	u.ServeJSON()
}

// @Title 配置庄
// @Description config redhat
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Param	body		body 	models.RedReq	true		"body for rancher"
// @Param	cancel		query 	bool	true		"cancle for rancher"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /config [post]
func (u *CattleController) Config() {
	print("action:configzhuang")
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
	cancel, err := u.GetBool("cancel", false)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsAssistant(mc.Id) && !room.IsAdmin(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}
		var req models.RedReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			if err := room.ConfigRedhat(&req, cancel); err != nil {
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
	print(fmt.Sprintf("action:qiangzhuang %v", mc.Id))
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
		if !room.IsCustom(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
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

// @Title  弃庄
// @Description 弃庄
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for rancher"
// @Success 200 {string} success
// @router /discard [post]
func (u *CattleController) Discard() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	print(fmt.Sprintf("action:fangqizhuang %v", mc.Id))
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
		if !room.IsAnyone(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}

		if err := room.Discard(); err != nil {
			u.CustomAbort(500, err.Error())
			return
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
// @Param	body		body 	models.DiverReq	true		"body for DiverReq"
// @Success 200 {body} models.Marks
// @Failure 403 query is empty
// @router /send [post]
func (u *CattleController) Distribute() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	print(fmt.Sprintf("action:fahongbao %v", mc.Id))
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {

		var req models.DiverReq
		err := json.Unmarshal(u.Ctx.Input.RequestBody, &req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {
			print(req)
			rs, err := room.Diver(mc.Id, &req)
			fmt.Println("distribute socre....")
			print(rs)
			print(err)
			if err != nil {
				u.CustomAbort(500, err.Error())
				return
			} else {
				//print(rs)
				u.Data["json"] = rs

			}
		}

	}

	u.ServeJSON()
}

// @Title Gain cattle
// @Description gain cattle
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	string	true		"The roomid for cattle"
// @Success 200 {object} models.ScoreUnion
// @Failure 403 query is empty
// @router /gain [get]
func (u *CattleController) Gain() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}

	print(fmt.Sprintf("action:qianghongbao %v", mc.Id))
	rid := u.GetString("roomid")
	room, ok := models.RoomList[rid]
	if !ok {
		u.CustomAbort(500, "the room is not exist")
		return
	} else {
		if !room.IsCustom(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}

		score, err := room.GetScore(mc.Id)
		fmt.Println("gain socre....")
		print(score)
		print(err)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		} else {

			u.Data["json"] = score

		}
	}

	u.ServeJSON()
}
