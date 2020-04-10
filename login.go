package iport

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/woshihot/go-lib/utils/http"
	"github.com/woshihot/go-lib/utils/log"
)

func cloudLogin(url, machineNo, mac, licenseKey string) cloudLoginResp {

	b, err := http.Get(url, map[string]string{"machineNo": machineNo, "macAddr": mac, "licenseKey": licenseKey})
	log.DF(TagDebug, "cloud login resp = %s , err =%v\n", string(b), err)
	if nil != err {
		return cloudLoginResp{}
	}
	resp := new(cloudLoginResp)
	_ = jsoniter.Unmarshal(b, &resp)
	return *resp
}

type cloudLoginResp struct {
	Token     string `json:"token"`
	AesKey    string `json:"aesKey"`
	Life      int64  `json:"life"`
	Date      string `json:"date"`
	GroupName string `json:"groupName"`
}
