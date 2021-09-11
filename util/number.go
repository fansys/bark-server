package util

import "strconv"

func Int64ToInt(i64 int64) int {
	i, _ := strconv.Atoi(strconv.FormatInt(i64, 10))
	return i
}

func IntToInt64(i int) int64 {
	return int64(i)
}
