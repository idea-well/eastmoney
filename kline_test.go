package eastmoney

import (
	"fmt"
	"testing"
)

func TestKLine(t *testing.T) {
	data, err := KLine("300547", 0)
	fmt.Println(err, len(data), data["2022-01-28"])
}
