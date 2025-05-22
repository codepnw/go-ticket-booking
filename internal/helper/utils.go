package helper

import "time"

var LocalTime = time.Now().Local().Format(time.DateTime)
