package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/WuKongIM/WuKongIM/pkg/wkhttp"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/version"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DemoServer struct {
	r    *wkhttp.WKHttp
	addr string
	s    *Server
	wklog.Log
}

// NewDemoServer new一个demo server
func NewDemoServer(s *Server) *DemoServer {
	// r := wkhttp.New()
	log := wklog.NewWKLog("DemoServer")
	r := wkhttp.NewWithLogger(wkhttp.LoggerWithWklog(log))
	r.Use(wkhttp.CORSMiddleware())

	ds := &DemoServer{
		r:    r,
		addr: s.opts.Demo.Addr,
		s:    s,
		Log:  log,
	}
	return ds
}

// Start 开始
func (s *DemoServer) Start() {

	s.r.GetGinRoute().Use(gzip.Gzip(gzip.DefaultCompression))

	st, _ := fs.Sub(version.DemoFs, "demo/chatdemo/dist")
	demoFS := http.FS(st)
	redirectToDemo := func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/chatdemo/")
		c.Abort()
	}
	serveDemo := func(c *gin.Context) {
		filePath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filePath == "" {
			indexFile, err := demoFS.Open("index.html")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer indexFile.Close()
			stat, err := indexFile.Stat()
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), indexFile)
			return
		}
		c.FileFromFS("/"+filePath, demoFS)
	}

	s.r.GetGinRoute().GET("/", redirectToDemo)
	s.r.GetGinRoute().GET("/chatdemo", redirectToDemo)
	s.r.GetGinRoute().GET("/chatdemo/*filepath", serveDemo)
	s.r.GetGinRoute().HEAD("/chatdemo/*filepath", serveDemo)

	s.setRoutes()
	go func() {
		err := s.r.Run(s.addr) // listen and serve
		if err != nil {
			panic(err)
		}
	}()
	s.Info("Demo server started", zap.String("addr", s.addr))

	_, port := parseAddr(s.addr)
	s.Info(fmt.Sprintf("Chat demo address： http://localhost:%d/chatdemo", port))
}

// Stop 停止服务
func (s *DemoServer) Stop() {
}

func (s *DemoServer) setRoutes() {

}
