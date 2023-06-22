package internals

import "os"

var IsDebug bool

func init() {
	IsDebug = os.Getenv("LUNA_ENV") == "debug"
}
