package eastmoney

import (
	"fmt"
	"testing"
)

func TestBusiness(t *testing.T) {
	datas, err := Business("sz300547")
	fmt.Println(err, len(datas), datas[0])
}
