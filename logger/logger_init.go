package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

type loggerConfig struct {
	encoder      zapcore.Encoder
	writerSyncer zapcore.WriteSyncer
	levelEnable  zapcore.Level
}

func InitLogger() {
	core := zapcore.NewCore(
		getEncoder(),
		zapcore.NewMultiWriteSyncer(
			getWriterSyncer(),
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(os.Stderr),
		),
		getLevelEnable(),
	)

	Logger = zap.New(core)
	SugarLogger = Logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncer() zapcore.WriteSyncer {
	file, _ := os.OpenFile("../output/output.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	return zapcore.AddSync(file)
}

func getLevelEnable() zapcore.Level {
	return zapcore.DebugLevel
}
