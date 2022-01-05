package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type guDongRes struct {
	Result struct {
		Data GuDongDatas `json:"data"`
	} `json:"result"`
}

type GuDongData struct {
	//FREE_HOLDNUM_RATIO: 1.44849464764
	//FREE_RATIO_QOQ: "-62.43717617"
	//HOLDER_STATE: null
	//HOLDER_STATEE: "减仓"
	//HOLD_CHANGE: "-4252789"
	//HOLD_NUM: 2558520
	//HOLD_NUM_CHANGE: "-4252789"
	//HOLD_RATIO: 1.1796
	//HOLD_RATIO_CHANGE: -1.9606
	//IS_HOLDORG: "1"
	//IS_MAX_REPORTDATE: "1"
	//IS_REPORT: "1"
	//LISTING_STATE: "0"
	//UPDATE_DATE: "2021-10-29 00:00:00"
	//
	//
	//
	//
	//DIRECTION string // 增减持
	//END_DATE string // 报告期
	//HOLDER_CODE string // 持有人编码
	//HOLDER_NAME string // 持有人姓名
	//HOLDER_NATURE string // 持有人属性
	//
	//HOLDER_NEW: "王欣"
	//HOLDER_NEWTYPE: "个人"
	//HOLDER_TYPE_ORG: "个人"
	//HOLDNUM_CHANGE_NAME: "不变"
	//HOLDNUM_CHANGE_RATIO: 0
	//HOLD_CHANGE: "不变"
	//
	//HOLD_NUM string // 持有数量
	//HOLD_NUM_CHANGE int64 // 增减数量
	//HOLD_RATIO float64 // 持股比例
	//HOLD_RATIO_CHANGE float64 //
	//HOLD_RATIO_YOY: 0
	//IS_MAX_REPORTDATE: "1"
	//MXID: "61f99d5d7b21c9fa4db52440054cb845"
	//NOTICE_DATE: "2021-10-29 00:00:00"
	//ORG_CODE: "10145994"
	//RANK: 3
	//REPORT_DATE_NAME: "2021年三季报"
	//REPORT_TYPE: "定期报告"
	//SECUCODE: "300547.SZ"
	//SECURITY_CODE: "300547"
	//SECURITY_NAME_ABBR: "川环科技"
	//SECURITY_TYPE_CODE: "058001001"
	//SHARES_TYPE: "流通A股,限售流通A股"
	//XZCHANGE: 0
}

type GuDongDatas []*GuDongData

// GuDong 十大股东
func GuDong(code, date string) (GuDongDatas, error) {
	return doFetchGuDong("RPT_DMSK_HOLDERS", code, date)
}

// FreeGuDong 十大流通股东
func FreeGuDong(code, date string) (GuDongDatas, error) {
	return doFetchGuDong("RPT_F10_EH_FREEHOLDERS", code, date)
}

const guDongApi = "https://datacenter-web.eastmoney.com/api/data/v1/get"

func doFetchGuDong(report, code, date string) (GuDongDatas, error) {
	var datass, page = make(GuDongDatas, 0), 1
	for {
		datas, err := doFetchGuDongPage(report, code, date, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

func doFetchGuDongPage(report, code, date string, pageNo int) (GuDongDatas, error) {
	var res = new(guDongRes)
	page := fmt.Sprintf("&pageSize=50&pageNumber=%d", pageNo)
	query := fmt.Sprintf("?reportName=%s&columns=ALL&filter=", report)
	if date != "" && strings.ToLower(date) != "all" {
		query += url.QueryEscape(fmt.Sprintf("(END_DATE='%s')", date))
	}
	if code != "" && strings.ToLower(code) != "all" {
		query += url.QueryEscape(fmt.Sprintf("(SECURITY_CODE=\"%s\")", code))
	}
	resp, err := http.Get(guDongApi + query + page)
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
