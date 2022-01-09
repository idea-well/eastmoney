package eastmoney

import (
	"fmt"
	"testing"
)

func TestYanBao(t *testing.T) {
	data1, err1 := GeGuYanBao("2021-12-31", "2022-01-02", "AP202201011537872199")
	data2, err2 := HangYeYanBao("2021-12-31", "2022-01-02", "")
	fmt.Println(err1, err2, len(data1), len(data2), data2[0])
}
