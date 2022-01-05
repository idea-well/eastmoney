package eastmoney

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
	"strings"
)

const pengPaiUrl = "https://www.thepaper.cn"

type PengPaiData struct {
	Id      string `json:"id"`       // 新闻ID
	Title   string `json:"title"`    // 新闻标题
	Content string `json:"content"`  // 新闻内容
	PubTime string `json:"pub_time"` // 发布时间
	SrcUrl  string `json:"src_url"`  // 原文链接
}

type PengPaiDatas []*PengPaiData

func (ds PengPaiDatas) indexOf(id string) int {
	for i, d := range ds {
		if d.Id == id {
			return i
		}
	}
	return -1
}

func (ds PengPaiDatas) fetchContent() error {
	var errs = make(Errors, 0)
	spider := newSpider(true)
	spider.OnHTML(".news_about", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ds[i].PubTime = e.ChildTexts("p")[1][0:16]
	})
	spider.OnHTML(".news_txt", func(e *colly.HTMLElement) {
		i, _ := strconv.Atoi(e.Request.URL.Fragment)
		ds[i].Content, _ = e.DOM.Html()
	})
	spider.OnError(func(resp *colly.Response, err error) {
		errs = append(errs, fmt.Errorf("fetch content error, status: %d, error: %v", resp.StatusCode, err))
	})
	for i, d := range ds {
		fmt.Println(d.Title)
		_ = spider.Visit(d.SrcUrl + fmt.Sprintf("#%d", i))
	}
	spider.Wait() // wait done
	return errs.first()
}

func PengPaiCaiJing(lastId string) (PengPaiDatas, error) {
	ds, err := fetchPengPai(lastId)
	return ds, callWithoutErr(err, ds.fetchContent)
}

func fetchPengPai(lastId string) (PengPaiDatas, error) {
	var datass, page = make(PengPaiDatas, 0), 1
	for {
		datas, err := fetchPengPaiPage(page)
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

func fetchPengPaiPage(pageNo int) (PengPaiDatas, error) {
	datas := make(PengPaiDatas, 0)
	spider := newSpider(false)
	spider.OnHTML(".news_li h2 a", func(e *colly.HTMLElement) {
		ss := strings.Split(e.Attr("href"), "_")
		datas = append(datas, &PengPaiData{
			Id:     ss[2],
			Title:  e.Text,
			SrcUrl: pengPaiUrl + "/" + e.Attr("href"),
		})
	})
	query := fmt.Sprintf("?channelID=25951&pageidx=%d", pageNo)
	err := spider.Visit(pengPaiUrl + "/load_index.jsp" + query)
	return datas, err
}
