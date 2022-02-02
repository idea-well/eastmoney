package eastmoney

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type banKuaiRes struct {
	Data struct {
		Diff BanKuaiDatas `json:"diff"`
	} `json:"data"`
}

type BanKuaiData struct {
	BkCode string `json:"f12"`
	BkName string `json:"f14"`
}

type BanKuaiDatas map[string]*BanKuaiData

func (ds BanKuaiDatas) indexByName() map[string]*BanKuaiData {
	map_ := make(map[string]*BanKuaiData)
	for _, data := range ds {
		map_[data.BkName] = data
	}
	return map_
}

const banKuaiApi = "https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=10000&fs=m:90"

func BanKuai() (BanKuaiDatas, error) {
	var res = new(banKuaiRes)
	resp, err := http.Get(banKuaiApi)
	return res.Data.Diff, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
