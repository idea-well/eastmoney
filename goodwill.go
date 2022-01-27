package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type goodwillRes struct {
	Result struct {
		Data []*GoodWillData `json:"data"`
	} `json:"result"`
}

type GoodWillData struct {
	SECURITY_CODE string  `json:"SECURITY_CODE"` // 股票代码
	GOODWILL      float64 `json:"GOODWILL"`      // 商誉值(元)
	NOTICE_DATE   string  `json:"NOTICE_DATE"`   // 公告日期
}

func GoodWill() ([]*GoodWillData, error) {
	datass, page := make([]*GoodWillData, 0), 1
	for {
		datas, err := doFetchGoodsWillPage(page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

const api = "https://datacenter-web.eastmoney.com/api/data/get?sty=ALL&type=RPT_GOODWILL_STOCKDETAILS"

func doFetchGoodsWillPage(pageNo int) ([]*GoodWillData, error) {
	var res = new(goodwillRes)
	page := fmt.Sprintf("ps=50&p=%d", pageNo)
	resp, err := http.Get(api + page)
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
