package eastmoney

import (
	"fmt"
	"testing"
)

func TestFenShi(t *testing.T) {
	datas, err := FenShi("300732", 0)
	fmt.Println(err, len(datas))
	kld := datas.KLineData()
	fmt.Println(kld.Pre(), kld.BuyPre(), kld.SellPre())
}
