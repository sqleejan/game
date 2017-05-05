package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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

func CreateDBUser(openid, nicname string) error {
	u := &DBUser{
		Id: openid,
	}
	pas := string(Krand(8, KC_RAND_KIND_LOWER))
	err := cemsdk.CreateAccount(openid, pas, nicname)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			cemsdk.DeleteAccount(openid)
		}
	}()
	u.Password = pas
	u.NickName = nicname
	var trans *gorp.Transaction
	trans, err = dBEngine.Begin()
	defer CheckAndCommit(trans, err)
	err = u.Insert(trans)
	if err != nil {
		return err
	}
	// token := ""
	// token, err = cemsdk.GetUserToken(u.Id, u.Password)
	// return token, err
	return nil
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
	Id        int64     `db:"id","primarykey,"autoincrement" json:"id"`
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

func (r *DBRecord) List(db gorp.SqlExecutor) ([]*DBRecord, error) {
	list := []*DBRecord{}
	query := fmt.Sprintf(`select * from %s where room_id = "%s"`, r.TableName(), r.RoomId)

	_, err := db.Select(&list, query)
	return list, err
}

func (r *DBRecord) ListDBRoom(db gorp.SqlExecutor) ([]*DBRecord, error) {
	list := []*DBRecord{}
	query := fmt.Sprintf(`select distinct room_id, room_name,master from %s `, r.TableName())

	_, err := db.Select(&list, query)
	return list, err
}

func Record(rid string, rname string, master string, body interface{}) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r := &DBRecord{
		RoomId: rid,
		Body:   string(bs),
	}
	return r.Insert(dBEngine)
}

var dBEngine *gorp.DbMap

func DBEngineInit() *gorp.DbMap {
	connectionString := fmt.Sprintf(
		"%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Asia%%2FShanghai&interpolateParams=true",
		"root:123456",
		"127.0.0.1",
		"3306",
		"test",
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
	dbEngine.AddTableWithName(userTable, userTable.TableName()).SetKeys(false, "Id", "NickName")
	dbEngine.AddTableWithName(recordTable, recordTable.TableName()).SetKeys(true, "Id")
	return dbEngine
}

func CheckAndCommit(db *gorp.Transaction, err error) {
	if err == nil {
		db.Commit()
	} else {
		db.Rollback()
	}
}