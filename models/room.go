package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	RoomList       = make(map[int]*Room)
	roomLock       sync.Mutex
	RoomNamePrefix = "yoyo_"
)

const (
	Stat_Kongxian = iota
	Stat_Qizhuang
	Stat_Qiangzhuang
	Stat_Keepzhuang
	Stat_Conifgzhuang
	Stat_Sendredpaper
	Stat_Getredpaper
)

type Room struct {
	id        int
	gid       string
	name      string
	active    bool
	users     map[string]*Player
	water     int
	base      int
	describe  string
	duration  int
	superman  bool
	reqStatus bool
	CreateAt  time.Time
	ActiveAt  time.Time
	// assistantSum int
	// assistantNum int
	admin        string
	startTime    time.Time
	endTime      time.Time
	locker       sync.Mutex
	redhats      chan *redhat
	hasRedhat    bool
	redhatMaster string
	score        chan int
	hasScore     bool
	echo         chan *result
	results      map[string]*result
	status       int
	jushu        int
	//	gainlimit    int
	Scope
}

func GetRIDFromName(name string) int {
	for _, r := range RoomList {
		if r.name == name {
			return r.id
		}
	}
	return 0
}

func (r *Room) reset() {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.redhatMaster = ""
	r.hasRedhat = false
	r.hasScore = false
	r.results = map[string]*result{}
	r.status = 0
	r.superman = false
}

func (r *Room) Super(start bool) error {
	if r.Active() {
		r.superman = start
		return nil
	}
	return fmt.Errorf("room is disbaled!")
}

type result struct {
	custom string
	score  int
	bay    int
}

type Marks struct {
	RedId   string
	Master  string
	Water   int
	NN      string
	Jushu   int
	Results []Mark
}

type Mark struct {
	Custom   string
	Score    string
	Pay      int
	Nickname string
}

func (room *Room) makeReport(rs []*result, redid string) *Marks {
	marks := []Mark{}
	if len(rs) == 0 {
		return &Marks{
			RedId: redid,
		}
	}
	water := 0
	nn := ""
	for _, r := range rs {
		var nicname string
		player, ok := room.users[r.custom]
		if ok {
			nicname = player.NicName
		}
		nn = niuText(niu(r.score))
		marks = append(marks, Mark{
			Custom:   r.custom,
			Score:    f2str(float32(r.score) / 100),
			Pay:      r.bay,
			Nickname: nicname,
		})
		water += r.bay
	}

	return &Marks{
		RedId:   redid,
		Master:  marks[len(rs)-1].Custom,
		Water:   water,
		NN:      nn,
		Jushu:   room.jushu,
		Results: marks,
	}
}

type redhat struct {
	amount  int
	count   int
	timeout time.Duration
	end     bool
	base    int
}

type RedReq struct {
	Number int
	//Timeout   int
	//Diver     int
	//Master    string
	//RedAmount float32
	//	ScoreLimt int
}

type DiverReq struct {
	Diver     int
	RedAmount float32
}

type RoomReq struct {
	//AssistantNum int
	Duration  int
	UserId    string
	Nickname  string
	UserLimit int
	RoomName  string
	Water     int
	Base      int
	Scope
}

type UserMark struct {
	*Player
	Uid string
}

type RoomRespone struct {
	Id        int
	RoomId    string
	RoomName  string
	Base      int
	Water     int
	Admin     string
	Banker    string
	CreateAt  string
	StartTime string
	EndTime   string
	ActiveAt  string
	Active    bool
	LenUser   int
	ScoreSum  int
	Users     []UserMark
	Status    int
	Isused    bool
	Fly       bool
	Expire    bool
	ReqStatus bool
	LifeTime  int
	Scope
}

type RoomResponeNoUsers struct {
	Id        int
	RoomId    string
	RoomName  string
	Base      int
	Water     int
	Admin     string
	Banker    string
	CreateAt  string
	StartTime string
	EndTime   string
	ActiveAt  string
	Active    bool
	LenUser   int
	ScoreSum  int
	Status    int
	Isused    bool
	Fly       bool
	Expire    bool
	ReqStatus bool
	LifeTime  int
	Scope
}

type sortUsers []UserMark

func (su sortUsers) Less(i, j int) bool {
	return su[i].Score > su[j].Score
}

func (su sortUsers) Len() int {
	return len(su)
}

func (su sortUsers) Swap(i, j int) {
	su[i], su[j] = su[j], su[i]
}

func timeFormat(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}

