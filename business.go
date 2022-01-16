package eastmoney

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type businessRes struct {
	Data []*BusinessData `json:"zygcfx"`
}

type BusinessData struct {
	ITEM_NAME            string  // 主营构成
	MAINOP_TYPE          string  // 1行业 2产品 3地区
	GROSS_RPOFIT_RATIO   float64 // 毛利率
	MAIN_BUSINESS_COST   float64 // 主营成本
	MAIN_BUSINESS_INCOME float64 // 主营收入
	MAIN_BUSINESS_RPOFIT float64 // 主营利润
	MBC_RATIO            float64 // 成本比例
	MBI_RATIO            float64 // 收入比例
	MBR_RATIO            float64 // 利润比例
	REPORT_DATE          string  // 报告日期
	SECURITY_CODE        string  // 股票代码
}

const businessApi = "https://emweb.securities.eastmoney.com/BusinessAnalysis/PageAjax"

func Business(code string) ([]*BusinessData, error) {
	var res = new(businessRes)
	resp, err := http.Get(businessApi + "?code=" + code)
	return res.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
