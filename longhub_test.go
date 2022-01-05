package eastmoney

import (
	"fmt"
	"testing"
)

func TestDateLongHub(t *testing.T) {
	datas, err := DateLongHub("2021-12-30")
	fmt.Println(err, len(datas), fmt.Sprintf("%#v", datas[0]))
}
