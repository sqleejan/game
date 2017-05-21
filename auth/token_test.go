package auth

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {

	fmt.Println(CodeUrl(123, false))
	return
	// fmt.Println(Gen("dddddd", "sfsdfsdfdfdsfdsfdsfdsfsdfdsfds", "xui"))
	myc, err := Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsZWVqYW4iLCJleHAiOjE0OTUyMjE2NTcsImp0aSI6Im9mYzNQdnBTZ2tlVjc1MGNmZ0M4Q2txNWVtREkiLCJpc3MiOiJodHRwOi8vd3gucWxvZ28uY24vbW1vcGVuL1BpYWp4U3FCUmFFSXplMkR2Y0VlNXQ5Z2d6akptcnMzNmhSSWZ0aWFMTlozUHVibk5FSE03b3FRZEVCbHV4T1o1VWZEcEtMU0dIQXluOXJqc2ljUDV1MWp3LzAifQ.JIHmA1szQs251x8VRqvWNoKVeuKM-SVz4Lx6p5AZcCQ")
	//myc, err := WXClaim("031P7vTv1hJRfd0IuoUv184qTv1P7vTS")
	fmt.Println(err)
	if myc != nil {
		fmt.Println(myc.Audience, myc.Issuer)
	}

	return
	ts := Gen("room1", "openid1", "YWMt4_B4Li-tEeeV6othz6ZiH8tIFFAr9RHnhdRRDfnjbLcXGQygL1YR57VbrxqfTP69AwMAAAFbzEiCpABPGgDDfUr_f64kQSqrrraMjp4zmSfSMGGU_lyCOftLMe_62Q")
	fmt.Println(len(ts))
	cc, _ := Parse(ts)
	fmt.Println(cc.Id)
}
