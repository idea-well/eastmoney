package eastmoney

import (
	"fmt"
	"testing"
)

func TestZhiYa(t *testing.T) {
	datas, err := ZhiYa("2022-01-07")
	fmt.Println(err, len(datas), datas[0])
}
