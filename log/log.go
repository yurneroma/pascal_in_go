package log

import (
	"fmt"
)

func Log(env string, args ...interface{}) {
	if env == "debug" {
		fmt.Println(args)
	}
}
