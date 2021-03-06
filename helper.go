package eastmoney

import (
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseFloat(s string) float64 {
	fs, _ := decimal.NewFromString(s)
	fv, _ := fs.Round(2).Float64()
	return fv
}

func TodayString() string {
	return time.Now().Format("20060102")
}

func MinString(s string, ss ...string) string {
	if len(ss) == 0 {
		return s
	}
	if s <= ss[0] {
		return MinString(s, ss[1:]...)
	}
	return MinString(ss[0], ss[1:]...)
}

func MaxString(s string, ss ...string) string {
	if len(ss) == 0 {
		return s
	}
	if s >= ss[0] {
		return MaxString(s, ss[1:]...)
	}
	return MaxString(ss[0], ss[1:]...)
}
