package eastmoney

import (
	"fmt"
	"testing"
)

func TestPengPaiCaiJing(t *testing.T) {
	datas, err := PengPaiCaiJing("16059604")
	fmt.Println(err, len(datas), datas[0])
}
