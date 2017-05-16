package auth

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	fmt.Println(CodeUrl("xxxx"))
	return
	ts := Gen("room1", "openid1", "YWMt4_B4Li-tEeeV6othz6ZiH8tIFFAr9RHnhdRRDfnjbLcXGQygL1YR57VbrxqfTP69AwMAAAFbzEiCpABPGgDDfUr_f64kQSqrrraMjp4zmSfSMGGU_lyCOftLMe_62Q")
	fmt.Println(len(ts))
	cc, _ := Parse(ts)
	fmt.Println(cc.Id)
}
