package auth

import (
	"net/http"
	"time"

	"fmt"
	"os"

	"sync"

	"github.com/chanxuehong/wechat.v2/mp/oauth2"
	authClient "github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/skip2/go-qrcode"
)

type freshToken struct {
	sync.RWMutex
	list map[string]time.Time
}

func (f *freshToken) Add(uid string) {
	f.Lock()
	defer f.Unlock()
	f.list[uid] = time.Now().Add(time.Second * 100)
	fmt.Println("add", f.list[uid], uid)
}

func (f *freshToken) Active(uid string) bool {
	fmt.Println("active", f.list[uid], uid)
	f.Lock()
	defer f.Unlock()
	if t, ok := f.list[uid]; ok {
		now := time.Now()
		if t.After(now) {
			delete(f.list, uid)
			return false
		}
		f.list[uid] = now.Add(time.Second * 100)
		return true
	}
	return false
}

var (
	refreshToken = &freshToken{
		list: map[string]time.Time{},
	}
	rootkey  = []byte("wcx")
	wxClient *authClient.Client
)

type MyCustomClaims struct {
	//WXToken string `json:"wx_token"`
	jwt.StandardClaims
}

func (mc *MyCustomClaims) Token() string {
	claims := mc.StandardClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ob, _ := token.SignedString(rootkey)
	return ob
}

func Gen(rid string, openid string, nicname string) string {

	// Create the Claims
	claims := jwt.StandardClaims{
	//WXToken: wxToken,
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	claims.Id = openid
	claims.Audience = nicname
	//claims.Issuer = wxToken
	claims.Subject = rid

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ob, _ := token.SignedString(rootkey)
	return ob
}

func Parse(tokenString string) (*MyCustomClaims, error) {

	// sample token is expired.  override time so it parses as valid

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return rootkey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		if !refreshToken.Active(claims.Id) {
			fmt.Println(claims.Id, "token is expired!")
			return nil, fmt.Errorf("token is expired!")
		}
		return &MyCustomClaims{*claims}, nil
	} else {
		return nil, err
	}
}

func QRCode(rid int) ([]byte, error) {
	return qrcode.Encode(loginURI+fmt.Sprintf("?roomid=%d&time=%v", rid, time.Now().Unix()), qrcode.Medium, 256)
}

var (
	appId       = "wx9f85ad10c59fc23c"
	appSecret   = "102576667f4d60f65c0d1405c8a04d4e"
	loginURI    = "http://game.highlifes.com/v1/auth/wx/checkin"
	redirectURI = "http://game.highlifes.com/v1/auth/wx/code"
	scope       = "snsapi_userinfo"
	state       = ""
)

func CodeUrl(rid int, ext bool) string {
	// uri := redirectURI
	// if rid != "" {
	// 	uri = uri + "&roomid=" + rid
	// }
	roomid := fmt.Sprintf("%d", rid)
	red := ""
	// if ext {
	// 	red = oauth2.AuthExtURL(appId, redirectURI, scope, roomid)
	// } else {
	red = oauth2.AuthCodeURL(appId, redirectURI, scope, roomid)
	// }

	fmt.Println(red)
	return red
}

func CodeUrlTest(rid int, ext bool) string {
	// uri := redirectURI
	// if rid != "" {
	// 	uri = uri + "&roomid=" + rid
	// }
	roomid := fmt.Sprintf("%d", rid)
	red := ""
	if ext {
		red = oauth2.AuthExtURL(appId, "http://game.highlifes.com/v1/auth/wx/codetest", scope, roomid)
	} else {
		red = oauth2.AuthCodeURL(appId, "http://game.highlifes.com/v1/auth/wx/codetest", scope, roomid)
	}

	fmt.Println(red)
	return red
}

func init() {
	endpoint := oauth2.NewEndpoint(appId, appSecret)
	wxClient = &authClient.Client{
		Endpoint:   endpoint,
		HttpClient: http.DefaultClient,
	}
	appid := os.Getenv("APPID")
	appsecret := os.Getenv("APPSECRET")
	fmt.Println(appid, "111"+appsecret)
	if appid != "" {
		appId = appid
	}
	if appsecret != "" {
		appSecret = appsecret
	}

}

func WXClaim(code string) (*MyCustomClaims, error) {
	token, err := wxClient.ExchangeToken(code)
	if err != nil {
		return nil, err
	}
	user, err := oauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("nicname:", user.Nickname)
	claims := &MyCustomClaims{}
	//这个时间是token最长的维持时长
	claims.ExpiresAt = time.Now().Add(time.Hour * 36).Unix()
	claims.Id = token.OpenId
	claims.Audience = user.Nickname
	claims.Issuer = user.HeadImageURL
	//claims.Subject = rid
	refreshToken.Add(claims.Id)
	return claims, nil
}