func (r *Room) Convert() *RoomRespone {
	sumScore := 0
	users := []UserMark{}
	for k, u := range r.users {
		sumScore += u.Score
		users = append(users, UserMark{
			Player: r.users[k],
			Uid:    k,
		})
	}
	now := time.Now() //"2006-01-02 15:04:05
	expired := !r.endTime.IsZero() && now.Add(time.Duration(30)*time.Minute).After(r.endTime)
	sort.Sort(sortUsers(users))
	lifetime := int(r.endTime.Sub(r.startTime).Hours()) + r.duration
	return &RoomRespone{
		Id:        r.id,
		RoomId:    r.gid,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		Banker:    r.redhatMaster,
		StartTime: timeFormat(r.startTime),
		EndTime:   timeFormat(r.endTime),
		Active:    r.active,
		LenUser:   len(r.users),
		ScoreSum:  sumScore,
		Users:     users,
		Status:    r.status,
		Isused:    !r.startTime.IsZero(),
		Fly:       r.superman,
		Expire:    expired,
		ReqStatus: r.reqStatus,
		CreateAt:  timeFormat(r.CreateAt),
		ActiveAt:  timeFormat(r.ActiveAt),
		LifeTime:  lifetime,
		Scope: Scope{
			CountUp:      r.CountUp,
			CountDown:    r.CountDown,
			RedDown:      r.RedDown,
			RedUp:        r.RedUp,
			RedCountDown: r.RedCountDown,
			RedCountUp:   r.RedCountUp,
			Timeout:      r.Timeout,
			ScoreLimit:   r.ScoreLimit,
			Describe:     r.Describe,
			RedInterval:  r.RedInterval,
		},
	}
}

func (r *Room) ConvertNoUsers() *RoomResponeNoUsers {
	sumScore := 0
	for _, u := range r.users {
		sumScore += u.Score
	}

	now := time.Now()
	expired := !r.endTime.IsZero() && now.Add(time.Duration(30)*time.Minute).After(r.endTime)
	lifetime := int(r.endTime.Sub(r.startTime).Hours()) + r.duration
	return &RoomResponeNoUsers{
		Id:        r.id,
		RoomId:    r.gid,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		Banker:    r.redhatMaster,
		StartTime: timeFormat(r.startTime),
		EndTime:   timeFormat(r.endTime),
		Active:    r.active,
		LenUser:   len(r.users),
		ScoreSum:  sumScore,
		Status:    r.status,
		Isused:    !r.startTime.IsZero(),
		Fly:       r.superman,
		ReqStatus: r.reqStatus,
		Expire:    expired,
		CreateAt:  timeFormat(r.CreateAt),
		ActiveAt:  timeFormat(r.ActiveAt),
		LifeTime:  lifetime,
		Scope: Scope{
			CountUp:      r.CountUp,
			CountDown:    r.CountDown,
			RedDown:      r.RedDown,
			RedUp:        r.RedUp,
			RedCountDown: r.RedCountDown,
			RedCountUp:   r.RedCountUp,
			Timeout:      r.Timeout,
			ScoreLimit:   r.ScoreLimit,
			RedInterval:  r.RedInterval,
		},
	}
}

func (r *Room) SetStatus(stat int) {
	r.locker.Lock()
	r.status = stat
	r.locker.Unlock()
}

func CreateRoom(req *RoomReq) (*Room, error) {
	if len(RoomList) > 1000 {
		return nil, fmt.Errorf("the number of rooms overflow!")
	}
	for _, rm := range RoomList {
		if rm.reqStatus && req.UserId == rm.admin {
			return nil, fmt.Errorf("you have room %d in request", rm.id)
		}
	}
	gid, err := cemsdk.AddGroup(req.RoomName, "", req.UserId, true, false, req.UserLimit, nil)
	if err != nil {
		return nil, err
	}
	// admin := NewUser(UserReq{
	// 	UserId: req.UserId,
	// 	//Username: req.Username,
	// })
	//roomid := string(Krand(8, 1))
	//now := time.Now()

	var zero time.Time
	room := &Room{
		gid:       gid,
		active:    false,
		reqStatus: true,
		base:      req.Base,
		users:     make(map[string]*Player),
		admin:     req.UserId,
		duration:  req.Duration,
		CreateAt:  time.Now(),
		// assistantSum: req.AssistantNum,
		// assistantNum: 0,
		startTime: zero,
		endTime:   zero,
		redhats:   make(chan *redhat, 10),
		echo:      make(chan *result),
		results:   make(map[string]*result),
		water:     req.Water,
		Scope: Scope{
			CountDown:    req.CountDown,
			CountUp:      req.CountUp,
			RedCountDown: req.RedCountDown,
			RedCountUp:   req.RedCountUp,
			RedDown:      req.RedDown,
			RedUp:        req.RedUp,
			ScoreLimit:   req.ScoreLimit,
			Timeout:      req.Timeout,
			Describe:     req.Describe,
		},
	}

	adminPlayer := &Player{
		Role:    Role_Admin,
		Active:  true,
		Score:   0,
		NicName: req.Nickname,
		Head:    req.RoomName,
	}
	room.users[req.UserId] = adminPlayer

	pi := 0
	for {
		if pi > 7 {
			break
		}
		pi = rand.Intn(59)
	}
	rsid, err := idUp()
	if err != nil {
		return nil, err
	}
	room.id = pi*1000 + rsid
	room.name = fmt.Sprintf("%s的 %d号房间", req.Nickname, room.id)
	//fmt.Println(room.Update())
	_, err = room.Insert()
	if err != nil {
		cemsdk.DelGroup(gid)
		fmt.Println(err)
		return nil, err
	}
	roomLock.Lock()
	RoomList[room.id] = room
	//fmt.Println("len(roomlist)=", len(RoomList))
	roomLock.Unlock()
	return room, nil
}

