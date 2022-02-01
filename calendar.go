package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type calendarRes struct {
	Data CalendarDatas `json:"data"`
}

type CalendarData struct {
	Date   string `json:"jyrq"` // 交易日期
	Status string `json:"jybz"` // 是否交易
}

func (cd *CalendarData) formatDate(layout string) string {
	if layout == "" || layout == "2006-01-02" {
		return cd.Date
	}
	t, _ := time.Parse("2006-01-02", cd.Date)
	return t.Format(layout)
}

func (cd *CalendarData) formatStatus() bool { return cd.Status == "1" }

// CalendarDatas
type CalendarDatas []*CalendarData

func (ds CalendarDatas) Format(layout string) map[string]bool {
	map_ := make(map[string]bool)
	for _, data := range ds {
		map_[data.formatDate(layout)] = data.formatStatus()
	}
	return map_
}

const calendarApi = "https://www.szse.cn/api/report/exchange/onepersistenthour/monthList?month=%d-%d"

func Calendar(year int) (CalendarDatas, error) {
	datass := make(CalendarDatas, 0)
	for i := 1; i <= 12; i++ {
		datas, err := calendarMonth(year, i)
		if err != nil || len(datas) == 0 {
			return datass, err
		}
		datass = append(datass, datas...)
	}
	return datass, nil
}

func calendarMonth(year, month int) (CalendarDatas, error) {
	var res = new(calendarRes)
	resp, err := http.Get(fmt.Sprintf(calendarApi, year, month))
	return res.Data, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
