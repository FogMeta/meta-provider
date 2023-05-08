package pocket

import (
	logging "github.com/ipfs/go-log/v2"
)

var log *logging.ZapEventLogger

func GetLog() *logging.ZapEventLogger {
	if log == nil {
		log = logging.Logger("pocket")
		logging.SetLogLevel("pocket", "INFO")
	}

	return log
}

func SetLogLevel(level string) {
	logging.SetLogLevel("pocket", level)
}
