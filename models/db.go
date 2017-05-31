package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"strings"

	"sort"

	"encoding/base64"

	"gopkg.in/gorp.v1"
)

type DBUser struct {
	Id        string    `db:"id","primarykey" json:"id"`
	NickName  string    `db:"nick_name" json:"nick_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"update_at" json:"update_at"`
	Password  string    `db:"password" json:"password"`
}

func (u DBUser) TableName() string {
	return "user"
}

func (u *DBUser) Insert(db gorp.SqlExecutor) error {
	if u != nil {
		now := time.Now()
		u.CreatedAt = now
		u.UpdatedAt = now
		return db.Insert(u)
	}
	return nil
}

func (u *DBUser) Update(db gorp.SqlExecutor) error {
	if u != nil {
		now := time.Now()
		u.UpdatedAt = now
		_, err := db.Update(u)
		return err
	}
	return nil
}

func (u *DBUser) Fetch(db gorp.SqlExecutor) error {

	query := fmt.Sprintf(`select * from %s where id = "%s"`, u.TableName(), u.Id)
	return db.SelectOne(u, query)

}

func CreateDBUser(openid, nicname string) (string, error) {
	u := &DBUser{
		Id: openid,
	}

	if err := u.Fetch(dBEngine); err == nil {
		if u.NickName == "" {
			if err := cemsdk.ChangeNickname(openid, nicname); err != nil {
				return "", err
			}
			return nicname, nil
		}
		return u.NickName, nil
	}

	pas := string(Krand(8, KC_RAND_KIND_LOWER))
	//fmt.Println("insert uid:", openid, "pas:", pas)
retry:
	err := cemsdk.CreateAccount(openid, pas, nicname)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("%s exists", openid)) {
			if err := cemsdk.DeleteAccount(openid); err != nil {
				return "", err
			}
			goto retry
		}
		return "", err
	}
	defer func() {
		if err != nil {
			cemsdk.DeleteAccount(openid)
		}
	}()
	u.Password = pas
	u.NickName = ""
	//fmt.Println("db insert uid:", u.Id, "pas:", u.Password)
	var trans *gorp.Transaction
	trans, err = dBEngine.Begin()
	defer func() {
		if err == nil {
			trans.Commit()
		} else {
			trans.Rollback()
		}
	}()
	err = u.Insert(trans)
	if err != nil {
		return "", err
	}
	// token := ""
	// token, err = cemsdk.GetUserToken(u.Id, u.Password)
	// return token, err
	return nicname, nil
}

func GetToken(openid string) (string, error) {
	u := &DBUser{
		Id: openid,
	}
	// trans, err := dBEngine.Begin()
	// defer CheckAndCommit(trans, err)
	err := u.Fetch(dBEngine)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	//fmt.Println("uid:", u.Id, "pass:", u.Password)
	token := ""
	token, err = cemsdk.GetUserToken(u.Id, u.Password)
	return token, err

}

type DBRecord struct {
	Id        int64     `db:"id","primarykey","autoincrement" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	RoomId    int       `db:"room_id" json:"room_id"`
	RoomName  string    `db:"room_name" json:"room_name"`
	Master    string    `db:"master" json:"master"`
	Water     int       `db:"water" json:"water"`
	Jusu      int       `db:"jushu" json:"jushu"`
	Body      string    `db:"body" json:"body"`
}

func (r *DBRecord) Insert(db gorp.SqlExecutor) error {
	if r != nil {
		r.CreatedAt = time.Now()
		return db.Insert(r)
	}
	return nil
}

func (r DBRecord) TableName() string {
	return "record"
}

type sortRecord []*DBRecord

func (s sortRecord) Less(i, j int) bool {
	return s[i].CreatedAt.After(s[j].CreatedAt)
}

func (s sortRecord) Len() int {
	return len(s)
}

func (s sortRecord) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (r *DBRecord) List(db gorp.SqlExecutor, page int, size int) (interface{}, error) {
	list := []*DBRecord{}
	query := fmt.Sprintf(`select * from %s where room_id = "%d"`, r.TableName(), r.RoomId)

	_, err := db.Select(&list, query)
	if err != nil {
		return nil, err
	}
	sum := 0
	for _, rec := range list {
		sum += rec.Water
	}
	sort.Sort(sortRecord(list))
	resp := &struct {
		Pagination
		Sum  int
		Data []*DBRecord `json:"data"`
	}{}
	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	start, end := PageLocate(len(list), size, page)
	resp.Total = len(list)
	resp.TotalPage = resp.Total / size
	if resp.Total%size != 0 {
		resp.TotalPage += 1
	}
	resp.Data = list[start:end]
	resp.Sum = -sum
	for i := range resp.Data {
		resp.Data[i].Body = decodeNic(resp.Data[i].Body)
	}
	return resp, nil
}

func (r *DBRecord) ListDBRoom(page int, size int) (interface{}, error) {
	list := []*DBRecord{}
	query := fmt.Sprintf(`select distinct room_id, room_name,master from %s `, r.TableName())

	_, err := dBEngine.Select(&list, query)
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Pagination
		Data []*DBRecord `json:"data"`
	}{}
	if size == 0 {
		size = 20
	}
	if page == 0 {
		page = 1
	}
	start, end := PageLocate(len(list), size, page)
	resp.Total = len(list)
	resp.TotalPage = resp.Total / size
	if resp.Total%size != 0 {
		resp.TotalPage += 1
	}
	resp.Data = list[start:end]
	return resp, nil
	//return list, err
}

type DBRoom struct {
	Id        int       `db:"id","primarykey" json:"id"`
	RoomId    string    `db:"room_id","primarykey"`
	RoomName  string    `db:"room_name"`
	CreateAt  time.Time `db:"create_at"`
	Base      int       `db:"base"`
	Water     int       `db:"water"`
	Admin     string    `db:"admin"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	Duration  int       `db:"duration"`
	ReqStatus bool      `db:"req_status"`
	ActiveAt  time.Time `db:"active_at"`
	Jushu     int       `db:"jushu"`
	//	Describe  string    `db:"discrible"`
	Scope
}

