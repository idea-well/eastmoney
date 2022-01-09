package eastmoney

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type yanBaoRes struct {
	Hits int         `json:"hits"`
	Size int         `json:"size"`
	Data YanBaoDatas `json:"data"`
	Year int         `json:"currentYear"`
}

func (r *yanBaoRes) setYear() {
	for _, data := range r.Data {
		data.CurrentYear = r.Year
	}
}

type YanBaoData struct {
	Title        string `json:"title"`        // 研报标题
	Content      string `json:"content"`      // 研报内容
	InfoCode     string `json:"infoCode"`     // 研报编号
	OrgName      string `json:"orgName"`      // 机构名称
	OrgSName     string `json:"orgSName"`     // 机构简称
	StockCode    string `json:"stockCode"`    // 股票代码
	StockName    string `json:"stockName"`    // 股票名称
	IndvInduCode string `json:"indvInduCode"` // 个股行业编号
	IndvInduName string `json:"indvInduName"` // 个股行业名称
	IndustryCode string `json:"industryCode"` // 行业行业编号
	IndustryName string `json:"industryName"` // 行业行业名称
	PublishDate  string `json:"publishDate"`  // 发布日期

	CurrentYear           int
	PredictThisYearEps    string `json:"predictThisYearEps"`    // 今年Eps
	PredictThisYearPe     string `json:"predictThisYearPe"`     // 今年Pe
	PredictNextYearEps    string `json:"predictNextYearEps"`    // 明年Eps
	PredictNextYearPe     string `json:"predictNextYearPe"`     // 明年Pe
	PredictNextTwoYearEps string `json:"predictNextTwoYearEps"` // 后年Eps
	PredictNextTwoYearPe  string `json:"predictNextTwoYearPe"`  // 后年Pe
}

func (d *YanBaoData) Predict() map[int]map[string]string {
	var map_ = make(map[int]map[string]string)
	map_[d.CurrentYear] = map[string]string{
		"pe": d.PredictThisYearPe, "eps": d.PredictThisYearEps,
	}
	map_[d.CurrentYear+1] = map[string]string{
		"pe": d.PredictNextYearPe, "eps": d.PredictNextYearEps,
	}
	map_[d.CurrentYear+2] = map[string]string{
		"pe": d.PredictNextTwoYearPe, "eps": d.PredictNextTwoYearEps,
	}
	return map_
}

func (d *YanBaoData) PubTimeFormat(layout string) string {
	t, _ := time.Parse("2006-01-02 15:04:05.999", d.PublishDate)
	return t.Format(layout)
}

func (d *YanBaoData) PdfUrl() string {
	return fmt.Sprintf("https://pdf.dfcfw.com/pdf/H3_%s_1.pdf", d.InfoCode)
}

func (d *YanBaoData) SrcUrl() string {
	return fmt.Sprintf("https://data.eastmoney.com/report/info/%s.html", d.InfoCode)
}

type YanBaoDatas []*YanBaoData

func (ds YanBaoDatas) indexOf(id string) int {
	for i, d := range ds {
		if d.InfoCode == id {
			return i
		}
	}
	return -1
}

func (ds YanBaoDatas) fetchContent() error {
	errs := make(Errors, 0)
	lock := make(chan struct{}, 5)
	spider := newSpider(true)
	spider.OnHTML("#ContentBody .newsContent", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ds[i].Content, _ = e.DOM.Html()
	})
	spider.OnError(func(resp *colly.Response, err error) {
		errs.add(fmt.Errorf("fetch content error, status: %d, error: %v", resp.StatusCode, err))
	})
	spider.OnResponse(func(_ *colly.Response) { <-lock })
	for i, d := range ds {
		lock <- struct{}{}
		frame := fmt.Sprintf("#%d", i)
		_ = spider.Visit(d.SrcUrl() + frame)
	}
	spider.Wait() // wait done
	return errs.first()
}

// GeGuYanBao 个股研报
func GeGuYanBao(begin, end, lastId string) (YanBaoDatas, error) {
	ds, err := doFetchYanBao(0, begin, end, lastId)
	return ds, callWithoutErr(err, ds.fetchContent)
}

// HangYeYanBao 行业研报
func HangYeYanBao(begin, end, lastId string) (YanBaoDatas, error) {
	ds, err := doFetchYanBao(1, begin, end, lastId)
	return ds, callWithoutErr(err, ds.fetchContent)
}

const yanBaoApi = "https://reportapi.eastmoney.com/report/list"

func doFetchYanBao(type_ int, begin, end, lastId string) (YanBaoDatas, error) {
	var query = fmt.Sprintf("?qType=%d&beginTime=%s&endTime=%s", type_, begin, end)
	var datass, page = make(YanBaoDatas, 0), 1
	for {
		datas, err := doFetchYanBaoPage(query, page)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		if i := datas.indexOf(lastId); i == -1 {
			datass = append(datass, datas...)
		} else {
			if i > 0 {
				datass = append(datass, datas[0:i]...)
			}
			return datass, nil
		}
		page++ // next page
	}
}

func doFetchYanBaoPage(query string, pageNo int) (YanBaoDatas, error) {
	var res = new(yanBaoRes)
	var page = fmt.Sprintf("&pageSize=50&pageNo=%d", pageNo)
	resp, err := http.Get(yanBaoApi + query + page)
	return res.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			defer res.setYear()
			return json.Unmarshal(bts, res)
		})
	})
}
