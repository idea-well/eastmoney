package eastmoney

import (
	"fmt"
	"testing"
)

func TestFenShi(t *testing.T) {
	datas, err := FenShi("300541", 0)
	fmt.Println(err, len(datas))
	fmt.Printf("%#v",datas.KLineData())
}