func (room *DBRoom) TableName() string {
	return "db_room"
}

func (room *DBRoom) Insert(db gorp.SqlExecutor) error {
	room.CreateAt = time.Now()
	fmt.Println("create time:", room.CreateAt)
	return db.Insert(room)
}

func (room *DBRoom) Fetch(db gorp.SqlExecutor) error {
	return db.SelectOne(room, fmt.Sprintf(`select * from %s where id="%v" `, room.TableName(), room.Id))
}

func (room *DBRoom) Update(db gorp.SqlExecutor) error {
	fmt.Println("update time:", room.CreateAt)
	_, err := db.Update(room)
	return err
}

type RoomUser struct {
	Id     int64  `db:"id","primarykey","autoincrement" json:"id"`
	Uid    string `db:"uid"`
	RoomId int    `db:"room_id"`
	Player
}

func (ruser *RoomUser) TableName() string {
	return "room_user"
}

func (ruser *RoomUser) Insert(db gorp.SqlExecutor) error {
	return db.Insert(ruser)
}

func (ruser *RoomUser) Fetch(db gorp.SqlExecutor) error {
	return db.SelectOne(ruser, fmt.Sprintf(`select * from %s where uid="%s" and room_id="%d"`, ruser.TableName(), ruser.Uid, ruser.RoomId))
}

func (ruser *RoomUser) Update(db gorp.SqlExecutor) error {
	_, err := db.Update(ruser)
	return err
}

func (r *Room) Insert() (int, error) {
	trans, err := dBEngine.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			trans.Rollback()
		} else {
			trans.Commit()
		}
	}()
	dbRoom := &DBRoom{
		Id:        r.id,
		RoomId:    r.gid,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		StartTime: r.startTime,
		EndTime:   r.endTime,
		Duration:  r.duration,
		ReqStatus: r.reqStatus,
		CreateAt:  r.CreateAt,
		Jushu:     r.jushu,
		//		Describe:  r.describe,
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
	err = dbRoom.Insert(trans)
	if err != nil {
		return 0, err
	}
	for uid, player := range r.users {
		ruser := &RoomUser{
			Uid:    uid,
			RoomId: dbRoom.Id,
			Player: *player,
		}
		ruser.NicName = encodeNic(ruser.NicName)
		err = ruser.Insert(trans)
		if err != nil {
			return 0, err
		}
	}
	return dbRoom.Id, err
}

func (r *Room) Update() error {
	trans, err := dBEngine.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trans.Rollback()
		} else {
			trans.Commit()
		}
	}()
	dbRoom := &DBRoom{
		Id:        r.id,
		RoomId:    r.gid,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		StartTime: r.startTime,
		EndTime:   r.endTime,
		Duration:  r.duration,
		ReqStatus: r.reqStatus,
		CreateAt:  r.CreateAt,
		ActiveAt:  r.ActiveAt,
		Jushu:     r.jushu,
		//		Describe:  r.describe,
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

	err = dbRoom.Update(trans)
	return err
}

