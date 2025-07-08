package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	TaskFilePath string
	LogLevel     LogLevel
}

type LogLevel string

const (
	LogDebug LogLevel = "DEBUG"
	LogInfo  LogLevel = "INFO"
	LogWarn  LogLevel = "WARN"
	LogError LogLevel = "ERROR"

	DefaultTaskFilePath = "tasks.json"
	DefaultLogLevel     = LogInfo
)

var (
	ErrInvalidLogLevel = errors.New("invalid log level")
)

func (l LogLevel) String() string {
	return string(l)
}

func (l LogLevel) IsValid() bool {
	switch l {
	case LogDebug, LogInfo, LogWarn, LogError:
		return true
	default:
		return false
	}
}

func GetTaskFilePath(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}

	envFilePath := os.Getenv("TASK_CLI_PATH")
	if envFilePath != "" {
		return envFilePath
	}

	return DefaultTaskFilePath
}

func GetLogLevel(flagValue string) (LogLevel, error) {
	if flagValue != "" {
		lvl := LogLevel(strings.ToUpper(flagValue))
		if !lvl.IsValid() {
			return "", ErrInvalidLogLevel
		}
		return lvl, nil
	}
	envVal := os.Getenv("TASK_CLI_LOG")
	if envVal != "" {
		lvl := LogLevel(strings.ToUpper(envVal))
		if !lvl.IsValid() {
			return "", ErrInvalidLogLevel
		}
		return lvl, nil
	}
	return DefaultLogLevel, nil
}

func LoadConfig(taskFileFlag, logLevelFlag string) (*Config, error) {
	taskFilePath := GetTaskFilePath(taskFileFlag)

	logLevel, err := GetLogLevel(logLevelFlag)
	if err != nil {
		return nil, err
	}

	return &Config{
		TaskFilePath: taskFilePath,
		LogLevel:     logLevel,
	}, nil
}