type RoomConfig struct {
	Base  int
	Water int
	Scope
}

type Scope struct {
	CountUp      int     `db:"count_up"`
	CountDown    int     `db:"count_down"`
	RedDown      float32 `db:"red_down"`
	RedUp        float32 `db:"red_up"`
	RedCountDown int     `db:"redcount_down"`
	RedCountUp   int     `db:"redcount_up"`
	Timeout      int     `db:"timeout"`
	ScoreLimit   int     `db:"score_limit"`
	Describe     string  `db:"describe"`
	RedInterval  int     `db:"red_interval"`
}

func (r *Room) Config(req *RoomConfig) error {
	if r.status != Stat_Kongxian {
		return fmt.Errorf("you cant config room now!")
	}
	r.locker.Lock()
	defer r.locker.Unlock()
	r.base = req.Base
	r.water = req.Water
	r.CountDown = req.CountDown
	r.CountUp = req.CountUp
	r.RedCountDown = req.RedCountDown
	r.RedCountUp = req.RedCountUp
	r.RedDown = req.RedDown
	r.RedUp = req.RedUp
	r.Timeout = req.Timeout
	r.ScoreLimit = req.ScoreLimit
	r.Describe = req.Describe
	r.RedInterval = req.RedInterval
	return r.Update()
}

func (r *Room) Active() bool {
	if r.active == false {
		return false
	}
	r.locker.Lock()
	if r.duration != 0 || (!r.endTime.IsZero() && time.Now().Before(r.endTime)) {
		r.active = true
	} else {
		r.active = false
	}
	r.locker.Unlock()
	return r.active
}

func (r *Room) SetActive() {
	fmt.Println(r.active, r.reqStatus)
	r.active = true
	r.reqStatus = false
	r.ActiveAt = time.Now()
	r.Update()
}

func (r *Room) Cancle() {

	// r.reset()
	// r.end()
	r.active = false
	closeRoom <- struct{}{}
	emsay(r.gid, `{"type":"msg","msg":"房间已被管理员注销"}`)
}

func (r *Room) Close() {
	// for _, u := range r.users {
	// 	delete(u.Rooms, r)
	// 	if len(u.Rooms) == 0 {
	// 		userLock.Lock()
	// 		delete(UserList, u.Id)
	// 		userLock.Unlock()
	// 	}
	// }
	//delete(RoomList, r.id)
	r.reset()
	r.end()
	emsay(r.gid, `{"type":"msg","msg":"房间已过期"}`)
	//go cemsdk.DelGroup(r.gid)
}

func (r *Room) IsAdmin(uid string) bool {
	return r.admin == uid
}

func (r *Room) Role(uid string) int {
	u, ok := r.users[uid]
	if !ok {
		return -1
	}
	return u.Role
}

func (r *Room) IsCustom(uid string) bool {
	u, ok := r.users[uid]
	if !ok {
		return false
	}
	return u.Role == Role_Custom && u.Active
}

func (r *Room) IsAssistant(uid string) bool {
	u, ok := r.users[uid]
	if !ok {
		return false
	}
	return u.Role == Role_Assistant && u.Active
}

func (r *Room) IsFinace(uid string) bool {
	u, ok := r.users[uid]
	if !ok {
		return false
	}
	return u.Role == Role_Finace && u.Active
}

func (r *Room) IsAnyone(uid string) bool {
	u, ok := r.users[uid]
	return ok && u.Active
}

func (r *Room) AppendUser(openid string, nicname string, head string) (string, error) {
	if r.Active() {

		token, err := GetToken(openid)
		if err != nil {
			fmt.Println("huanxin:", err)
			return "", err
		}
		_, ok := r.users[openid]
		if ok {
			//r.users[openid].NicName = nicname
			r.users[openid].Head = head
			return token, nil
		}
		// if token == "" {
		// 	token, err = CreateDBUser(openid)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// }
		// err = cemsdk.AddUserToGroup(r.id, openid)
		// if err != nil && !strings.Contains(err.Error(), "already in group") {
		// 	return token, err
		// }
		r.users[openid] = &Player{
			Role:    Role_Custom,
			NicName: nicname,
			Head:    head,
		}
		emsay2user(r.admin, fmt.Sprintf(`{"type":"goin","uname":"%s","room":"%d"}`, nicname, r.id))
		// user, ok := r.users[ur.UserId]
		// if !ok {
		// 	user = NewUser(ur)
		// 	r.locker.Lock()
		// 	r.users[user.Id] = user
		// 	r.locker.Unlock()
		// }
		// user.AppendRoom(r, Role_Custom)
		return token, nil
	}

	return "", fmt.Errorf("room is disable")
}

