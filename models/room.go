package models

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	RoomList       = make(map[string]*Room)
	roomLock       sync.Mutex
	RoomNamePrefix = "yoyo_"
)

const (
	Stat_Kongxian = iota
	Stat_Qizhuang
	Stat_Qiangzhuang
	Stat_Conifgzhuang
	Stat_Sendredpaper
	Stat_Getredpaper
)

type Room struct {
	id     string
	name   string
	active bool
	users  map[string]*Player
	water  int
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
	gainlimit    int
}

type result struct {
	custom string
	score  int
	bay    int
}

type Marks struct {
	Master  string
	Water   int
	Results []Mark
}

type Mark struct {
	Custom string
	Score  float32
	Pay    int
}

func MakeReport(rs []*result) *Marks {
	marks := []Mark{}
	water := 0
	for _, r := range rs {
		marks = append(marks, Mark{
			Custom: r.custom,
			Score:  float32(r.score) / 100,
			Pay:    r.bay,
		})
		water += r.bay
	}

	return &Marks{
		Master:  marks[0].Custom,
		Water:   water,
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
	Number    int
	Timeout   int
	Diver     int
	Master    string
	Base      int
	RedAmount float32
	//	ScoreLimt int
}

type DiverReq struct {
	Diver     int
	RedAmount float32
}

type RoomReq struct {
	//AssistantNum int
	Duration   int
	UserId     string
	UserLimit  int
	RoomName   string
	Water      int
	ScoreLimit int
}

type RoomRespone struct {
	RoomId    string
	RoomName  string
	Admin     string
	Banker    string
	StartTime time.Time
	EndTime   time.Time
	Active    bool
	LenUser   int
	ScoreSum  int
	Users     map[string]*Player
}

func (r *Room) Convert() *RoomRespone {
	sumScore:=0
	for _,u:=range r.users{
		sumScore+=u.Score
	}
	return &RoomRespone{
		RoomId:    r.id,
		RoomName:  r.name,
		Admin:     r.admin,
		Banker:    r.redhatMaster,
		StartTime: r.startTime,
		EndTime:   r.endTime,
		Active:    r.active,
		LenUser:   len(r.users),
		ScoreSum: sumScore,
		Users:     r.users,
	}
}

func (r *Room) SetStatus(stat int) {
	r.locker.Lock()
	r.status = stat
	r.locker.Unlock()
}

func CreateRoom(req *RoomReq) (*Room, error) {
	gid, err := cemsdk.AddGroup(req.RoomName, "", req.UserId, true, false, req.UserLimit, nil)
	if err != nil {
		return nil, err
	}
	// admin := NewUser(UserReq{
	// 	UserId: req.UserId,
	// 	//Username: req.Username,
	// })
	//roomid := string(Krand(8, 1))
	now := time.Now()
	room := &Room{
		id:     gid,
		active: true,
		users:  make(map[string]*Player),
		admin:  req.UserId,
		// assistantSum: req.AssistantNum,
		// assistantNum: 0,
		startTime: now,
		endTime:   now.Add(time.Duration(req.Duration) * time.Hour),
		redhats:   make(chan *redhat, 10),
		echo:      make(chan *result),
		results:   make(map[string]*result),
		water:     req.Water,
		name:      req.RoomName,
		gainlimit: req.ScoreLimit,
	}

	adminPlayer := &Player{
		Role:   Role_Admin,
		Active: true,
		Score:  0,
	}
	room.users[req.UserId] = adminPlayer
	roomLock.Lock()
	RoomList[gid] = room
	roomLock.Unlock()
	return room, nil
}

func (r *Room) Active() bool {
	r.locker.Lock()
	r.active = r.endTime.After(time.Now())
	r.locker.Unlock()
	return r.active
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
	delete(RoomList, r.id)
	go cemsdk.DelGroup(r.id)
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

func (r *Room) IsAnyone(uid string) bool {
	u, ok := r.users[uid]
	return ok && u.Active
}

func (r *Room) AppendUser(openid string, nicname string) (string, error) {
	if r.Active() {

		token, err := GetToken(openid)
		if err != nil {
			return "", err
		}
		_, ok := r.users[openid]
		if ok {
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
		}

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

func (r *Room) ActiveUser(openid string, req *UserActiveReq) error {

	if r.Active() {
		player, ok := r.users[openid]
		if !ok {
			return fmt.Errorf("you must join room at first!")
		}
		if req.Active {
			if player.Active != req.Active {
				err := cemsdk.AddUserToGroup(r.id, openid)
				if err != nil && !strings.Contains(err.Error(), "already in group") {
					return err
				}
			}
			player.Active = true
			player.Score = req.Score
			if player.NicName != req.Nicname {
				err := cemsdk.ChangeNickname(openid, req.Nicname)
				if err != nil {
					return err
				}
				u:=&DBUser{}
				u.Id=openid
				err=u.Fetch(dBEngine)
				if err != nil {
					return err
				}
				u.NickName=req.Nicname
				err=u.Update(dBEngine)
				if err != nil {
					return err
				}
			}
			player.NicName = req.Nicname

		} else {
			if player.Active != req.Active {
				err := cemsdk.DelUserFromGroup(r.id, openid)
				if err != nil {
					return nil
				}
			}
			player.Active = false

		}
		return nil

	}

	return fmt.Errorf("room is disable")

}

func (r *Room) Assistant(uid string,role int) error {
	if role>1{
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
		if !ok {
			return fmt.Errorf("%s is not in Room: %s", uid, r.id)
		}
		r.users[uid].Role = role
		return nil
	}
	return fmt.Errorf("room is disable")
}

func (r *Room) SendRedhat() error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}
	if r.status != Stat_Kongxian {
		return fmt.Errorf("the room stat is %d!", r.status)
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
	cemsdk.SendMessage("admin", "chatgroups", []string{r.id}, map[string]string{
		"type": "txt",
		"msg":  fmt.Sprintf("%s[%s] 抢到庄家", nicname, master),
	}, map[string]string{})
	return nil
}

func (r *Room) ConfigRedhat(master string, rr *RedReq) error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}

	// if !r.HaveRedhat() {
	// 	return fmt.Errorf("have not redhat!")
	// }
	if r.status != Stat_Qiangzhuang {
		return fmt.Errorf("the room stat is %d!", r.status)
	}
	if r.redhatMaster != master {
		return fmt.Errorf("master is %s!", r.redhatMaster)
	}
	if rr.Number > 10 || rr.Number < 2 {
		return fmt.Errorf("number overflow")
	}
	//r.gainlimit = rr.ScoreLimt

	if rr.Number == 1 {
		r.redhats <- &redhat{
			amount:  int(rr.RedAmount * 100),
			count:   rr.Diver,
			timeout: time.Duration(rr.Timeout) * time.Minute,
			base:    rr.Base,
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
			amount:  int(rr.RedAmount * 100),
			count:   rr.Diver,
			timeout: time.Duration(rr.Timeout) * time.Minute,
			base:    rr.Base,
		}
	}

	r.redhats <- &redhat{
		amount:  int(rr.RedAmount * 100),
		count:   rr.Diver,
		timeout: time.Duration(rr.Timeout) * time.Minute,
		base:    rr.Base,
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
	r.scoreClear()
	r.redhatClear()
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

	if req.Diver < 2 {
		return nil, fmt.Errorf("Diver < 2")
	}

	if req.RedAmount < 1 || req.RedAmount > 9.99 {
		return nil, fmt.Errorf("RedAmount overflow")
	}

	rd := <-r.redhats
	rd.amount = int(req.RedAmount * 100)
	rd.count = req.Diver
	r.score = make(chan int, rd.count+1)
	close(r.echo)
	r.echo = make(chan *result)

	GenerateScore(rd.amount, rd.count, r.score)
	masterscore := <-r.score
	response := []*result{&result{
		custom: master,
		score:  masterscore,
	}}
	r.locker.Lock()
	r.hasScore = true
	r.locker.Unlock()

	for {
		select {
		case <-time.After(rd.timeout):
			close(r.score)
			r.locker.Lock()
			r.hasScore = false
			r.locker.Unlock()
			r.redhats <- rd
			return nil, fmt.Errorf("diver timeout!")
		case rs := <-r.echo:
			if rs != nil {
				if rs.score < 0 {

					rs.score = -rs.score
					response = append(response, rs)
					r.scoreClear()
					r.juge(response, rd.base, r.water)
					reports := MakeReport(response)
					Record(r.id, r.name, r.admin, reports)
					if rd.end {
						r.redhatClear()
					}
					return reports, nil
				}
				response = append(response, rs)

			}

		}
	}

}

func (r *Room) GetScore(custom string) (int, error) {
	if !r.Active() {
		return 0, fmt.Errorf("the room is disable")
	}

	if !r.HaveScore() {
		return 0, fmt.Errorf("have no score")
	}

	//	score, isClosed := <-r.score
	//	if isClosed {
	//		r.locker.Lock()
	//		r.hasScore = false
	//		r.locker.Unlock()
	//		return 0, fmt.Errorf("diver timeout")
	//	}
	if _, ok := r.results[custom]; ok {
		return 0, fmt.Errorf("you have a score")
	}
	if custom == r.redhatMaster {
		return 0, fmt.Errorf("you are master of the work")
	}
	r.results[custom] = nil
	if r.users[custom].Score < r.gainlimit {
		return 0, fmt.Errorf("your score is below the limit")
	}
	score := <-r.score

	if score < 0 {
		r.locker.Lock()
		r.echo <- &result{
			custom: custom,
			score:  score,
		}
		r.hasScore = false
		r.locker.Unlock()
		return -score, nil
	}

	r.echo <- &result{
		custom: custom,
		score:  score,
	}
	return score, nil

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

func GenerateScore(sum int, count int, score chan int) {
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
	for i := 0; i < count; i++ {
		if i == count-1 {
			score <- -suarray[i]
		} else {
			score <- suarray[i]
		}

	}
	//score <- -sum

}

func (r *Room) juge(rs []*result, base int, water int) {
	length := len(rs)
	if length < 2 {
		return
	}
	begin := rs[0]
	//begin.score = -begin.score
	master := niu(begin.score)
	sum := 0
	for i := range rs {

		if i == 0 {
			continue
		}

		if rs[i].score < 0 {
			rs[i].score = -rs[i].score
		}
		ba := niu(rs[i].score)
		if ba < master {
			rs[i].bay = -master * base
			sum += master * base

		}
		if ba > master {
			rs[i].bay = ba * base
			sum -= ba * base
		}
		r.users[rs[i].custom].Score += rs[i].bay
		//rs[i].score = float32(rs[i].score) / 100
	}
	begin.bay = sum - water*len(rs)
	r.users[begin.custom].Score += begin.bay
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
	if b == 0 && c == 1 {
		return 11
	}
	if b+c == 10 {
		return 10
	}
	return (b + c) % 10
}

func Clear() {
	for _, room := range RoomList {
		if !room.Active() {
			if room.endTime.Before(room.startTime.Add(time.Since(room.endTime))) {
				room.Close()
			}
		}
	}
}

func init() {

	go func() {
		tricker := time.NewTicker(time.Second * 30)

		for {
			<-tricker.C
			Clear()
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
	list, err := cemsdk.FetchAllGroupFromApp()
	if err != nil {
		panic(err)
	}
	for _, v := range list.Data {
		if strings.HasPrefix(v.Groupname, RoomNamePrefix) {
			cemsdk.DelGroup(v.Groupid)
		}
	}

}
