package auth

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {

	//fmt.Println(CodeUrl(123, false))
	//return
	// fmt.Println(Gen("dddddd", "sfsdfsdfdfdsfdsfdsfdsfsdfdsfds", "xui"))
	myc, err := Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsZWVqYW4iLCJleHAiOjE0OTU0NzU5ODIsImp0aSI6Im9tbFZTd2VLZFNRQ28zR29rRkJYdjJHYXA5RVkiLCJpc3MiOiJodHRwOi8vd3gucWxvZ28uY24vbW1vcGVuL1BpYWp4U3FCUmFFSTh6N1FEYUdBY082dDNCT1FheUVYZWljaWJJdmRNUDduQkJHcXB2ZkVxVUZPSE9DaWNKQXRWaWNwaWNVd2ljVU5xU1UzODFWV0YzZGgzMHZsdy8wIn0.5TuYtCAqMM-mgHwvVhKRJt7RB-1bYLHWCnwnSQBlHsU")
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
