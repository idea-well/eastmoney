package eastmoney

import (
	"fmt"
	"testing"
)

func TestFenShi(t *testing.T) {
	datas, err := FenShi("300921", 0)
	fmt.Println(err, len(datas), datas[0])
	kld := datas.KLineData()
	fmt.Println(kld.Pre(), kld.BuyPre(), kld.SellPre())
}