func (r *Room) ModifyScore(openid string, score int) error {
	if r.Active() {
		player, ok := r.users[openid]
		if !ok {
			return fmt.Errorf("you must join room at first!")
		}
		player.Score = score
		err := player.Update(r.id, openid)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("room is disable")
}

func (r *Room) ActiveUser(openid string, req *UserActiveReq) error {
	if r.admin == openid {
		return fmt.Errorf("uid %s is room master", openid)
	}
	if r.Active() {
		player, ok := r.users[openid]
		if !ok {
			return fmt.Errorf("you must join room at first!")
		}
		if req.Active {
			var creater bool
			var updater bool
			var err error

			if player.Active != req.Active {
				errt := cemsdk.AddUserToGroup(r.gid, openid)
				defer func() {
					if err != nil {
						cemsdk.DelUserFromGroup(r.gid, openid)
					}
				}()
				if errt != nil && !strings.Contains(errt.Error(), "already in group") {
					err = errt
					return err
				}
				creater = true

			}
			player.Active = true

			if player.NicName != req.Nicname {
				err = cemsdk.ChangeNickname(openid, req.Nicname)
				if err != nil {
					return err
				}
				player.NicName = req.Nicname
				updater = true

			}
			if creater {
				err = player.Insert(r.id, openid)
				if err != nil {
					return err
				}
			} else if updater {
				err = player.Update(r.id, openid)
				if err != nil {
					return err
				}
			}

			count := 0
			for _, v := range r.users {
				if v.Active {
					count++
				}
			}
			// cemsdk.SendMessage("admin", "chatgroups", []string{r.id}, map[string]string{
			// 	"type": "txt",
			// 	"msg":  fmt.Sprintf("当前玩家数 %d", count),
			// }, map[string]string{})
			// cemsdk.SendMessage("admin", "chatgroups", []string{r.id}, map[string]string{
			// 	"type": "txt",
			// 	"msg":  fmt.Sprintf("玩家%s 加入房间", player.NicName),
			// }, map[string]string{})
			emsay(r.gid, fmt.Sprintf(`[{"type":"join"},{"user":"%s"},{"num":%d}]`, player.NicName, count))
			emsay2user(openid, fmt.Sprintf(`[{"type":"active"},{"active":1},{"room":%d}]`, r.id))

		} else {
			if player.Active != req.Active {
				err := cemsdk.DelUserFromGroup(r.gid, openid)
				if err != nil {
					return nil
				}
			}
			player.Active = false
			emsay2user(openid, `[{"type":"active"},{"active":1}`)
		}
		return nil

	}

	return fmt.Errorf("room is disable")

}

func (r *Room) Renew(dur int) {
	if dur == 0 {
		return
	}
	// r.locker.Lock()
	// defer r.locker.Unlock()
	// now := time.Now()
	// if r.endTime.Before(now) {
	// 	r.endTime = now.Add(time.Hour * time.Duration(dur))
	// } else {
	// 	r.endTime = r.endTime.Add(time.Hour * time.Duration(dur))
	// }
	// r.active = true
	if r.Active() && !r.endTime.IsZero() && time.Now().Before(r.endTime) {

		r.endTime = r.endTime.Add(time.Hour * time.Duration(dur))
		r.Update()

	} else {
		r.duration = r.duration + dur
		// r.active = true
		// r.reqStatus = false
		r.SetActive()
	}

	return
}

func (r *Room) Assistant(uid string, role int) error {
	if role == Role_Admin {
		return fmt.Errorf("cant Understand the role %d", role)
	}
	if r.Active() {
		// user, ok := r.users[uid]
		// if !ok {
		// 	return fmt.Errorf("the user not in room")
		// } else {
		// 	user.Rooms[r].Role = Role_Assistant
		// }
		_, ok := r.users[uid]
		if !ok || r.users[uid].Role == Role_Admin {
			return fmt.Errorf("%s is not in Room: %d", uid, r.id)
		}
		r.users[uid].Role = role
		emsay2user(uid, fmt.Sprintf(`[{"type":"role"},{"room":"%d"},{"role":%d}]`, r.id, role))
		return nil
	}
	return fmt.Errorf("room is disable")
}

func (r *Room) start() {
	r.locker.Lock()
	now := time.Now()
	r.startTime = now
	r.endTime = now.Add(time.Duration(r.duration) * time.Hour)
	r.duration = 0
	r.locker.Unlock()
	r.Update()
}

func (r *Room) end() {
	// if r.duration != 0 {
	// 	r.duration = 0
	// }
	// if !r.startTime.IsZero() {
	// 	r.endTime = time.Now()
	// }
	var zero time.Time
	r.locker.Lock()
	r.duration = 0
	r.startTime = zero
	r.endTime = zero
	r.active = false
	r.locker.Unlock()
	r.Update()
}

func (r *Room) SendRedhat() error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}
	if r.status != Stat_Kongxian {
		return fmt.Errorf("the room stat is %d!", r.status)
	}
	if r.duration != 0 {
		r.start()
	}
	r.SetStatus(Stat_Qizhuang)
	return nil
	/*
		if r.HaveRedhat() {
			return fmt.Errorf("have redhat!")
		}
		if rr.Number > 10 {
			return fmt.Errorf("number overflow")
		}
		if rr.Diver < 2 {
			return fmt.Errorf("Diver < 2")
		}

		// if rr.Number == 1 {
		// 	r.redhats <- &redhat{
		// 		count:   rr.Diver,
		// 		timeout: time.Duration(rr.Timeout) * time.Minute,
		// 		end:     true,
		// 	}
		// 	r.locker.Lock()
		// 	r.hasRedhat = true
		// 	r.locker.Unlock()
		// 	return nil
		// }
		for i := 0; i < rr.Number-1; i++ {
			r.redhats <- &redhat{
				count:   rr.Diver,
				timeout: time.Duration(rr.Timeout) * time.Minute,
				base:    rr.Base,
			}
		}

		r.redhats <- &redhat{
			count:   rr.Diver,
			timeout: time.Duration(rr.Timeout) * time.Minute,
			base:    rr.Base,
			end:     true,
		}
		r.locker.Lock()
		r.hasRedhat = true
		r.locker.Unlock()
		return nil
	*/

}

