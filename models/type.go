package models

import "time"

type ListReq struct {
	Like   string `json:"like"`
	Req    bool   `json:"req"`
	IsUsed bool   `json:"is_used"`
	Expire bool   `json:"expire"`
	Active bool   `json:"active"`
	Cancle bool   `json:"cancle"`
}

type TmpRoomReq struct {
	Duration     int
	UserId       string
	UserLimit    int
	RoomName     string
	Water        int
	Base         int
	Nickname     string
	ScoreLimit   int
	CountUp      int
	CountDown    int
	RedDown      float32
	RedUp        float32
	RedCountDown int
	RedCountUp   int
	Timeout      int
	Describe     string
	RedInterval  int
}

type TmpRoomConfig struct {
	Base         int
	Water        int
	CountUp      int
	CountDown    int
	RedDown      float32
	RedUp        float32
	RedCountDown int
	RedCountUp   int
	Timeout      int
	ScoreLimit   int
	Describe     string
	RedInterval  int
}

type TmpRespone struct {
	RoomId    string
	RoomName  string
	Base      int
	Water     int
	Admin     string
	Banker    string
	CreateAt  time.Time
	StartTime time.Time
	EndTime   time.Time
	ActiveAt  time.Time
	Active    bool
	LenUser   int
	ScoreSum  int
	Status    int
	Fly       bool
	Expire    bool
	ReqStatus bool
	//Users        map[string]*Player
	CountUp      int
	CountDown    int
	RedDown      float32
	RedUp        float32
	RedCountDown int
	RedCountUp   int
	Timeout      int
	ScoreLimit   int
	Describe     string
	RedInterval  int
}
