package gin

import (
	"context"
	"github.com/PittYao/stream_burn/components/config"
	"github.com/PittYao/stream_burn/components/log"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func NewGinEngine(middlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(config.C.Server.Mode)
	engine := gin.New()
	engine.RemoveExtraSlash = true

	for _, middleware := range middlewares {
		engine.Use(middleware)
	}

	return engine
}

func Run(app http.Handler) {
	addr := ":" + config.C.Server.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	// 启动 http server
	go func() {
		var ln net.Listener
		var err error
		if strings.ToLower(strings.Split(addr, ":")[0]) == "unix" {
			ln, err = net.Listen("unix", strings.Split(addr, ":")[1])
			if err != nil {
				panic(err)
			}
		} else {
			ln, err = net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
		}
		if err := srv.Serve(ln); err != nil {
			log.L.Sugar().Errorf("Server runing error: %s", err.Error())
		}
	}()
	log.L.Sugar().Infof("Server is running on %s", srv.Addr)

	// 监听中断信号， WriteTimeout 时间后优雅关闭服务
	// syscall.SIGTERM 不带参数的 kill 命令
	// syscall.SIGINT ctrl-c kill -2
	// syscall.SIGKILL 是 kill -9 无法捕获这个信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.L.Sugar().Info("Server is shutting down.")

	// 创建一个 context 用于通知 server 3 秒后结束当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.L.Sugar().Error("Server shutdown with error: " + err.Error())
	}
	log.L.Sugar().Info("Server exit.")
}

func Load() error {
	app := NewGinEngine(DefaultGinMiddlewares()...)
	// 注册路由
	RegisterRouter(app)
	// 启动http服务
	Run(app)
	log.L.Info("http服务启动")
	return nil
}
