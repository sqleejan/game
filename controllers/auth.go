package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"
	"time"

	"github.com/astaxie/beego"
)

// Operations about auth
type AuthController struct {
	beego.Controller
}

// @Title Token
// @Description create token
// @Param	body		body 	models.TmpClaims	true		"The object content"
// @Success 200 {string} token
// @Failure 403 body is empty
// @router /token [post]
func (o *AuthController) Token() {
	var ob auth.MyCustomClaims
	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	o.Data["json"] = ob.Token()
	o.ServeJSON()
}

// @Title 临时创建用户
// @Description 临时创建用户
// @Param	uid		query 	string	true		"The uid for user"
// @Param	nicname	query 	string	true		"The nicname for user"
// @Success 200 {string} set success
// @Failure 403 uid is null
// @router /user/create [post]
func (u *AuthController) Create() {
	uid := u.GetString("uid")
	nicname := u.GetString("nicname")
	_, err := models.CreateDBUser(uid, nicname)
	if err != nil {
		u.CustomAbort(500, err.Error())
		return
	}
	u.Data["json"] = "ok"

	u.ServeJSON()
}

// @Title 管理员登陆
// @Description 管理员登陆
// @Success 200 {string} set success
// @router /admin/login [post]
func (u *AuthController) Login() {
	u.Ctx.Request.ParseForm()
	input := u.Input()
	username := input.Get("username")
	password := input.Get("password")
	if pw, err := models.GetPassword(username); err != nil {
		u.CustomAbort(405, err.Error())
		return
	} else if pw != password {
		u.CustomAbort(405, "sercret is wrong")
		return
	}
	var ob auth.MyCustomClaims
	ob.Id = "admin"
	ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	u.Data["json"] = ob.Token()
	//u.Data["json"] = "ok"

	u.ServeJSON()
}

// @Title 管理员密码更新
// @Description 管理员密码更新
// @Success 200 {string} set success
// @router /admin/update [post]
func (u *AuthController) Update() {
	u.Ctx.Request.ParseForm()
	input := u.Input()
	username := input.Get("username")
	password := input.Get("password")
	if pw, err := models.GetPassword(username); err != nil {
		u.CustomAbort(405, err.Error())
		return
	} else if pw != password {
		u.CustomAbort(405, "sercret is wrong")
		return
	}
	var ob auth.MyCustomClaims
	ob.Id = "admin"
	ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	u.Data["json"] = ob.Token()
	//u.Data["json"] = "ok"

	u.ServeJSON()
}


