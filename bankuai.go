package eastmoney

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type banKuaiRes struct {
	Hits int            `json:"hits"`
	Size int            `json:"size"`
	Data []*BanKuaiData `json:"data"`
}

type BanKuaiData struct {
	BkCode      string `json:"bkCode"`
	BkName      string `json:"bkName"`
	FubkCode    string `json:"fubkCode"`
	PublishCode string `json:"publishCode"`
	FirstLetter string `json:"firstLetter"`
}

// HangYeBanKuai 行业板块
func HangYeBanKuai() ([]*BanKuaiData, error) {
	return doFetchBanKuai("016")
}

// GaiNianBanKuai 概念板块
func GaiNianBanKuai() ([]*BanKuaiData, error) {
	return doFetchBanKuai("007")
}

// DiYuBanKuai 地域板块
func DiYuBanKuai() ([]*BanKuaiData, error) {
	return doFetchBanKuai("020")
}

const banKuaiApi = "https://reportapi.eastmoney.com/report/bk"

func doFetchBanKuai(code string) ([]*BanKuaiData, error) {
	var res = new(banKuaiRes)
	resp, err := http.Get(banKuaiApi + "?bkCode=" + code)
	return res.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		}, func() error {
			return assertError(res.Hits == len(res.Data), "miss hits")
		})
	})
}