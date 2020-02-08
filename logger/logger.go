// Copyright 2019 Trần Anh Dũng <chiro@fkguru.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"encoding/hex"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var once sync.Once
var onceSugar sync.Once
var sugar *zap.SugaredLogger

//GetLogger get singleton logger
func GetLogger() *zap.Logger {
	once.Do(func() {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = config.Build()
	})
	defer logger.Sync()
	return logger
}

//GetSugarLogger get singleton sugar logger
func GetSugarLogger() *zap.SugaredLogger {
	logger := GetLogger()
	onceSugar.Do(func() {
		sugar = logger.Sugar()
	})
	return sugar
}

//HexDump for debug purpose
func HexDump(title string, data []byte) {
	sugar := GetSugarLogger()
	content := hex.Dump(data)
	sugar.Debugf("%s\n%s", title, content[:len(content)-1])
}