func (r *Room) MasterRedhat(master string) error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}
	// if !r.HaveRedhat() {
	// 	return fmt.Errorf("have not redhat!")
	// }
	if r.status != Stat_Qizhuang {
		return fmt.Errorf("the room stat is %d!", r.status)
	}
	r.locker.Lock()
	defer r.locker.Unlock()
	if r.redhatMaster != "" {
		return fmt.Errorf("master is someone else!")
	}
	r.redhatMaster = master
	r.status = Stat_Qiangzhuang
	var nicname string
	u, _ := cemsdk.GetUser(master)
	if u != nil {
		nicname = u.Nicname
	}
	mp, ok := r.users[master]
	if ok {
		nicname = mp.NicName
	}
	// cemsdk.SendMessage("admin", "chatgroups", []string{r.id}, map[string]string{
	// 	"type": "txt",
	// 	"msg":  fmt.Sprintf("%s[%s] 抢到庄家", nicname, master),
	// }, map[string]string{})
	emsay(r.gid, fmt.Sprintf(`{"type":"message","msg":"\"%s\" 抢到庄家.请点击连庄按钮配置连庄次数"}`, nicname))
	return nil
}

func (r *Room) KeepZhuang(master string) error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}
	if r.redhatMaster != master {
		return fmt.Errorf("you are not master!")
	}
	if r.status != Stat_Qiangzhuang {
		return fmt.Errorf("the room stat is %d!", r.status)
	}
	r.status = Stat_Keepzhuang
	return nil
}

func (r *Room) ConfigRedhat(rr *RedReq, cancel bool) error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}

	// if !r.HaveRedhat() {
	// 	return fmt.Errorf("have not redhat!")
	// }
	if r.status != Stat_Keepzhuang {
		return fmt.Errorf("the room stat is %d!", r.status)
	}
	if cancel {
		r.status = Stat_Qiangzhuang
		return nil
	}
	// if r.redhatMaster != master {
	// 	return fmt.Errorf("master is %s!", r.redhatMaster)
	// }
	if rr.Number > r.CountUp || rr.Number < r.CountDown || rr.Number < 1 {
		return fmt.Errorf("count overflow")
	}
	//r.gainlimit = rr.ScoreLimt

	if rr.Number == 1 {
		r.redhats <- &redhat{
			//amount:  int(rr.RedAmount * 100),
			count:   1,
			timeout: time.Duration(r.Timeout) * time.Minute,
			base:    r.base,
			end:     true,
		}
		r.locker.Lock()
		r.hasRedhat = true
		r.locker.Unlock()
		return nil
	}
	r.redhats = make(chan *redhat, rr.Number)
	for i := 0; i < rr.Number-1; i++ {
		r.redhats <- &redhat{
			//amount:  int(rr.RedAmount * 100),
			count:   rr.Number - i,
			timeout: time.Duration(r.Timeout) * time.Minute,
			base:    r.base,
		}
	}

	r.redhats <- &redhat{
		//amount:  int(rr.RedAmount * 100),
		count:   1,
		timeout: time.Duration(r.Timeout) * time.Minute,
		base:    r.base,
		end:     true,
	}
	// r.locker.Lock()
	// r.hasRedhat = true
	// r.locker.Unlock()
	r.SetStatus(Stat_Conifgzhuang)
	return nil

}

func (r *Room) Discard() error {
	if !(r.status > Stat_Kongxian) {
		return fmt.Errorf("the room stat is idle!")
	}
	if r.HaveScore() {
		return fmt.Errorf("red leaved!")
	}
	r.scoreClear()
	r.redhatClear()
	emsay(r.gid, "庄家弃庄")
	return nil
}

