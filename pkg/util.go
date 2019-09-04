package promtail

import "go.uber.org/zap"

var _logger *zap.Logger

func Info(msg string, fields ...zap.Field) {
    if _logger != nil {
        _logger.Info(msg, fields...)
    }
}

func Warnf(msg string, err error) {
    if _logger != nil {
        _logger.Warn(msg, zap.Error(err))
    }
}

func Errof(msg string, err error) {
    if _logger != nil {
        _logger.Error(msg, zap.Error(err))
    }
}

func Error(msg string, fields ...zap.Field) {
    if _logger != nil {
        _logger.Info(msg, fields...)
    }
}

func InitLogger(logger *zap.Logger) {
    _logger = logger
}