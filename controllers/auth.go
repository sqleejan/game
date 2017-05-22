package controllers

import (
	"encoding/json"
	"game/auth"
	"game/models"
	"html/template"
	"time"

	"fmt"

	"strings"

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

// @Title 微信认证
// @Description 微信认证
// @Param	roomid		query 	int		false		"房间ID"
// @Success 200 {string} token
// @router /wx/login [get]
func (o *AuthController) WXAuth() {
	uagent := o.Ctx.Input.Header("User-Agent")
	fmt.Println(uagent)
	rid, err := o.GetInt("roomid", 0)
	if err != nil {
		o.CustomAbort(407, err.Error())
		return
	}

	redirectUrl := auth.CodeUrl(rid, false)
	if !strings.Contains(uagent, "MicroMessenger") {
		redirectUrl = auth.CodeUrl(rid, true)
	}
	o.Redirect(redirectUrl, 302)
	// var ob auth.MyCustomClaims
	// json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	// ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	// o.Data["json"] = ob.Token()
	// o.ServeJSON()
}

// @Title 微信认证
// @Description 微信认证
// @Param	state		query 	int		false		"房间ID"
// @Param	code		query 	string		false		"微信code"
// @Success 200 {string} token
// @router /wx/code [get]
func (o *AuthController) WXCode() {
	roomid, err := o.GetInt("state", 0)
	if err != nil {
		o.CustomAbort(407, err.Error())
		return
	}

	code := o.GetString("code")
	fmt.Println("code=", code)
	if code == "" {
		o.CustomAbort(405, "weixin auth failed!")
		return
	}
	mc, err := auth.WXClaim(code)
	if err != nil {
		o.CustomAbort(405, err.Error())
		return
	}
	_, err = models.CreateDBUser(mc.Id, mc.Audience)
	if err != nil {
		o.CustomAbort(500, err.Error())
		return
	}

	o.Redirect("/fg/redir.html?token="+mc.Token()+fmt.Sprintf("&state=%d", roomid), 302)
	//o.Data["json"] = mc.Token()
	//o.ServeJSON()

	// var ob auth.MyCustomClaims
	// json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	// ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	// o.Data["json"] = ob.Token()
	// o.ServeJSON()
}

// @Title 微信跳转
// @Description 微信跳转
// @Param	state		query 	string		false		"房间ID"
// @Param	token		query 	string		false		"token"
// @Success 200 {string} token
// @router /wx/redirect [get]
func (o *AuthController) Redi() {
	roomid := o.GetString("state")
	token := o.GetString("token")
	if token == "" {
		o.Data["json"] = "ok"
		o.ServeJSON()
		return
	}

	o.Redirect("/fg/index.html?token="+token+"&roomid="+roomid, 302)
	//o.Data["json"] = mc.Token()
	//o.ServeJSON()

	// var ob auth.MyCustomClaims
	// json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	// ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	// o.Data["json"] = ob.Token()
	// o.ServeJSON()
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
		u.CustomAbort(408, err.Error())
		return
	} else if pw != password {
		u.CustomAbort(405, "sercret is wrong")
		return
	}
	var ob auth.MyCustomClaims
	ob.Id = "admin"
	ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	data := struct {
		Token string
	}{
		Token: ob.Token(),
	}

	//u.Data["json"] = "ok"
	//u.Redirect("/swagger/bak/admin.html", 302)
	t := template.Must(template.New("roomreq").Parse(bakhtml))
	t.Execute(u.Ctx.ResponseWriter, data)

	return
}

// @Title 管理员登陆页
// @Description 管理员登陆页
// @Success 200 {string} set success
// @router /admin/login [get]
func (u *AuthController) LoginGet() {
	u.Redirect("/static/login.html", 302)
	return
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
		u.CustomAbort(408, err.Error())
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

// @Title 房间二维码
// @Description 房间二维码
// @Param	roomid		query 	int		false		"房间id"
// @Success 200 {string} image
// @router /qrcode.png [get]
func (o *AuthController) QRCode() {
	roomid, err := o.GetInt("roomid")
	if err != nil {
		o.CustomAbort(500, err.Error())
		return
	}
	bt, err := auth.QRCode(roomid)
	if err != nil {
		o.CustomAbort(500, err.Error())
		return
	}
	o.Ctx.Output.Header("Content-Type", "image/png")
	o.Ctx.Output.Body(bt)
	// var ob auth.MyCustomClaims
	// json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
	// ob.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	// o.Data["json"] = ob.Token()
	// o.ServeJSON()
}
