package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type longHubRes struct {
	Result struct {
		Data LongHubDatas `json:"data"`
	} `json:"result"`
}

type LongHubData struct {
	ACCUM_AMOUNT       float64              `json:"ACCUM_AMOUNT"`       // 总成交额
	BILLBOARD_DEAL_AMT float64              `json:"BILLBOARD_DEAL_AMT"` // 龙虎榜成交额
	BILLBOARD_BUY_AMT  float64              `json:"BILLBOARD_BUY_AMT"`  // 龙虎榜买入额
	BILLBOARD_SELL_AMT float64              `json:"BILLBOARD_SELL_AMT"` // 龙虎榜卖出额
	BILLBOARD_NET_AMT  float64              `json:"BILLBOARD_NET_AMT"`  // 龙虎榜净买额
	DEAL_AMOUNT_RATIO  float64              `json:"DEAL_AMOUNT_RATIO"`  // 成交额占总成交比例
	FREE_MARKET_CAP    float64              `json:"FREE_MARKET_CAP"`    // 流通市值
	SECUCODE           string               `json:"SECUCODE"`           // 股票代码
	SECURITY_CODE      string               `json:"SECURITY_CODE"`      // 股票代码
	SECURITY_NAME_ABBR string               `json:"SECURITY_NAME_ABBR"` // 股票名称
	TRADE_DATE         string               `json:"TRADE_DATE"`         // 交易日期
	TURNOVERRATE       float64              `json:"TURNOVERRATE"`       // 换手率
	CHANGE_RATE        float64              `json:"CHANGE_RATE"`        // 涨跌幅
	CLOSE_PRICE        float64              `json:"CLOSE_PRICE"`        // 收盘价
	EXPLANATION        string               `json:"EXPLANATION"`        // 上榜理由
	BUY_DETAILS        []*LongHubDetailData // 买入明细
	SELL_DETAILS       []*LongHubDetailData // 卖出明细
}

func (d *LongHubData) fetchBuyDetails() (err error) {
	d.BUY_DETAILS, err = fetchLongHubBuyDetail(d.TRADE_DATE[0:10], d.SECURITY_CODE)
	return
}

func (d *LongHubData) fetchSellDetails() (err error) {
	d.SELL_DETAILS, err = fetchLongHubSellDetail(d.TRADE_DATE[0:10], d.SECURITY_CODE)
	return
}

type LongHubDatas []*LongHubData

// STDatas ST龙虎榜
func (ds LongHubDatas) STDatas() LongHubDatas {
	datas := make(LongHubDatas, 0)
	for _, d := range ds {
		if strings.Contains(d.SECURITY_NAME_ABBR, "ST") {
			datas = append(datas, d)
		}
	}
	return datas
}

// HSDatas 换手率龙虎榜
func (ds LongHubDatas) HSDatas() LongHubDatas {
	return ds
}

// ZDDatas 涨跌幅龙虎榜
func (ds LongHubDatas) ZDDatas() LongHubDatas {
	return ds
}

func (ds LongHubDatas) fetchDetail() error {
	errs := make(Errors, 0)
	line := pool.AddLine(func(i interface{}) {
		errs.add(
			ds[i.(int)].fetchBuyDetails(),
			ds[i.(int)].fetchSellDetails(),
		)
	})
	for index := range ds {
		line.Submit(index)
	}
	line.Wait() // wait done
	return errs.first()
}

const longHubApi = "https://datacenter-web.eastmoney.com/api/data/v1/get"

func DateLongHub(date string) (LongHubDatas, error) {
	ds, err := fetchLongHub(date)
	return ds, callWithoutErr(err, ds.fetchDetail)
}

func fetchLongHub(date string) (LongHubDatas, error) {
	var datass, page = make(LongHubDatas, 0), 1
	for {
		datas, err := fetchLongHubPage(date, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

func fetchLongHubPage(date string, pageNo int) (LongHubDatas, error) {
	var res = new(longHubRes)
	page := fmt.Sprintf("&pageSize=50&pageNumber=%d", pageNo)
	filter := fmt.Sprintf("(TRADE_DATE>='%s')(TRADE_DATE<='%s')", date, date)
	query := "?reportName=RPT_DAILYBILLBOARD_DETAILS&columns=ALL"
	resp, err := http.Get(longHubApi + query + page + "&filter=" + url.QueryEscape(filter))
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}

type longHubDetailRes struct {
	Result struct {
		Data []*LongHubDetailData `json:"data"`
	} `json:"result"`
}

type LongHubDetailData struct {
	BUY              float64 `json:"BUY"`              // 买入额
	SELL             float64 `json:"SELL"`             // 卖出额
	NET              float64 `json:"NET"`              // 净买额
	OPERATEDEPT_NAME string  `json:"OPERATEDEPT_NAME"` // 交易机构
	OPERATEDEPT_CODE string  `json:"OPERATEDEPT_CODE"` // 机构编号
	TOTAL_BUYRIO     float64 `json:"TOTAL_BUYRIO"`     // 买入占比
	TOTAL_SELLRIO    float64 `json:"TOTAL_SELLRIO"`    // 卖出占比
}

func fetchLongHubBuyDetail(date, code string) ([]*LongHubDetailData, error) {
	var res = new(longHubDetailRes)
	filter := fmt.Sprintf("(TRADE_DATE='%s')(SECURITY_CODE=\"%s\")", date, code)
	query := "?reportName=RPT_BILLBOARD_DAILYDETAILSBUY&columns=ALL&pageSize=50&pageNumber=1"
	resp, err := http.Get(longHubApi + query + "&filter=" + url.QueryEscape(filter))
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}

func fetchLongHubSellDetail(date, code string) ([]*LongHubDetailData, error) {
	var res = new(longHubDetailRes)
	filter := fmt.Sprintf("(TRADE_DATE='%s')(SECURITY_CODE=\"%s\")", date, code)
	query := "?reportName=RPT_BILLBOARD_DAILYDETAILSSELL&columns=ALL&pageSize=50&pageNumber=1"
	resp, err := http.Get(longHubApi + query + "&filter=" + url.QueryEscape(filter))
	return res.Result.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
