package eastmoney

import (
	"fmt"
	"testing"
	"time"
)

func TestPanKou(t *testing.T) {
	PanKou(time.Second*3, 10, testHandler, testLogger)
}

func testHandler(ds []*PanKouData) error {
	fmt.Println(len(ds), ds[0])
	return nil
}

func testLogger(err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(time.Now().String())
	}
}
