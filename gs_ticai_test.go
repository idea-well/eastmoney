package eastmoney

import (
	"fmt"
	"testing"
)

func TestGsTiCai(t *testing.T) {
	datas, err := GsTiCai()
	fmt.Println(err, len(datas), datas[0])
}
