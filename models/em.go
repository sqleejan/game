package models

import (
	"github.com/links123com/emsdk"
)

const (
	hxId     = "YXA6y0gUUCv1EeeF1FEN-eNstw"
	hxSecret = "YXA6GmZWVoGxr-MzUVRUWWZIKA8Bsmo"
)

var cemsdk *emsdk.Client

func newEm() (*emsdk.Client, error) {
	return emsdk.New("1189170428115308", "jifenyouxi", hxId, hxSecret)
}