func (r *Room) Diver(master string, req *DiverReq) (*Marks, error) {

	if !r.Active() {
		return nil, fmt.Errorf("the room is disable")
	}

	// if !r.HaveRedhat() {
	// 	return nil, fmt.Errorf("have not redhat!")
	// }

	if r.status != Stat_Conifgzhuang {
		return nil, fmt.Errorf("the room stat is %d!", r.status)
	}

	if r.redhatMaster != master {
		return nil, fmt.Errorf("can't dive redhat!")
	}

	if r.HaveScore() {
		return nil, fmt.Errorf("releave score!")
	}

	if req.Diver < r.RedCountDown || req.Diver > r.RedCountUp {
		return nil, fmt.Errorf("Diver Count overflow")
	}

	if req.RedAmount < r.RedDown || req.RedAmount > r.RedUp {
		return nil, fmt.Errorf("RedAmount overflow")
	}

	rd := <-r.redhats
	leaveRed := rd.count
	rd.amount = int(req.RedAmount * 100)
	//rd.count = req.Diver
	r.score = make(chan int, req.Diver+1)
	if r.echo != nil {
		close(r.echo)
	}
	r.echo = make(chan *result)

	//传递信息给抢红的用户
	var mNic string
	pl, ok := r.users[master]
	if ok {
		mNic = pl.NicName
	}
	r.results = make(map[string]*result)
	r.results["000--"] = &result{
		custom: mNic,
		score:  rd.amount,
		bay:    req.Diver,
	}

	redid := string(Krand(8, 3))
	r.results["000-+"] = &result{
		custom: redid,
	}
	fmt.Printf("score amount=====: %d [%s]\n", rd.amount, r.id)
	GenerateScore(rd.amount, req.Diver, r.score, r.superman)
	var masterNic string
	mPlayer, ok := r.users[master]
	if ok {
		masterNic = mPlayer.NicName
	}
	// masterscore := <-r.score
	// response := []*result{&result{
	// 	custom: master,
	// 	score:  masterscore,
	// }}
	r.jushu += 1
	emsay(r.gid, fmt.Sprintf(`{"type":"redhat","redid":"%s","nicname":"%s","master":"%s","amount":%v,"diver":%d,"jushu":%d}`, redid, master, masterNic, req.RedAmount, req.Diver, r.jushu))
	response := []*result{}
	r.locker.Lock()
	r.hasScore = true
	r.locker.Unlock()
	nowcount := 0
	timecount := 0
	//exitcount=int(rd.timeout.Seconds())*2
	ticker := time.NewTicker(time.Second * 30)
	//ot:=time.After(rd.timeout)
	defer ticker.Stop()
	defer r.Update()
	for {
		select {
		// case <-ot:
		// 	if r.score != nil {
		// 		close(r.score)
		// 	}

		// 	// r.locker.Lock()
		// 	// r.hasScore = false
		// 	// r.locker.Unlock()
		// 	r.scoreClear()
		// 	r.redhats <- rd
		// 	//ticker.Stop()
		// 	fmt.Println("==============diver timeout...==============")
		// 	return nil, fmt.Errorf("diver timeout!")
		case <-ticker.C:
			timecount += 1
			leave := req.Diver - nowcount
			ltime := int(rd.timeout.Seconds()) - 30*timecount
			if ltime <= 0 {
				//without lock
				r.hasScore = false
				r.results = make(map[string]*result)
				//
				if r.score != nil {
					close(r.score)
				}

				r.locker.Lock()
				// r.locker.Lock()
				// r.hasScore = false
				// r.locker.Unlock()

				r.redhats <- rd
				if rd.end != true {
					for {
						rdl := <-r.redhats
						r.redhats <- rdl
						if rdl.end {
							break
						}
					}
				}
				r.locker.Unlock()
				//r.scoreClear()
				//ticker.Stop()
				emsay(r.gid, fmt.Sprintf(`{"type":"red","count":%d,"time": "红包超时"}`, leave))
				emsay(r.gid, fmt.Sprintf(`{"type":"msg","msg":"剩余坐庄次数%v"}`, leaveRed))
				//fmt.Println("==============diver timeout ticker.C...==============")
				return nil, fmt.Errorf("%s diver timeout!", redid)
			}
			emsay(r.gid, fmt.Sprintf(`{"type":"red","count":%d,"time": %d }`, leave, ltime))

		case rs := <-r.echo:
			nowcount += 1
			if rs != nil {
				if rs.score < 0 && len(response) >= (req.Diver-1) {
					if r.score != nil {
						close(r.score)
					}
					rs.score = -rs.score
					response = append(response, rs)
					r.scoreClear()
					reports := &Marks{}
					//if len(response) >= req.Diver {
					res := r.juge(response, rd.base, r.water, r.redhatMaster)
					reports = r.makeReport(res, redid)
					print(reports)
					Record(r.id, r.name, r.admin, reports, reports.Water, r.jushu)
					//}

					if rd.end {
						r.redhatClear()
						emsay(r.gid, `{"type":"message","msg":"本轮坐庄结束"}`)
					} else {
						emsay(r.gid, fmt.Sprintf(`{"type":"msg","msg":"剩余坐庄次数%v"}`, leaveRed-1))
					}

					// nicname:=""
					// um,ok1:=r.users[master]
					// if ok1{
					// 	nicname=um.NicName
					// }
					return reports, nil
				} else if rs.score < 0 && len(response) < (req.Diver-1) {
					if r.score != nil {
						close(r.score)
					}
					r.scoreClear()
					return nil, fmt.Errorf("leave red!")
				}
				response = append(response, rs)

			} else {
				continue
			}

		}
	}

}

