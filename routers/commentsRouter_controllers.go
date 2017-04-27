package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/create`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:CattleController"] = append(beego.GlobalControllerRouter["game/controllers:CattleController"],
		beego.ControllerComments{
			Method: "Distribute",
			Router: `/send`,
			AllowHTTPMethods: []string{"get"},
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
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:RoomController"] = append(beego.GlobalControllerRouter["game/controllers:RoomController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:roomid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["game/controllers:UserController"] = append(beego.GlobalControllerRouter["game/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
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
			Method: "GetUser",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
