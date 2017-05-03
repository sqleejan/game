package models

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	cli, err := newEm()
	if err != nil {
		t.Error(err)
		fmt.Println(err)
		return
	}

	//fmt.Println(cli.AddGroup("test1", "", "leejan", true, false, 10, nil))
	// fmt.Println(cli.CreateAccount("leejan", "123456", "Lee"))
	// fmt.Println(cli.AddUserToGroup("14915347742721", "leejan"))
	//fmt.Println(cli.FetchUserFromGroup("15078152798209"))
	//fmt.Println(cli.FetchAllGroupFromApp())
	fmt.Println(cli.GetUserToken("leejan", "123456"))
}