func emsay(rid string, msg string) {
	fmt.Printf("========huanxin message=====[%s]\n", msg)
	cemsdk.SendMessage("admin", "chatgroups", []string{rid}, map[string]string{
		"type": "txt",
		"msg":  msg,
	}, map[string]string{})
}

func emsay2user(uid string, msg string) {
	cemsdk.SendMessage("admin", "users", []string{uid}, map[string]string{
		"type": "txt",
		"msg":  msg,
	}, map[string]string{})
}

type ScoreUnion struct {
	Master string
	Score  string
	Amount string
	Count  int
	RedId  string
}

func (r *Room) GetScore(custom string) (*ScoreUnion, error) {
	if !r.Active() {
		return nil, fmt.Errorf("the room is disable")
	}

	if !r.HaveScore() {
		return nil, fmt.Errorf("have no score")
	}

	//	score, isClosed := <-r.score
	//	if isClosed {
	//		r.locker.Lock()
	//		r.hasScore = false
	//		r.locker.Unlock()
	//		return 0, fmt.Errorf("diver timeout")
	//	}
	if _, ok := r.results[custom]; ok {
		return nil, fmt.Errorf("you have a score")
	}
	if custom == r.redhatMaster {
		return nil, fmt.Errorf("you are master of the work")
	}
	r.results[custom] = nil
	if r.users[custom].Score < r.ScoreLimit {
		return nil, fmt.Errorf("your score is below the limit")
	}
	redid := ""
	redresult, ok := r.results["000-+"]
	if ok {
		redid = redresult.custom
	}
	var nicname string
	pl, ok := r.users[custom]
	if ok {
		nicname = pl.NicName
	}
	r.locker.Lock()
	defer r.locker.Unlock()
	score, isOk := <-r.score
	if !isOk {
		return nil, fmt.Errorf("have no score")
	}

	if score < 0 {

		r.echo <- &result{
			custom: custom,
			score:  score,
		}
		r.hasScore = false
		//r.locker.Unlock()
		res := &ScoreUnion{}
		res.Score = f2str(float32(-score) / 100)
		//获取红信息
		if rt, ok := r.results["000--"]; ok {
			res.Master = rt.custom
			res.Count = rt.bay
			res.Amount = f2str(float32(rt.score) / 100)
			res.RedId = redid
		}

		RedInsert(r.id, redid, custom, nicname, res.Score)
		return res, nil
		//return -score, nil
	}
	r.echo <- &result{
		custom: custom,
		score:  score,
	}
	//r.locker.Unlock()
	res := &ScoreUnion{}
	res.Score = f2str(float32(score) / 100)
	//获取红信息
	if rt, ok := r.results["000--"]; ok {
		res.Master = rt.custom
		res.Count = rt.bay
		res.Amount = f2str(float32(rt.score) / 100)
		res.RedId = redid
	}
	RedInsert(r.id, redid, custom, nicname, res.Score)
	return res, nil

}

func (r *Room) NeedClean() bool {
	return r.active == false && (r.duration != 0 || !r.endTime.IsZero())
}
func (r *Room) redhatClear() {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.redhatMaster = ""
	r.hasRedhat = false
	//	r.gainlimit = -3000
	r.status = Stat_Kongxian
}

func (r *Room) scoreClear() {
	r.locker.Lock()
	defer r.locker.Unlock()
	r.hasScore = false
	r.results = make(map[string]*result)
}

func (r *Room) HaveScore() bool {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.hasScore

}

func (r *Room) HaveRedhat() bool {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.hasRedhat

}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func maxHead(is []int) {
	max := 0
	index := 0
	for i := range is {
		nc := niu(is[i])
		if nc > max {
			max = nc
			index = i
		}
	}
	is[len(is)-1], is[index] = is[index], is[len(is)-1]
}
func GenerateScore(sum int, count int, score chan int, is bool) {
	s1 := rand.New(rand.NewSource(time.Now().Unix() + rand.Int63n(100)))
	//s2 := rand.New(rand.NewSource(time.Now().Unix() - 100))
	suarray := make([]int, count)
	if count != 1 {
		for i := 0; i < count-1; i++ {
			// c := rand.Intn(4) + 1
			// var rs int
			// for j := 0; j < c; j++ {
			rs := s1.Intn((sum+i-count-1)/2) + 1
			// }
			//rs = rs

			sum = sum - rs
			suarray[i] = rs
			//score <- rs
		}
	}
	suarray[count-1] = sum
	if is {
		maxHead(suarray)
	}
	for i := 0; i < count; i++ {
		fmt.Printf("generate score: %d\n", suarray[i])
		if i == count-1 {
			score <- -suarray[i]
		} else {
			score <- suarray[i]
		}

	}
	//score <- -sum

}

