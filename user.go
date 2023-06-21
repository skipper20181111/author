package main

import (
	"author/internal/config"
	"author/internal/handler"
	"author/internal/svc"
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpc"
	"net/http"
	"time"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	go refresscache()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
func refresscache() {
	for true {
		fmt.Println("开始刷新")
		time.Sleep(time.Second * 3)
		urlPath := "http://localhost:8888/refresh/refreshat"
		data := ForceRefresh{Force: false}
		resp, _ := httpc.Do(context.Background(), http.MethodPost, urlPath, data)
		if resp == nil {
			time.Sleep(time.Second * 50)
			continue
		} else {
			fmt.Println("结束刷新", resp)
			fmt.Println(resp.Body.Close())
			time.Sleep(time.Second * 50)
		}
	}
}

type ForceRefresh struct {
	Force bool `json:"force"`
}
