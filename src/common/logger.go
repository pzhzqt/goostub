// Copyright (c) 2021 Qitian Zeng
// 
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package common

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

var Logger log.Logger

func InitLogger(opt level.Option) {
	Logger = log.NewLogfmtLogger(os.Stderr)
	Logger = level.NewFilter(Logger, opt)
	Logger = log.With(Logger, "ts", log.DefaultTimestampUTC)
}
