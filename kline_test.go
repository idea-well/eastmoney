package eastmoney

import (
	"fmt"
	"testing"
)

func TestKLine(t *testing.T) {
	data, err := KLine("300547", 0)
	fmt.Println(err, len(data), data["20220128"])
	fmt.Println(KLineDate("300547", 0, "20220128"))
}
