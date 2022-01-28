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
	Time   int `json:"t"`  // 成交时间
	Type   int `json:"bs"` // 1卖 2买
	Price  int `json:"p"`  // 成交价格
	Volume int `json:"v"`  // 成交手数
}

type FengShiDatas []*FengShiData

func (ds FengShiDatas) realData() FengShiDatas {
	for i, d := range ds {
		if d.Type == 1 || d.Type == 2 {
			return ds[i:]
		}
	}
	return ds
}

func (ds FengShiDatas) KLineData() (kld KLineData) {
	kld.Low = ds[0].Price
	kld.Close = ds[len(ds)-1].Price
	for i, d := range ds {
		if d.Type == 2 {
			kld.Buy.Volume += d.Volume
			kld.Buy.Amount += d.Volume * d.Price
		} else {
			kld.Sell.Volume += d.Volume
			kld.Sell.Amount += d.Volume * d.Price
		}
		if kld.Open == 0 && d.Time >= 93000 {
			kld.Open = ds[i-1].Price
		}
		if kld.Low > d.Price {
			kld.Low = d.Price
		}
		if kld.High < d.Price {
			kld.High = d.Price
		}
	}
	return
}

type KLineData struct {
	Open  int `json:"open"`  // 开盘
	Close int `json:"close"` // 收盘
	High  int `json:"high"`  // 最高
	Low   int `json:"low"`   // 最低
	Buy   struct {
		Volume int `json:"volume"`
		Amount int `json:"amount"`
	} `json:"buy"`
	Sell struct {
		Volume int `json:"volume"`
		Amount int `json:"amount"`
	} `json:"sell"`
}

func (kld *KLineData) Pre() float64 {
	return float64(kld.Buy.Amount+kld.Sell.Amount) / float64(kld.Buy.Volume+kld.Sell.Volume)
}

func (kld *KLineData) BuyPre() float64 {
	return float64(kld.Buy.Amount) / float64(kld.Buy.Volume)
}

func (kld *KLineData) SellPre() float64 {
	return float64(kld.Sell.Amount) / float64(kld.Sell.Volume)
}

func FenShi(code string, market int) (FengShiDatas, error) {
	var datass, page, size = make(FengShiDatas, 0), 0, 1000
	for {
		datas, err := doFetchFenShiPage(code, market, page, size)
		if len(datas) > 0 {
			datass = append(datass, datas...)
		}
		if err != nil || len(datas) < size {
			return datass.realData(), err
		}
		page++ // next page
	}
}

const fenShiApi = "http://push2ex.eastmoney.com/getStockFenShi" +
	"?ut=7eea3edcaed734bea9cbfc24409ed989&dpt=wzfscj&sort=1&ft=1"

func doFetchFenShiPage(code string, market, page, size int) (FengShiDatas, error) {
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
