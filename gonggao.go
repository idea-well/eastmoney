package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type gongGaoRes struct {
	Data struct {
		List GongGaoDatas `json:"list"`
	} `json:"data"`
}

type GongGaoData struct {
	Title      string `json:"title"`       // 公告标题
	Content    string `json:"content"`     // 公告内容
	ArtCode    string `json:"art_code"`    // 公告编号
	NoticeDate string `json:"notice_date"` // 公告日期
	Codes      []struct {
		StockCode string `json:"stock_code"` // 股票代码
		ShortName string `json:"short_name"` // 股票名称
	} `json:"codes"`
}

func (d *GongGaoData) PdfUrl() string {
	return fmt.Sprintf("https://pdf.dfcfw.com/pdf/H2_%s_1.pdf", d.ArtCode)
}

type GongGaoDatas []*GongGaoData

// AllGonGao 公告查询
// market: SHA, SZA, BJA, CYB, KCB
// fNode: 1.财务报告 2.融资公告 3.风险提示 4.信息变更 5.重大事项 6.资产重组 7.持股变动
func AllGonGao(market, begin, end string, fNode string) (GongGaoDatas, error) {
	var query = fmt.Sprintf(
		"?begin_time=%s&end_time=%s&ann_type=%s&f_node=%s",
		begin, end, market, fNode,
	)
	var datass, page = make(GongGaoDatas, 0), 1
	for {
		datas, err := doFetchGongGaoPage(query, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
		page++ // next page
	}
}

const gongGaoApi = "https://np-anotice-stock.eastmoney.com/api/security/ann"

func doFetchGongGaoPage(query string, pageNo int) (GongGaoDatas, error) {
	var res = new(gongGaoRes)
	var page = fmt.Sprintf("&page_size=50&page_index=%d", pageNo)
	resp, err := http.Get(gongGaoApi + query + page)
	return res.Data.List, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}

// gongGaoContentRes 公告内容
type gongGaoContentRes struct {
	Data struct {
		NoticeContent string `json:"notice_content"`
	} `json:"data"`
}

type gongGaoNoticeContent struct {
}
