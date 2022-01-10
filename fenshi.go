package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type fengShiRes struct {
	Data struct {
		Data []*FengShiData `json:"data"`
	} `json:"data"`
}

type FengShiData struct {
	Type   int `json:"bs"` // 1买 2卖
	Time   int `json:"t"`  // 时间
	Price  int `json:"p"`  // 价格
	Volume int `json:"v"`  // 手数
}

func FenShi(code string, market int) ([]*FengShiData, error) {
	var datass, page = make([]*FengShiData, 0), 1
	for {
		datas, err := doFetchFenShiPage(code, market, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

const fenShiApi = "http://push2ex.eastmoney.com/getStockFenShi" +
	"?pagesize=100&ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wzfscj&sort=1&ft=1"

func doFetchFenShiPage(code string, market, page int) ([]*FengShiData, error) {
	var res = new(fengShiRes)
	query := fmt.Sprintf("&code=%s&market=%d&pageindex=%d", code, market, page)
	resp, err := http.Get(fenShiApi + query)
	return res.Data.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
