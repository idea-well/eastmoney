package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type banKuaiRes struct {
	Data struct {
		Diff BanKuaiDataMap `json:"diff"`
	} `json:"data"`
}

type BanKuaiData struct {
	BkCode string `json:"f12"`
	BkName string `json:"f14"`
	BkType int
}

// BanKuaiDataMap
type BanKuaiDataMap map[string]*BanKuaiData

func (dm BanKuaiDataMap) fillType(typ int) {
	for _, d := range dm {
		d.BkType = typ
	}
}

// BanKuaiDatas
type BanKuaiDatas []*BanKuaiData

func (ds BanKuaiDatas) indexByName() map[string]*BanKuaiData {
	map_ := make(map[string]*BanKuaiData)
	for _, data := range ds {
		map_[data.BkName] = data
	}
	return map_
}

const banKuaiApi = "https://push2.eastmoney.com/api/qt/clist/get?pn=1&pz=10000&fs=m:90+t:%d"

func BanKuai() (BanKuaiDatas, error) {
	ds1, er1 := doFetchBanKuaiType(1)
	ds2, er2 := doFetchBanKuaiType(2)
	ds3, er3 := doFetchBanKuaiType(3)
	err := firstError(er1, er2, er3)
	return mergeBanKuaiDatas(ds1, ds2, ds3), err
}

func mergeBanKuaiDatas(dss ...BanKuaiDataMap) BanKuaiDatas {
	ss := make(BanKuaiDatas, 0)
	for _, ds := range dss {
		for _, d := range ds {
			ss = append(ss, d)
		}
	}
	return ss
}

func doFetchBanKuaiType(typ int) (BanKuaiDataMap, error) {
	var res = new(banKuaiRes)
	res.Data.Diff = make(BanKuaiDataMap)
	resp, err := http.Get(fmt.Sprintf(banKuaiApi, typ))
	defer res.Data.Diff.fillType(typ)
	return res.Data.Diff, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
