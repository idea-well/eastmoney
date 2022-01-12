package eastmoney

import (
	"fmt"
	"testing"
)

func TestDateLongHub(t *testing.T) {
	datas, _ := DateLongHub("2022-01-10", "069001002002")
	for _, g := range datas.Groups() {
		fmt.Println(g.Explanation, len(g.Datas), g.Datas.Codes())
	}
}
