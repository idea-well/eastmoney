package eastmoney

import (
	"fmt"
	"testing"
)

func TestDateLongHub(t *testing.T) {
	datas, err := DateLongHub("2022-01-10")
	fmt.Println(err, len(datas.STDatas()))
}
