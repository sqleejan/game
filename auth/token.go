package auth

import (
	"net/http"
	"time"

	"fmt"
	"os"

	"github.com/chanxuehong/wechat.v2/mp/oauth2"
	authClient "github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/skip2/go-qrcode"
)

var (
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
		return &MyCustomClaims{*claims}, nil
	} else {
		return nil, err
	}
}

func QRCode(rid int) ([]byte, error) {
	return qrcode.Encode(loginURI+fmt.Sprintf("?roomid=%s", rid), qrcode.Medium, 256)
}

var (
	appId       = "wx9f85ad10c59fc23c"
	appSecret   = "102576667f4d60f65c0d1405c8a04d4e"
	loginURI    = "http://game.highlifes.com/v1/auth/wx/login"
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
	if ext {
		red = oauth2.AuthExtURL(appId, redirectURI, scope, roomid)
	} else {
		red = oauth2.AuthCodeURL(appId, redirectURI, scope, roomid)
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
	claims := &MyCustomClaims{}
	claims.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	claims.Id = token.OpenId
	claims.Audience = user.Nickname
	claims.Issuer = user.HeadImageURL
	//claims.Subject = rid
	return claims, nil
}
