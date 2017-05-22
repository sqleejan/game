package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"

	"fmt"

	"github.com/astaxie/beego"
)

// Operations about Rooms
type RoomController struct {
	beego.Controller
}

// @Title CreateRoom
// @Description create room
// @Param	body		body 	models.TmpRoomReq	true		"body for room content"
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

	var req models.RoomReq
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &req)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
		//u.Data["json"] = err.Error()
	} else {
		if mc.Id != req.UserId {
			u.CustomAbort(408, "permission is not allow!")
			return
		}

		print(req)
		room, err := models.CreateRoom(&req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
		u.Data["json"] = room.Convert()
	}

	u.ServeJSON()
}

// @Title ConfigRoom
// @Description config room
// @Param	body		body 	models.TmpRoomConfig	true		"body for room content"
// @Param	roomid		query 	int		true		"The id for room"
// @Param	token		query 	string	true		"The token for user"
// @Success 200 {string} ok
// @Failure 403 body is empty
// @router /config [post]
func (u *RoomController) Config() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	var req models.RoomConfig
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &req)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
		//u.Data["json"] = err.Error()
	}
	print(req)
	roomid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		}
		if !room.IsAdmin(mc.Id) {
			u.CustomAbort(408, "permission is not allow!")
			return
		}
		err = room.Config(&req)
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
		u.Data["json"] = "ok"

	}

	u.ServeJSON()
}

// @Title GetAll
// @Description get all Rooms
// @Param	token		query 	string	true		"The token for user"
// @Param	limit		query 	int		false		"The default is 20"
// @Param	page		query 	int		false		"The default is 1"
// @Param	body		body 	models.ListReq	true		"body for list content"
// @Success 200 {object} models.TmpRespone
// @router /list [post]
func (u *RoomController) GetAll() {
	token := u.GetString("token")
	_, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	// if mc.Id != "admin" {
	// 	u.CustomAbort(405, "permission is not allow!")
	// 	return
	// }

	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 10)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	var lrq models.ListReq
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &lrq)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
		//u.Data["json"] = err.Error()
	} else {
		fmt.Println("len(roomlist)=", len(models.RoomList))
		u.Data["json"] = models.RLConvert(models.RoomList).Convert(&lrq, page, limit)
	}
	// rooms := []interface{}{}
	// for _, r := range models.RoomList {
	// 	if r.Active() {
	// 		rooms = append(rooms, r.Convert())
	// 	}
	// }

	u.ServeJSON()
}

// @Title Get
// @Description get user by roomid
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		path 	int		true		"The key for staticblock"
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
	roomid, err := u.GetInt(":roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			if !room.IsAnyone(mc.Id) && mc.Id != "admin" {
				u.CustomAbort(408, "permission is not allow!")
				return
			}
			print(room.Convert())
			u.Data["json"] = room.Convert()
		}
	}
	u.ServeJSON()
}

// @Title 激活
// @Description 激活
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		query 	int		true		"The id for room"
// @Success 200 {bool} success
// @router /active [get]
func (u *RoomController) SetActive() {
	fmt.Println("-----------------------------")
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	roomid, err := u.GetInt("roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	fmt.Println("roomid=", roomid)
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			if mc.Id != "admin" {
				u.CustomAbort(408, "permission is not allow!")
				return
			}
			//print(room.Convert())
			room.SetActive()
			u.Data["json"] = room.Convert()
		}
	}
	u.ServeJSON()
}

// @Title 续房间
// @Description 续房间
// @Param	token		query 	string	true		"The token for user"
// @Param	duration	query 	int		true		"房间延续时间"
// @Param	roomid		path 	string	true		"The id for room"
// @Success 200 {string} ok
// @Failure 403 :roomid is empty
// @router /:roomid/renew [post]
func (u *RoomController) Renew() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	dur, err := u.GetInt("duration", 1)
	if err != nil {
		u.CustomAbort(407, "duration format is wrong!")
		return
	}
	roomid, err := u.GetInt(":roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			if mc.Id != "admin" {
				u.CustomAbort(408, "permission is not allow!")
				return
			}
			room.Renew(dur)
			u.Data["json"] = "ok"
		}
	}
	u.ServeJSON()
}

// @Title Delete Room
// @Description delete room
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		path 	int		true		"The roomid "
// @Success 200 {string} ok
// @Failure 403 :roomid is empty
// @router /:roomid [delete]
func (u *RoomController) DeleteRoom() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(408, "permission is not allow!")
		return
	}
	roomid, err := u.GetInt(":roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			room.Close()
			delete(models.RoomList, roomid)
			fmt.Println(room.DeleteDB())
			u.Data["json"] = room.Convert()
		}
	}
	u.ServeJSON()
}

