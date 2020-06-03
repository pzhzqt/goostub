package common

import (
    "os"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"
)

var Logger log.Logger

func InitLogger(opt level.Option) {
    Logger = log.NewLogfmtLogger(os.Stderr)
    Logger = level.NewFilter(Logger, opt)
    Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
}
