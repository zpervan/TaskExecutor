package util

import (
	"fmt"
	"time"
)

func CurrentDateTime() string {
	now := time.Now()

	return fmt.Sprintf("%d/%d/%d_%d:%d:%d",
		now.Month(),
		now.Day(),
		now.Year(),
		now.Hour(),
		now.Hour(),
		now.Second())
}
