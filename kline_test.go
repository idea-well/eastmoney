package eastmoney

import (
	"fmt"
	"testing"
)

func TestKLine(t *testing.T) {
	data, err := KLine("300547", 0)
	fmt.Println(err, len(data), data.at("20220128"))
	fmt.Println(KLineDate("600112", 1, "20220104"))
}

func TestMLine(t *testing.T) {
	data, err := MLine("BK0539", 90)
	fmt.Println(err, len(data), data.at("0931"))
}
