package models

import "time"

type TmpRoomReq struct {
	Duration     int
	UserId       string
	UserLimit    int
	RoomName     string
	Water        int
	Base         int
	ScoreLimit   int
	CountUp      int
	CountDown    int
	RedDown      float32
	RedUp        float32
	RedCountDown int
	RedCountUp   int
	Timeout      int
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
}

type TmpRespone struct {
	RoomId    string
	RoomName  string
	Base      int
	Water     int
	Admin     string
	Banker    string
	StartTime time.Time
	EndTime   time.Time
	Active    bool
	LenUser   int
	ScoreSum  int
	Status    int
	//Users        map[string]*Player
	CountUp      int
	CountDown    int
	RedDown      float32
	RedUp        float32
	RedCountDown int
	RedCountUp   int
	Timeout      int
	ScoreLimit   int
}
