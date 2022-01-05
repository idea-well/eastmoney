package eastmoney

import (
	"fmt"
	"testing"
)

func TestBanKuai(t *testing.T) {
	data1, err1 := HangYeBanKuai()
	data2, err2 := GaiNianBanKuai()
	data3, err3 := DiYuBanKuai()
	fmt.Println(err1, err2, err3, len(data1), len(data2), len(data3))
}
