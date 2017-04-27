package models

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	RoomList = make(map[string]*Room)
	roomLock sync.Mutex
)

type Room struct {
	id           string
	active       bool
	users        map[string]*User
	assistantSum int
	assistantNum int
	startTime    time.Time
	endTime      time.Time
	locker       sync.Mutex
	redhats      chan *redhat
	hasRedhat    bool
	score        chan int
	hasScore     bool
	echo         chan *result
}

type result struct {
	custom string
	score  int
}

type Mark struct {
	Custom string
	Score  int
}

func MakeReport(rs []*result) []Mark {
	marks := []Mark{}
	for _, r := range rs {
		marks = append(marks, Mark{
			Custom: r.custom,
			Score:  r.score,
		})
	}
	return marks
}

type redhat struct {
	master  string
	count   int
	timeout time.Duration
	end     bool
}

type RedReq struct {
	Number  int
	Timeout int
	Diver   int
	Master  string
}

type RoomReq struct {
	AssistantNum int
	Duration     int
	UserId       string
	Username     string
}

type RoomRespone struct {
	RoomId    string
	StartTime time.Time
	EndTime   time.Time
	Active    bool
	LenUser   int
}

func (r *Room) Convert() *RoomRespone {
	return &RoomRespone{
		RoomId:    r.id,
		StartTime: r.startTime,
		EndTime:   r.endTime,
		Active:    r.active,
		LenUser:   len(r.users),
	}
}

func CreateRoom(req *RoomReq) *Room {
	admin := NewUser(UserReq{
		UserId:   req.UserId,
		Username: req.Username,
	})
	roomid := string(Krand(8, 1))
	now := time.Now()
	room := &Room{
		id:           roomid,
		active:       true,
		users:        make(map[string]*User),
		assistantSum: req.AssistantNum,
		assistantNum: 0,
		startTime:    now,
		endTime:      now.Add(time.Duration(req.Duration) * time.Hour),
		redhats:      make(chan *redhat, 10),
		echo:         make(chan *result),
	}
	admin.Rooms[room] = &Player{
		Role:  Role_Admin,
		Score: 0,
	}
	room.users[admin.Id] = admin
	roomLock.Lock()
	RoomList[roomid] = room
	roomLock.Unlock()
	return room
}

func (r *Room) Active() bool {
	r.locker.Lock()
	r.active = r.endTime.After(time.Now())
	r.locker.Unlock()
	return r.active
}

func (r *Room) Close() {
	for _, u := range r.users {
		delete(u.Rooms, r)
		if len(u.Rooms) == 0 {
			userLock.Lock()
			delete(UserList, u.Id)
			userLock.Unlock()
		}
	}
	delete(RoomList, r.id)
}

func (r *Room) AppendUser(ur UserReq) error {
	if r.Active() {
		user, ok := r.users[ur.UserId]
		if !ok {
			user = NewUser(ur)
			r.locker.Lock()
			r.users[user.Id] = user
			r.locker.Unlock()
		}
		user.AppendRoom(r, Role_Custom)
		return nil
	}

	return fmt.Errorf("room is disable")
}

func (r *Room) Assistant(uid string) error {
	if r.Active() {
		user, ok := r.users[uid]
		if !ok {
			return fmt.Errorf("the user not in room")
		} else {
			user.Rooms[r].Role = Role_Assistant
		}
	}
	return nil
}

func (r *Room) SendRedhat(rr *RedReq) error {
	if !r.Active() {
		return fmt.Errorf("the room is disable")
	}

	if r.HaveRedhat() {
		return fmt.Errorf("have redhat!")
	}
	if rr.Number > 10 {
		return fmt.Errorf("number overflow")
	}

	if rr.Number == 1 {
		r.redhats <- &redhat{
			master:  rr.Master,
			count:   rr.Diver,
			timeout: time.Duration(rr.Timeout) * time.Minute,
			end:     true,
		}
		r.locker.Lock()
		r.hasRedhat = true
		r.locker.Unlock()
		return nil
	}
	for i := 0; i < rr.Number-1; i++ {
		r.redhats <- &redhat{
			master:  rr.Master,
			count:   rr.Diver,
			timeout: time.Duration(rr.Timeout) * time.Minute,
		}
	}

	r.redhats <- &redhat{
		master:  rr.Master,
		count:   rr.Diver,
		timeout: time.Duration(rr.Timeout) * time.Minute,
		end:     true,
	}
	r.locker.Lock()
	r.hasRedhat = true
	r.locker.Unlock()
	return nil

}

func (r *Room) Diver(master string) ([]*result, error) {

	if !r.Active() {
		return nil, fmt.Errorf("the room is disable")
	}

	if !r.HaveRedhat() {
		return nil, fmt.Errorf("have not redhat!")
	}

	rd := <-r.redhats
	if rd.master != master {
		return nil, fmt.Errorf("can't dive redhat!")
	}
	if r.HaveScore() {
		return nil, fmt.Errorf("releave score!")
	}

	r.score = make(chan int, rd.count+1)
	close(r.echo)
	r.echo = make(chan *result)

	GenerateScore(rd.count, r.score)
	r.locker.Lock()
	r.hasScore = true
	r.locker.Unlock()

	response := []*result{}
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
					if rd.end {
						r.locker.Lock()
						r.hasRedhat = false
						r.locker.Unlock()
					}
					rs.score = -rs.score
					response = append(response, rs)
					r.locker.Lock()
					r.hasScore = false
					r.locker.Unlock()
					return response, nil
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

func GenerateScore(count int, score chan int) {
	if count != 1 {
		for i := 0; i < count-1; i++ {
			score <- 1
		}
	}
	score <- -1

}

func Clear() {
	for _, room := range RoomList {
		if !room.active {
			if room.endTime.Before(room.startTime.Add(time.Since(room.endTime))) {
				room.Close()
			}
		}
	}
}
