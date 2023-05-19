package util

import (
	"time"
)

func ConvertAgeIntoEpochTime(age int64) int64 {
	now := time.Now()
	then := now.AddDate(int(-age), 0, 0)
	return then.Unix()
}