func (r *Room) juge(rs []*result, base int, water int, master string) []*result {
	length := len(rs)
	if length < 1 {
		return nil
	}
	last := rs[length-1]
	//begin.score = -begin.score
	masterscore := niu(last.score)
	sum := 0
	for i := range rs {

		if i == length-1 {
			break
		}

		if rs[i].score < 0 {
			rs[i].score = -rs[i].score
		}
		ba := niu(rs[i].score)
		if ba < masterscore || (masterscore <= 2 && ba == masterscore) {
			rs[i].bay = -masterscore * base
			sum += masterscore * base
		}
		if ba > masterscore {
			rs[i].bay = ba * base
			sum -= ba * base
		}
		r.users[rs[i].custom].Score += rs[i].bay
		r.users[rs[i].custom].Update(r.id, rs[i].custom)
		//rs[i].score = float32(rs[i].score) / 100
	}
	last.bay = 0
	r.users[master].Score += sum - water*len(rs)
	r.users[master].Update(r.id, master)

	masterPlayer := &result{
		custom: master,
		score:  last.score,
		bay:    sum - water*len(rs),
	}
	rs = append(rs, masterPlayer)
	return rs
	//last.score = float32(last.score) / 100
}

func niu(score int) int {
	a := score / 100
	c := score % 10
	b := (score - a*100 - c) / 10
	//fmt.Println(a, b, c)
	if a+1 == b && b+1 == c {
		return 16
	}
	if a == b && b == c {
		return 14
	}
	if b == c && c == 0 {
		return 13
	}
	if b == c {
		return 12
	}
	if a != 0 && b == 0 && c == 1 {
		return 11
	}

	if b == 1 && c == 0 {
		return 11
	}
	if b+c == 10 {
		return 10
	}

	return (b + c) % 10
}

func niuText(i int) string {
	switch i {
	case 10:
		return "牛牛"
	case 11:
		return "超牛"
	case 12:
		return "对子"
	case 13:
		return "满牛"
	case 14:
		return "豹子"
	case 16:
		return "顺子"
	default:
		return fmt.Sprintf("牛%d", i)
	}
}

func abs(x int64) int64 {
	if x >= 0 {
		return x
	}
	return -x
}

func Clear() {
	for _, room := range RoomList {
		if room.reqStatus {
			return
		}
		if !room.Active() {
			if room.NeedClean() {
				room.Close()
			}
			//if room.endTime.Before(room.startTime.Add(time.Since(room.endTime))) {
			//room.DeleteDB()
			//}
		} else {
			if room.endTime.IsZero() {
				continue
			}
			now := time.Now().Unix()
			end := room.endTime.Unix()
			late := end - now

			if abs(late-1800) <= 2 { //30分钟
				emsay(room.gid, `{"late":30}`)
			} else if abs(late-1200) <= 2 { //20分钟
				emsay(room.gid, `{"late":20}`)
			} else if abs(late-600) <= 2 { //10分钟
				emsay(room.gid, `{"late":10}`)
			} else if abs(late-300) <= 2 { //5分钟
				emsay(room.gid, `{"late":5}`)
			} else if abs(late-60) <= 2 { //1分钟
				emsay(room.gid, `{"late":1}`)
			}

		}
	}
}

var closeRoom = make(chan struct{}, 5)

func init() {

	go func() {
		tricker := time.NewTicker(time.Second * 3)

		for {
			select {
			case <-tricker.C:
				Clear()
			case <-closeRoom:
				Clear()
			}
		}
	}()
	for i := 0; i < 3; i++ {
		if client, err := newEm(); err != nil {
			if i == 2 {
				fmt.Println(err)
				panic("emsdk init fail")
			}
		} else {
			cemsdk = client
			break
		}
	}
	dBEngine = DBEngineInit()
	if err := dBEngine.CreateTablesIfNotExists(); err != nil {
		panic("db init fail")
	}
	adminInsert()
	RoomInit(dBEngine)
	FilesInit()
	fmt.Println("Init finished...")
	// list, err := cemsdk.FetchAllGroupFromApp()
	// if err != nil {
	// 	panic(err)
	// }
	// for _, v := range list.Data {
	// 	if strings.HasPrefix(v.Groupname, RoomNamePrefix) {
	// 		cemsdk.DelGroup(v.Groupid)
	// 	}
	// }

}

func print(js interface{}) {
	bt, _ := json.MarshalIndent(js, "", " ")
	fmt.Println("print:", string(bt))
}

func f2str(f float32) string {
	return fmt.Sprintf("%.2f", f)
}
