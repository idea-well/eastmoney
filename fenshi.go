package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

type fengShiRes struct {
	Data struct {
		Data []*FengShiData `json:"data"`
	} `json:"data"`
}

type FengShiData struct {
	Time   int     `json:"t"`  // 成交时间
	Type   int     `json:"bs"` // 1卖 2买
	Price  float64 `json:"p"`  // 成交价格
	Volume float64 `json:"v"`  // 成交手数
}

func (fs FengShiData) CNY() float64 {
	return fs.Price / float64(1000)
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
	kld.Low = ds[0].CNY()
	kld.Open = ds[0].CNY()
	kld.High = ds[0].CNY()
	kld.Close = ds[len(ds)-1].CNY()
	for _, d := range ds {
		if d.Type == 2 {
			kld.Buy.Count += 1
			kld.Buy.Volume += d.Volume
			kld.Buy.Amount += d.Volume * d.CNY() * 100
		} else {
			kld.Sell.Count += 1
			kld.Sell.Volume += d.Volume
			kld.Sell.Amount += d.Volume * d.CNY() * 100
		}
		kld.Volume += d.Volume
		kld.Amount += d.Volume * d.CNY() * 100
		kld.Low = math.Min(kld.Low, d.CNY())
		kld.High = math.Max(kld.High, d.CNY())
	}
	return
}

type KLineData struct {
	Open      float64 `json:"open"`      // 开盘
	Close     float64 `json:"close"`     // 收盘
	High      float64 `json:"high"`      // 最高
	Low       float64 `json:"low"`       // 最低
	Volume    float64 `json:"volume"`    // 成交量
	Amount    float64 `json:"amount"`    // 成交额
	Change    float64 `json:"change"`    // 日涨幅
	Turnover  float64 `json:"turnover"`  // 换手率
	Amplitude float64 `json:"amplitude"` // 日振幅
	Buy       struct {
		Count  float64 `json:"count"`
		Volume float64 `json:"volume"`
		Amount float64 `json:"amount"`
	} `json:"buy"`
	Sell struct {
		Count  float64 `json:"count"`
		Volume float64 `json:"volume"`
		Amount float64 `json:"amount"`
	} `json:"sell"`
}

func (kld *KLineData) AvgPrice() float64 {
	return kld.Amount / (kld.Volume * 100)
}

func (kld *KLineData) BuyAvgPrice() float64 {
	return kld.Buy.Amount / (kld.Buy.Volume * 100)
}

func (kld *KLineData) SellAvgPrice() float64 {
	return kld.Sell.Amount / (kld.Sell.Volume * 100)
}

func (kld *KLineData) BuyAvgCount() float64 {
	return kld.Buy.Volume / kld.Buy.Count
}

func (kld *KLineData) SellAvgCount() float64 {
	return kld.Sell.Volume / kld.Sell.Count
}

func (kld *KLineData) BuyAvgAmount() float64 {
	return kld.Buy.Amount / kld.Buy.Count
}

func (kld *KLineData) SellAvgAmount() float64 {
	return kld.Sell.Amount / kld.Sell.Count
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
