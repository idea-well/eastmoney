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

func (d *FengShiData) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", d.Type, d.Time, d.Price, d.Volume)
}

func FenShi(code string, market int) ([]*FengShiData, error) {
	var datass, page, size = make([]*FengShiData, 0), 0, 1000
	for {
		datas, err := doFetchFenShiPage(code, market, page, size)
		if len(datas) > 0 {
			datass = append(datass, datas...)
		}
		if err != nil || len(datas) < size {
			return datass, err
		}
		page++ // next page
	}
}

const fenShiApi = "http://push2ex.eastmoney.com/getStockFenShi" +
	"?ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wzfscj&sort=1&ft=1"

func doFetchFenShiPage(code string, market, page, size int) ([]*FengShiData, error) {
	var res = new(fengShiRes)
	query := fmt.Sprintf("&code=%s&market=%d&pageindex=%d&pagesize=%d", code, market, page, size)
	resp, err := http.Get(fenShiApi + query)
	return res.Data.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
