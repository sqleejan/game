package models

import (
	"fmt"
	"testing"

	"encoding/json"
)

func TestClient(t *testing.T) {
	ch := make(chan int, 5)
	GenerateScore(5, ch)
	rs1 := []*result{}
	for i := 0; i < 5; i++ {
		score := <-ch
		//fmt.Println(score)
		rs1 = append(rs1, &result{score: score})
	}
	juge(rs1, 10, 10)
	reps := MakeReport(rs1)
	bb, _ := json.Marshal(reps)
	fmt.Println(string(bb))
	// for _, v := range reps {
	// 	fmt.Println(v.Score, v.Pay)
	// }
	return
	cli, err := newEm()
	if err != nil {
		t.Error(err)
		fmt.Println(err)
		return
	}

	fmt.Println(CreateDBUser("openid3", "cc"))
	return

	//fmt.Println(cli.AddGroup("test1", "", "leejan", true, false, 10, nil))
	// fmt.Println(cli.CreateAccount("leejan", "123456", "Lee"))
	//fmt.Println(cli.AddUserToGroup("14915347742721", "leejan"))
	//fmt.Println(cli.FetchUserFromGroup("15078152798209"))
	//fmt.Println(cli.FetchAllGroupFromApp())
	//fmt.Println(cli.GetUserToken("leejan", "123456"))
	//fmt.Println(GetToken("test"))
	//fmt.Println(Record("room1", `{"Name":"leejan"}`))
	list, _ := cli.FetchAllGroupFromApp()
	for _, v := range list.Data {
		fmt.Println(v.Groupname)
	}
	req := RoomReq{
		Duration:  10,
		UserId:    "leejan",
		UserLimit: 100,
		RoomName:  "yoyo_test1",
	}
	CreateRoom(&req)
	fmt.Println("-----------------------------------")
	list, _ = cli.FetchAllGroupFromApp()
	for _, v := range list.Data {
		fmt.Println(v.Groupname)
	}
	//fmt.Println(cli.FetchAllGroupFromApp())
}
