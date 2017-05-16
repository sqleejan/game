package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["game/controllers:AuthController"] = append(beego.GlobalControllerRouter["game/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Token",
			Router: `/token`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:AuthController"] = append(beego.GlobalControllerRouter["game/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Create",
			Router: `/user/create`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:AuthController"] = append(beego.GlobalControllerRouter["game/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/admin/login`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:AuthController"] = append(beego.GlobalControllerRouter["game/controllers:AuthController"],
		beego.ControllerComments{
			Method: "LoginGet",
			Router: `/admin/login`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:AuthController"] = append(beego.GlobalControllerRouter["game/controllers:AuthController"],
		beego.ControllerComments{
			Method: "Update",
			Router: `/admin/update`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/create`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Config",
			Router: `/config`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Master",
			Router: `/master`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Discard",
			Router: `/discard`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Distribute",
			Router: `/send`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Gain",
			Router: `/gain`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:ObjectController"] = append(beego.GlobalControllerRouter["game/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:ObjectController"] = append(beego.GlobalControllerRouter["game/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:ObjectController"] = append(beego.GlobalControllerRouter["game/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:ObjectController"] = append(beego.GlobalControllerRouter["game/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:ObjectController"] = append(beego.GlobalControllerRouter["game/controllers:ObjectController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Config",
			Router: `/config`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:roomid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "DeleteRoom",
			Router: `/:roomid`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Bill",
			Router: `/:roomid/bill`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "ListDB",
			Router: `/listdb`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Request",
			Router: `/request`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "RequestDelete",
			Router: `/request/:id`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "RequestAccept",
			Router: `/request/:id`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "RequestList",
			Router: `/request/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "Join",
			Router: `/join`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "Active",
			Router: `/active`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "Score",
			Router: `/score`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "Assistant",
			Router: `/assistant`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetSelf",
			Router: `/self`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
