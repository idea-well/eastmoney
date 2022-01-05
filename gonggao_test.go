package eastmoney

import (
	"fmt"
	"testing"
)

func TestGongGao(t *testing.T) {
	data, err := AllGonGao("2022-01-01", "2022-01-01", "1")
	fmt.Println(err, len(data), data[0])
}
