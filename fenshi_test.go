package eastmoney

import (
	"fmt"
	"testing"
)

func TestFenShi(t *testing.T) {
	datas, pre, err := FenShi("300732", 0)
	fmt.Println(err, len(datas), pre, datas[0])
}
