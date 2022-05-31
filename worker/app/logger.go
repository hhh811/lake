/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"
	"github.com/apache/incubator-devlake/logger"

	"go.temporal.io/sdk/log"
)

type TemporalLogger struct {
	log logger.Logger
}

func NewTemporalLogger(log logger.Logger) log.Logger {
	return &TemporalLogger{
		log,
	}
}

func (l *TemporalLogger) Log(lv logger.LogLevel, msg string, keyvals ...interface{}) {
	if l.log.IsLevelEnabled(lv) {
		for i := 0; i < len(keyvals); i += 2 {
			msg += fmt.Sprintf(" %s %v", keyvals[i], keyvals[i+1])
		}
		l.log.Log(lv, msg)
	}
}

func (l *TemporalLogger) Debug(msg string, keyvals ...interface{}) {
	l.Log(logger.LOG_DEBUG, msg, keyvals...)
}

func (l *TemporalLogger) Info(msg string, keyvals ...interface{}) {

	l.Log(logger.LOG_INFO, msg, keyvals...)
}

func (l *TemporalLogger) Warn(msg string, keyvals ...interface{}) {

	l.Log(logger.LOG_WARN, msg, keyvals...)
}

func (l *TemporalLogger) Error(msg string, keyvals ...interface{}) {
	l.Log(logger.LOG_ERROR, msg, keyvals...)
}