func (r *Room) Fetch() error {
	dbRoom := &DBRoom{}
	dbRoom.Id = r.id
	trans, err := dBEngine.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trans.Rollback()
		} else {
			trans.Commit()
		}
	}()
	err = dbRoom.Fetch(trans)
	if err != nil {
		return err
	}
	rlist := []*RoomUser{}
	_, err = trans.Select(&rlist, fmt.Sprintf(`select * from %s where  room_id="%d"`, (&RoomUser{}).TableName(), r.id))
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	r.jushu = dbRoom.Jushu
	r.gid = dbRoom.RoomId
	r.name = dbRoom.RoomName
	r.base = dbRoom.Base
	r.water = dbRoom.Water
	r.admin = dbRoom.Admin
	r.startTime = dbRoom.StartTime
	r.endTime = dbRoom.EndTime
	r.CountDown = dbRoom.CountDown
	r.CountUp = dbRoom.CountUp
	r.RedCountDown = dbRoom.RedCountDown
	r.RedCountUp = dbRoom.RedCountUp
	r.RedDown = dbRoom.RedDown
	r.RedUp = dbRoom.RedUp
	r.Timeout = dbRoom.Timeout
	r.ScoreLimit = dbRoom.ScoreLimit
	r.Describe = dbRoom.Describe
	r.duration = dbRoom.Duration
	r.RedInterval = dbRoom.RedInterval
	r.reqStatus = dbRoom.ReqStatus
	r.CreateAt = dbRoom.CreateAt
	r.ActiveAt = dbRoom.ActiveAt

	//	r.describe= dbRoom.Describe
	r.users = make(map[string]*Player)
	for i, u := range rlist {
		rlist[i].Player.NicName = decodeNic(rlist[i].Player.NicName)
		r.users[u.Uid] = &(rlist[i].Player)
	}
	if p, ok := r.users[r.admin]; ok {
		p.Role = Role_Admin
	}
	return nil

}

func (r *Room) DeleteDB() error {
	dbRoom := &DBRoom{}
	dbRoom.Id = r.id
	trans, err := dBEngine.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			trans.Rollback()
		} else {
			trans.Commit()
		}
	}()
	_, err = trans.Delete(dbRoom)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`delete from %s where room_id="%d"`, (&RoomUser{}).TableName(), r.id)

	_, err = trans.Exec(query)

	return err

}

func (p *Player) Insert(rid int, uid string) error {
	u := &RoomUser{
		RoomId: rid,
		Uid:    uid,
		Player: *p,
	}
	u.NicName = encodeNic(u.NicName)
	return u.Insert(dBEngine)
}

func (p *Player) Update(rid int, uid string) error {
	u := &RoomUser{
		RoomId: rid,
		Uid:    uid,
	}
	err := u.Fetch(dBEngine)
	if err != nil {
		return err
	}
	u.Player = *p
	u.NicName = encodeNic(u.NicName)
	return u.Update(dBEngine)
}

func Record(rid int, rname string, master string, body interface{}, water int, jushu int) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r := &DBRecord{
		RoomId:   rid,
		RoomName: rname,
		Master:   master,
		Water:    water,
		Jusu:     jushu,
		Body:     encodeNic(string(bs)),
	}
	return r.Insert(dBEngine)
}

func BillList(rid int, page int, limit int) (interface{}, error) {
	rs := &DBRecord{
		RoomId: rid,
	}
	return rs.List(dBEngine, page, limit)
}

var dBEngine *gorp.DbMap

func DBEngineInit() *gorp.DbMap {
	connectionString := fmt.Sprintf(
		"%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FShanghai&interpolateParams=true",
		"root:11223344Asdf", //root:11223344Asdf
		"127.0.0.1",
		"3306",
		"game", //game
	)
	//connectionString := ""
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(30)
	dbEngine := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	userTable := DBUser{}
	recordTable := DBRecord{}
	roomTable := DBRoom{}
	roomUserTable := RoomUser{}
	reqTable := DBRoomPost{}
	redTable := DBRed{}
	dbEngine.AddTableWithName(userTable, userTable.TableName()).SetKeys(false, "Id", "NickName")
	dbEngine.AddTableWithName(recordTable, recordTable.TableName()).SetKeys(true, "Id")
	dbEngine.AddTableWithName(roomTable, roomTable.TableName()).SetKeys(false, "Id")
	dbEngine.AddTableWithName(roomUserTable, roomUserTable.TableName()).SetKeys(true, "Id")
	dbEngine.AddTableWithName(reqTable, reqTable.TableName()).SetKeys(false, "UserId")
	dbEngine.AddTableWithName(redTable, redTable.TableName()).SetKeys(true, "Id")
	return dbEngine
}

func RoomInit(db gorp.SqlExecutor) {
	roomlist := []int{}
	qurey := fmt.Sprintf(`select id from %s `, (&DBRoom{}).TableName())

	_, err := db.Select(&roomlist, qurey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(qurey, roomlist)
	for _, rid := range roomlist {
		room := &Room{
			id: rid,
		}
		err := room.Fetch()
		if err != nil {
			//fmt.Println(err)
			continue
		}

		if room.reqStatus {
			room.active = false
		} else {
			room.active = true
		}

		room.Active()
		RoomList[rid] = room
	}
}

func CheckAndCommit(db *gorp.Transaction, err *error) {
	if err == nil {
		db.Commit()
	} else {
		db.Rollback()
	}
}

func encodeNic(nic string) string {
	return base64.StdEncoding.EncodeToString([]byte(nic))
}

func decodeNic(nic string) string {
	bat, err := base64.StdEncoding.DecodeString(nic)
	if err != nil {
		return nic
	}
	return string(bat)
}
