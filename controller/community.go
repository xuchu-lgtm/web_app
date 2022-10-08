package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"strconv"
	"time"
	"web_app/dao/mysql"
	"web_app/logic"
)

func IndexHandler(c *gin.Context) {

	//tracer := opentracing.GlobalTracer()
	//span := tracer.StartSpan("local index")

	//c.HTML(http.StatusOK, "index.html", nil)

	//span.Finish()

	span := opentracing.SpanFromContext(c.Request.Context())

	ctx := opentracing.ContextWithSpan(context.Background(), span)
	Foo(ctx)
}

func Foo(ctx context.Context) {
	// 开始一个span, 设置span的operation_name=Foo
	span, ctx := opentracing.StartSpanFromContext(ctx, "Foo")
	defer span.Finish()
	// 将context传递给Bar
	Bar(ctx)
	// 模拟执行耗时
	time.Sleep(1 * time.Second)
}
func Bar(ctx context.Context) {
	// 开始一个span，设置span的operation_name=Bar
	span, ctx := opentracing.StartSpanFromContext(ctx, "Bar")
	defer span.Finish()
	// 模拟执行耗时
	time.Sleep(2 * time.Second)

	// 假设Bar发生了某些错误
	err := errors.New("something wrong")
	span.LogFields(
		log.String("event", "error"),
		log.String("message", err.Error()),
	)
	span.SetTag("error", true)
}

func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed err,", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed err,", zap.Error(err))
		if errors.Is(err, mysql.ErrorInvalidId) {
			ResponseError(c, CodeInvalidId)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
