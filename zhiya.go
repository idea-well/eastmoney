package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type zhiYaRes struct {
	Result struct {
		Data []*ZhiYaData `json:"data"`
	} `json:"result"`
}

type ZhiYaData struct {
	SECUCODE                     string  `json:"SECUCODE"`                     // 股票代码
	SECURITY_CODE                string  `json:"SECURITY_CODE"`                // 股票代码
	SECURITY_NAME_ABBR           string  `json:"SECURITY_NAME_ABBR"`           // 股票名称
	PLEDGE_DEAL_NUM              int     `json:"PLEDGE_DEAL_NUM"`              // 质押笔数
	PLEDGE_MARKET_CAP            float64 `json:"PLEDGE_MARKET_CAP"`            // 质押市值
	PLEDGE_RATIO                 float64 `json:"PLEDGE_RATIO"`                 // 质押比例
	REPURCHASE_BALANCE           float64 `json:"REPURCHASE_BALANCE"`           // 质押手数
	REPURCHASE_LIMITED_BALANCE   float64 `json:"REPURCHASE_LIMITED_BALANCE"`   // 限售质押手数
	REPURCHASE_UNLIMITED_BALANCE float64 `json:"REPURCHASE_UNLIMITED_BALANCE"` // 不限售质押手数
}

func ZhiYa(date string) ([]*ZhiYaData, error) {
	var datass, page = make([]*ZhiYaData, 0), 1
	for {
		datas, err := doFetchZhiYaPage(date, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

const zhiYaApi = "https://datacenter-web.eastmoney.com/api/data/v1/get?pageSize=50&reportName=RPT_CSDC_LIST&columns=ALL"

func doFetchZhiYaPage(date string, page int) ([]*ZhiYaData, error) {
	var res = new(zhiYaRes)
	filter := url.QueryEscape(fmt.Sprintf("(TRADE_DATE='%s')", date))
	query := fmt.Sprintf("&pageNumber=%d&filter=%s", page, filter)
	resp, err := http.Get(zhiYaApi + query)
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
