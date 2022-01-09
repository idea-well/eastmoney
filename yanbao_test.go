package eastmoney

import (
	"fmt"
	"testing"
)

func TestYanBao(t *testing.T) {
	data1, err1 := GeGuYanBao("2020-12-31", "2021-01-01", "")
	data2, err2 := HangYeYanBao("2021-12-31", "2022-01-01", "")
	fmt.Println(err1, err2, len(data1), len(data2), data1[0].Predict())
}
