package logger

import (
	"os"
	"time"

	"filter-core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logo = `
      __        ____     __    _____  ___-  _______     ______  ___________- ___      ___  _______  
     /""\      /  " \   |" \  (\"   \|"  \ |   _  "\   /    " \("     _   ")|"  \    /"  |/" __   ) 
    /    \    /__|| |   ||  | |.\\   \    |(. |_)  :) // ____  \)__/  \\__/  \   \  //  /(__/ _) ./ 
   /' /\  \      |: |   |:  | |: \.   \\  ||:     \/ /  /    ) :)  \\_ /      \\  \/. ./     /  //  
  //  __'  \    _\  |   |.  | |.  \    \. |(|  _  \\(: (____/ //   |.  |       \.    //   __ \_ \\  
 /   /  \\  \  /" \_|\  /\  |\|    \    \ ||: |_)  :)\        /    \:  |        \\   /   (: \__) :\ 
(___/    \___)(_______)(__\_|_)\___|\____\)(_______/  \"_____/      \__|         \__/     \_______) 
                                                                                                    
`
)

func init() {
	conf := config.Conf.LogConf

	var (
		encoder     zapcore.Encoder
		writeSyncer zapcore.WriteSyncer
	)
	//lumberJackLogger := &lumberjack.Logger{
	//	Filename:   conf.FileName,
	//	MaxSize:    conf.MaxSize,
	//	MaxAge:     conf.MaxAge,
	//	MaxBackups: conf.MaxBackups,
	//	Compress:   false,
	//}
	//if conf.IsStdout {
	//	writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	//} else {
	//	writeSyncer = zapcore.AddSync(lumberJackLogger)
	//}
	writeSyncer = zapcore.AddSync(os.Stdout)
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoder = zapcore.NewConsoleEncoder(encodeConfig)

	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	var logger *zap.Logger
	if conf.IsStackTrace {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddCaller())
	}

	zap.ReplaceGlobals(logger)
	return
}

func Logo() {
	zap.S().Info(logo)
}
