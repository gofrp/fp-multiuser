package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/gofrp/fp-multiuser/pkg/server/controller"
	"golang.org/x/text/language"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cfg  controller.HandleController
	s    *http.Server
	done chan struct{}
}

func New(cfg controller.HandleController) (*Server, error) {
	s := &Server{
		cfg:  cfg,
		done: make(chan struct{}),
	}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Run() error {
	bindAddress := s.cfg.CommonInfo.PluginAddr + ":" + strconv.Itoa(s.cfg.CommonInfo.PluginPort)
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return err
	}
	log.Printf("HTTP server listen on %s", l.Addr().String())
	go func() {
		if err = s.s.Serve(l); !errors.Is(http.ErrServerClosed, err) {
			log.Printf("error shutdown HTTP server: %v", err)
		}
	}()
	<-s.done
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown HTTP server error: %v", err)
	}
	log.Printf("HTTP server exited")
	close(s.done)
	return nil
}

func (s *Server) init() error {
	if err := s.initHTTPServer(); err != nil {
		log.Printf("init HTTP server error: %v", err)
		return err
	}
	return nil
}

func GinI18nLocalize() gin.HandlerFunc {
	acceptLanguage := []language.Tag{language.Chinese, language.English}
	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "./assets/lang",
			AcceptLanguage:   acceptLanguage,
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "json",
			UnmarshalFunc:    json.Unmarshal,
		}),
		ginI18n.WithGetLngHandle(
			func(context *gin.Context, defaultLng string) string {
				header := context.GetHeader("Accept-Language")
				lang := strings.Split(strings.Split(header, ";")[0], ",")
				for _, tag := range acceptLanguage {
					if len(lang) == 1 && tag.String() == lang[0] {
						return lang[0]
					} else if len(lang) == 2 && tag.String() == lang[1] {
						return lang[1]
					}
				}
				return defaultLng
			},
		),
	)
}

func (s *Server) initHTTPServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(GinI18nLocalize())
	s.s = &http.Server{
		Handler: engine,
	}
	controller.NewHandleController(&s.cfg).Register(engine)
	return nil
}
