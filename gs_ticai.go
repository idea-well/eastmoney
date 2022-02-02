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
	BK                string         `json:"BK"`                // 所属板块
	BUSINSCOPE        string         `json:"BUSINSCOPE"`        // 经营范围
	COMPPROFILE       string         `json:"COMPPROFILE"`       // 公司简介
	COMPSCOPE         string         `json:"COMPSCOPE"`         // 公司沿革
	MAINBUSIN         string         `json:"MAINBUSIN"`         // 主营业务
	ZYCP              string         `json:"ZYCP"`              // 主营产品
	SECURITYCODE      string         `json:"SECURITYCODE"`      // 股票代码
	SECURITYSHORTNAME string         `json:"SECURITYSHORTNAME"` // 股票简称
	LISTINGDATE       string         `json:"LISTINGDATE"`       // 上市时间
	LTGB              string         `json:"ltgb"`              // 流通股本
	ZGB               string         `json:"zgb"`               // 总股本
	BanKuai           []*BanKuaiData // 所属板块
}

func (d *GsTiCaiData) fillBK(idx map[string]*BanKuaiData) {
	bks := strings.Split(d.BK, ",")
	for _, bk := range bks {
		if v, ok := idx[bk]; ok {
			d.BanKuai = append(d.BanKuai, v)
		}
	}
}

type GsTiCaiDatas []*GsTiCaiData

func (gs GsTiCaiDatas) fetchBanKuai() error {
	datas, err := BanKuai()
	return callWithoutErr2(err, func() {
		dataIdx := datas.indexByName()
		for _, g := range gs {
			g.fillBK(dataIdx)
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
