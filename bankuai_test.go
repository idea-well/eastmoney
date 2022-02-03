package eastmoney

import (
	"fmt"
	"testing"
)

func TestBanKuai(t *testing.T) {
	data, err := BanKuai()
	fmt.Println(err, len(data), data["100"])
}
