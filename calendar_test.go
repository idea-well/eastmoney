package eastmoney

import (
	"fmt"
	"testing"
)

func TestCalendar(t *testing.T) {
	datas, err := Calendar(2005)
	fmt.Println(err, len(datas), datas.Format("20060102"))
}