// @Title Cancle Room
// @Description Cancle room
// @Param	token		query 	string	true		"The token for user"
// @Param	roomid		path 	int		true		"The roomid "
// @Success 200 {string} ok
// @Failure 403 :roomid is empty
// @router /:roomid/cancle [post]
func (u *RoomController) Cancle() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(408, "permission is not allow!")
		return
	}
	roomid, err := u.GetInt(":roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok {
			u.CustomAbort(500, "the room is not exist")
			return
		} else {
			room.Cancle()
			u.Data["json"] = "ok"
		}
	}
	u.ServeJSON()
}

// @Title Bill
// @Description get bill for roomid
// @Param	token		query 	string	true		"The token for user"
// @Param	limit		query 	int		false		"The default is 20"
// @Param	page		query 	int		false		"The default is 1"
// @Param	roomid		path 	int		true		"The key for staticblock"
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
	roomid, err := u.GetInt(":roomid", 0)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 20)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	if roomid != 0 {
		room, ok := models.RoomList[roomid]
		if !ok && mc.Id != "admin" {
			u.CustomAbort(500, "the room is not exist")
			return
		}
		if !(room != nil && room.IsAnyone(mc.Id)) && mc.Id != "admin" {
			u.CustomAbort(408, "permission is not allow!")
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
		u.CustomAbort(408, "permission is not allow!")
		return
	}
	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 20)
	if err != nil {
		u.CustomAbort(407, err.Error())
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

// @Title 房间申请
// @Description 申请房间
// @Param	token		query 	string	true		"The token for user"
// @Param	body		body 	models.DBRoomPost	true		"body for room content"
// @Success 200 {string} name
// @router /request [post]
func (u *RoomController) Request() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	var req models.DBRoomPost
	err = json.Unmarshal(u.Ctx.Input.RequestBody, &req)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
		//u.Data["json"] = err.Error()
	} else {
		req.UserId = mc.Id
		req.RoomName = models.GenerateName(mc.Id)
		err = req.Insert()
		if err != nil {
			u.CustomAbort(500, err.Error())
			return
		}
	}

	u.Data["json"] = req.RoomName

	u.ServeJSON()
}

// @Title 房间申请删除
// @Description 删除房间申请
// @Param	token		query 	string	true		"The token for user"
// @Param	id			path 	string	true		"id for room request"
// @Success 200 {string} ok
// @router /request/:id [delete]
func (u *RoomController) RequestDelete() {
	token := u.GetString("token")
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(408, "permission is not allow!")
		return
	}

	del := models.DBRoomPost{}
	del.UserId = u.GetString(":id")
	err = del.Delete()
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
	}

	u.Data["json"] = "ok"

	u.ServeJSON()
}

// @Title 接受房间申请
// @Description 接受房间申请
// @Param	token		query 	string	true		"The token for user"
// @Param	id			path 	string	true		"id for room request"
// @Success 200 {string} ok
// @router /request/:id [post]
func (u *RoomController) RequestAccept() {
	token := u.GetString("token")
	//fmt.Println("id=", u.GetString(":id"))
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(408, "permission is not allow!")
		return
	}

	accept := models.DBRoomPost{}
	accept.UserId = u.GetString(":id")

	err = accept.Fetch()

	if err != nil {
		u.CustomAbort(408, err.Error())
		return
	}

	err = accept.CreateRoom()
	if err != nil {
		u.CustomAbort(406, err.Error())
		return
	}

	u.Data["json"] = "ok"

	u.ServeJSON()
}

// @Title 房间申请列表
// @Description 房间申请列表
// @Param	token		query 	string	true		"The token for user"
// @Param	limit		query 	int		false		"The default is 20"
// @Param	page		query 	int		false		"The default is 1"
// @Success 200 {string} ok
// @router /request/list [get]
func (u *RoomController) RequestList() {
	token := u.GetString("token")
	fmt.Println("token:", token)
	mc, err := auth.Parse(token)
	if err != nil {
		u.CustomAbort(405, err.Error())
		return
	}
	if mc.Id != "admin" {
		u.CustomAbort(408, "permission is not allow!")
		return
	}

	page, err := u.GetInt("page", 1)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}
	limit, err := u.GetInt("limit", 10)
	if err != nil {
		u.CustomAbort(407, err.Error())
		return
	}

	body, err := models.ListRoomPosts(page, limit)

	if err != nil {
		u.CustomAbort(500, err.Error())
		return
	}

	u.Data["json"] = body

	u.ServeJSON()
}

func print(v interface{}) {
	if e, ok := v.(error); ok {
		fmt.Println(e.Error())
	}
	return
	bts, _ := json.MarshalIndent(v, "", " ")
	fmt.Println(string(bts))
}
