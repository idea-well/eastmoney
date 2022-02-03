package eastmoney

import (
	"encoding/json"
	"fmt"
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
	BkType int
}

type BanKuaiDatas map[string]*BanKuaiData

func (ds BanKuaiDatas) fillType(typ int) {
	for _, d := range ds {
		d.BkType = typ
	}
}

func (ds BanKuaiDatas) append(ds2 BanKuaiDatas) BanKuaiDatas {
	for key, val := range ds2 {
		ds[key] = val
	}
	return ds
}

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
	return ds3.append(ds2).append(ds1), err
}

func doFetchBanKuaiType(typ int) (BanKuaiDatas, error) {
	var res = new(banKuaiRes)
	res.Data.Diff = make(BanKuaiDatas)
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
