package log

import (
	"context"
	"filter-core/config"
	"filter-core/util/xcontext"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
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
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder = zapcore.NewConsoleEncoder(encodeConfig)

	level := new(zapcore.Level)
	err := level.UnmarshalText([]byte(conf.Level))
	if err != nil {
		panic(err)
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	z := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(z)
	zap.S().Info(logo)

	return
}

func Infoc(ctx context.Context, template string, args ...interface{}) {
	zap.S().Infof(fmt.Sprintf("trace_id:%v\t%v", xcontext.GetTraceId(ctx), template), args...)
}

func Warnc(ctx context.Context, template string, args ...interface{}) {
	zap.S().Warnf(fmt.Sprintf("trace_id:%v\t%v", xcontext.GetTraceId(ctx), template), args...)
}

func Errorc(ctx context.Context, template string, args ...interface{}) {
	zap.S().Errorf(fmt.Sprintf("trace_id:%v\t%v", xcontext.GetTraceId(ctx), template), args...)
}
