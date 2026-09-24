package commavalue

import (
	"strconv"
	"strings"
	"time"
)

func Trim(str string) string {
	return strings.TrimSpace(strings.ToUpper(str))
}

func Int64(str string) (int64, error) {
	return strconv.ParseInt(Trim(str), 10, 64)
}

func Float64(str string) (float64, error) {
	return strconv.ParseFloat(Trim(str), 64)
}

func ToFloat64(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func Date(str string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, Trim(str))
}
