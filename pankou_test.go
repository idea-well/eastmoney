package eastmoney

import (
	"fmt"
	"testing"
	"time"
)

func TestPanKou(t *testing.T) {
	PanKou(time.Minute, 100, testHandler, testLogger)
}

func testHandler(ds []*PanKouData) error {
	fmt.Println(len(ds), ds[0])
	return nil
}

func testLogger(err error) {
	fmt.Println(err)
}
