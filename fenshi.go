package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type fengShiRes struct {
	Data struct {
		Details  fenShiStrings `json:"details"`
		PrePrice float64       `json:"prePrice"`
	} `json:"data"`
}

// fenShiStrings
type fenShiStrings []string

func (fs fenShiStrings) toData() FengShiDatas {
	ds := make(FengShiDatas, 0)
	for _, str := range fs {
		ss := strings.Split(str, ",")
		tp := ParseInt(ss[4])
		if tp != 1 && tp != 2 {
			continue
		}
		ds = append(ds, &FengShiData{
			Time:   strings.ReplaceAll(ss[0], ":", ""),
			Price:  ParseFloat(ss[1]),
			Volume: ParseInt(ss[2]),
			Count:  ParseInt(ss[3]),
			Type:   ParseInt(ss[4]),
		})
	}
	return ds
}

// FengShiData
type FengShiData struct {
	Time   string  `json:"t"`     // 成交时间
	Type   int     `json:"bs"`    // 1卖 2买
	Price  float64 `json:"p"`     // 成交价格
	Volume int     `json:"v"`     // 成交手数
	Count  int     `json:"count"` // 成交笔数
}

// FengShiDatas
type FengShiDatas []*FengShiData

const fenShiApi = "https://push2.eastmoney.com/api/qt/stock/details/get" +
	"?fields1=f4&fields2=f51,f52,f53,f54,f55&pos=0&secid=%d.%s"

func FenShi(code string, market int) (FengShiDatas, float64, error) {
	var res = new(fengShiRes)
	url := fmt.Sprintf(fenShiApi, market, code)
	resp, err := http.Get(url)
	err = callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
	return res.Data.Details.toData(), res.Data.PrePrice, err
}
