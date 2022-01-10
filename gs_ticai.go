package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type gsTiCaiRes struct {
	Data []*GsTiCaiData `json:"Data"`
}

type GsTiCaiData struct {
	BK                string       `json:"BK"`                // 所属板块
	BUSINSCOPE        string       `json:"BUSINSCOPE"`        // 经营范围
	COMPPROFILE       string       `json:"COMPPROFILE"`       // 公司简介
	COMPSCOPE         string       `json:"COMPSCOPE"`         // 公司沿革
	MAINBUSIN         string       `json:"MAINBUSIN"`         // 主营业务
	ZYCP              string       `json:"ZYCP"`              // 主营产品
	SECURITYCODE      string       `json:"SECURITYCODE"`      // 股票代码
	SECURITYSHORTNAME string       `json:"SECURITYSHORTNAME"` // 股票简称
	DIBK              BanKuaiDatas // 地域板块
	HYBK              BanKuaiDatas // 行业板块
	GNBK              BanKuaiDatas // 概念板块
}

func (d *GsTiCaiData) fillBK(dy, hy, gn map[string]*BanKuaiData) {
	bks := strings.Split(d.BK, ",")
	for _, bk := range bks {
		if v, ok := dy[bk]; ok {
			d.DIBK = append(d.DIBK, v)
		}
		if v, ok := hy[bk]; ok {
			d.HYBK = append(d.HYBK, v)
		}
		if v, ok := gn[bk]; ok {
			d.GNBK = append(d.GNBK, v)
		}
	}
}

type GsTiCaiDatas []*GsTiCaiData

func (gs GsTiCaiDatas) fetchBanKuai() error {
	var errs = new(Errors)
	dyBK, err1 := DiYuBanKuai()
	hyBK, err2 := HangYeBanKuai()
	gnBK, err3 := GaiNianBanKuai()
	errs.add(err1, err2, err3)
	return callWithoutErr2(errs.first(), func() {
		dyIndex := dyBK.indexByName()
		hyIndex := hyBK.indexByName()
		gnIndex := gnBK.indexByName()
		for _, g := range gs {
			g.fillBK(dyIndex, hyIndex, gnIndex)
		}
	})
}

const gsTiCaiApi = "https://data.eastmoney.com/dataapi/gstc/search"

func GsTiCai() ([]*GsTiCaiData, error) {
	datas, err := doFetchGsTiCai()
	return datas, callWithoutErr(err, datas.fetchBanKuai)
}

func doFetchGsTiCai() (GsTiCaiDatas, error) {
	var dataLocker = new(sync.Mutex)
	var errs, page = make(Errors, 0), 1
	var datass = make(GsTiCaiDatas, 0)
	line := pool.AddLine(func(i interface{}) {
		datas, err := doFetchGsTiCaiPage(i.(int))
		errs.add(callWithoutErr2(err, func() {
			dataLocker.Lock()
			datass = append(datass, datas...)
			dataLocker.Unlock()
		}))
	})
	for ; page <= 95; page++ {
		line.Submit(page)
	}
	line.Wait() // wait done
	for {
		datas, err := doFetchGsTiCaiPage(page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

func doFetchGsTiCaiPage(pageNo int) (GsTiCaiDatas, error) {
	var res = new(gsTiCaiRes)
	page := fmt.Sprintf("?ps=50&p=%d", pageNo)
	resp, err := http.Get(gsTiCaiApi + page)
	return res.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
