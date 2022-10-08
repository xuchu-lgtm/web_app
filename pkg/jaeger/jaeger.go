package jaeger

import (
	"fmt"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"web_app/settings"
)

var closer io.Closer

// Init 将jaeger tracer设置为全局tracer
func Init(serviceName string, conf *settings.JaegerConfig) (err error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		// 将采样频率设置为1，每一个span都记录
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true, //是否打印日志
			LocalAgentHostPort: fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		},
		ServiceName: serviceName,
	}

	closer, err = cfg.InitGlobalTracer(serviceName, jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		return
	}
	return
}

func Close() {
	_ = closer.Close()
}
