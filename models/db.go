package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"strings"

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
	err := cemsdk.CreateAccount(openid, pas, nicname)
	if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("%s exists", openid)) {
		return "", err
	}
	defer func() {
		if err != nil {
			cemsdk.DeleteAccount(openid)
		}
	}()
	u.Password = pas
	u.NickName = ""
	var trans *gorp.Transaction
	trans, err = dBEngine.Begin()
	defer CheckAndCommit(trans, err)
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

	token := ""
	token, err = cemsdk.GetUserToken(u.Id, u.Password)
	return token, err

}

type DBRecord struct {
	Id        int64     `db:"id","primarykey","autoincrement" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	RoomId    string    `db:"room_id" json:"room_id"`
	RoomName  string    `db:"room_name" json:"room_name"`
	Master    string    `db:"master" json:"master"`
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

func (r *DBRecord) List(db gorp.SqlExecutor, page int, size int) (interface{}, error) {
	list := []*DBRecord{}
	query := fmt.Sprintf(`select * from %s where room_id = "%s"`, r.TableName(), r.RoomId)

	_, err := db.Select(&list, query)
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
	//	Id        int64     `db:"id","primarykey","autoincrement" json:"id"`
	RoomId    string    `db:"room_id","primarykey"`
	RoomName  string    `db:"room_name"`
	Base      int       `db:"base"`
	Water     int       `db:"water"`
	Admin     string    `db:"admin"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
	//	Describe  string    `db:"discrible"`
	Scope
}

func (room *DBRoom) TableName() string {
	return "db_room"
}

func (room *DBRoom) Insert(db gorp.SqlExecutor) error {
	return db.Insert(room)
}

func (room *DBRoom) Fetch(db gorp.SqlExecutor) error {
	return db.SelectOne(room, fmt.Sprintf(`select * from %s where room_id="%s" `, room.TableName(), room.RoomId))
}

func (room *DBRoom) Update(db gorp.SqlExecutor) error {
	_, err := db.Update(room)
	return err
}

type RoomUser struct {
	Id     int64  `db:"id","primarykey","autoincrement" json:"id"`
	Uid    string `db:"uid"`
	RoomId string `db:"room_id"`
	Player
}

func (ruser *RoomUser) TableName() string {
	return "room_user"
}

func (ruser *RoomUser) Insert(db gorp.SqlExecutor) error {
	return db.Insert(ruser)
}

func (ruser *RoomUser) Fetch(db gorp.SqlExecutor) error {
	return db.SelectOne(ruser, fmt.Sprintf(`select * from %s where uid="%s" and room_id="%s"`, ruser.TableName(), ruser.Uid, ruser.RoomId))
}

func (ruser *RoomUser) Update(db gorp.SqlExecutor) error {
	_, err := db.Update(ruser)
	return err
}

func (r *Room) Insert() error {
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
		RoomId:    r.id,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		StartTime: r.startTime,
		EndTime:   r.endTime,
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
		},
	}
	err = dbRoom.Insert(trans)
	if err != nil {
		return err
	}
	for uid, player := range r.users {
		ruser := &RoomUser{
			Uid:    uid,
			RoomId: r.id,
			Player: *player,
		}
		err = ruser.Insert(trans)
		if err != nil {
			return err
		}
	}
	return err
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
		RoomId:    r.id,
		RoomName:  r.name,
		Base:      r.base,
		Water:     r.water,
		Admin:     r.admin,
		StartTime: r.startTime,
		EndTime:   r.endTime,
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
		},
	}
	err = dbRoom.Update(trans)
	return err
}

func (r *Room) Fetch() error {
	dbRoom := &DBRoom{}
	dbRoom.RoomId = r.id
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
	rlist := []RoomUser{}
	_, err = trans.Select(&rlist, fmt.Sprintf(`select * from %s where  room_id="%s"`, (&RoomUser{}).TableName(), r.id))
	if err != nil {
		return err
	}

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
	//	r.describe= dbRoom.Describe
	r.users = make(map[string]*Player)
	for _, u := range rlist {
		r.users[u.Uid] = &u.Player
	}
	return nil

}

func (r *Room) DeleteDB() error {
	dbRoom := &DBRoom{}
	dbRoom.RoomId = r.id
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
	query := fmt.Sprintf(`delete from %s where room_id="%s"`, (&RoomUser{}).TableName(), r.id)

	_, err = trans.Exec(query)

	return err

}

func (p *Player) Insert(rid string, uid string) error {
	u := &RoomUser{
		RoomId: rid,
		Uid:    uid,
		Player: *p,
	}
	return u.Insert(dBEngine)
}

func (p *Player) Update(rid string, uid string) error {
	u := &RoomUser{
		RoomId: rid,
		Uid:    uid,
	}
	err := u.Fetch(dBEngine)
	if err != nil {
		return err
	}
	u.Player = *p
	return u.Update(dBEngine)
}

func Record(rid string, rname string, master string, body interface{}) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r := &DBRecord{
		RoomId:   rid,
		RoomName: rname,
		Master:   master,
		Body:     string(bs),
	}
	return r.Insert(dBEngine)
}

func BillList(rid string, page int, limit int) (interface{}, error) {
	rs := &DBRecord{
		RoomId: rid,
	}
	return rs.List(dBEngine, page, limit)
}

var dBEngine *gorp.DbMap

func DBEngineInit() *gorp.DbMap {
	connectionString := fmt.Sprintf(
		"%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FShanghai&interpolateParams=true",
		"root:123456", //root:11223344Asdf
		"127.0.0.1",
		"3306",
		"test", //game
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
	dbEngine.AddTableWithName(userTable, userTable.TableName()).SetKeys(false, "Id", "NickName")
	dbEngine.AddTableWithName(recordTable, recordTable.TableName()).SetKeys(true, "Id")
	dbEngine.AddTableWithName(roomTable, roomTable.TableName()).SetKeys(false, "RoomId")
	dbEngine.AddTableWithName(roomUserTable, roomUserTable.TableName()).SetKeys(true, "Id")
	dbEngine.AddTableWithName(reqTable, reqTable.TableName()).SetKeys(false, "UserId")
	return dbEngine
}

func RoomInit(db gorp.SqlExecutor) {
	roomlist := []string{}
	qurey := fmt.Sprintf(`select room_id from %s `, (&DBRoom{}).TableName())
	_, err := db.Select(&roomlist, qurey)
	if err != nil {
		fmt.Println(err)
	}
	for _, rid := range roomlist {
		room := &Room{
			id: rid,
		}
		err := room.Fetch()
		if err != nil {
			fmt.Println(err)
			continue
		}
		RoomList[rid] = room
	}
}

func CheckAndCommit(db *gorp.Transaction, err error) {
	if err == nil {
		db.Commit()
	} else {
		db.Rollback()
	}
}
