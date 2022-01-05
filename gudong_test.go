package eastmoney

import (
	"fmt"
	"testing"
)

func TestGuDong(t *testing.T) {
	data2, err2 := GuDong("300547", "all")
	fmt.Println(err2, len(data2))
}
