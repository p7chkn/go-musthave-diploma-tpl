package utils

import (
	"strconv"
	"time"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CalculateAdditionTime(count int) time.Duration {
	return time.Second * time.Duration(30*(count+1))
}
