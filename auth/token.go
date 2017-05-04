package auth

import (
	"net/http"
	"time"

	"github.com/chanxuehong/wechat.v2/mp/oauth2"
	authClient "github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/skip2/go-qrcode"
)

var (
	rootkey  = []byte("renmindemingyi")
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

func Gen(rid string, openid string, wxToken string) string {

	// Create the Claims
	claims := jwt.StandardClaims{
	//WXToken: wxToken,
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 10).Unix()
	claims.Id = openid
	claims.Audience = rid
	claims.Issuer = wxToken
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

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return &MyCustomClaims{*claims}, nil
	} else {
		return nil, err
	}
}

func QRCode(rid string) ([]byte, error) {
	return qrcode.Encode("https://example.org/auth/"+rid, qrcode.Medium, 256)
}

const (
	appId       = ""
	appSecret   = ""
	redirectURI = ""
	scope       = ""
	state       = ""
)

func CodeUrl(rid string) string {
	uri := redirectURI
	if rid != "" {
		uri = uri + "&roomid=" + rid
	}
	return oauth2.AuthCodeURL(appId, uri, scope, state)
}

func init() {
	endpoint := oauth2.NewEndpoint(appId, appSecret)
	wxClient = &authClient.Client{
		Endpoint:   endpoint,
		HttpClient: http.DefaultClient,
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
