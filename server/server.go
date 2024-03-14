package server

import (
	"buffer-handler/handlers"
	"buffer-handler/logging"
	"buffer-handler/metrics"
	"buffer-handler/models"
	"fmt"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func StartServer(serverConfig []ServerConfig, streamerChannel chan models.EntryData, buffererChannel chan models.EntryData) *fasthttp.Server {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		handlers.RequestHandler(ctx, buffererChannel, streamerChannel)
	}

	fastHttpServer := &fasthttp.Server{
		Handler: requestHandler,
	}

	for _, config := range serverConfig {

		switch config.(type) {
		case *HttpServerConfig:
			httpConfig := config.(*HttpServerConfig)
			tcpAddr := fmt.Sprintf("%s:%d", httpConfig.Address, httpConfig.Port)
			logging.Logger.Sugar().Infof("HTTP Server is listening on %s\n", tcpAddr)
			go func(addr string) {
				if err := fastHttpServer.ListenAndServe(addr); err != nil {
					logging.Logger.Sugar().Fatalf("Error in ListenAndServe: %s", err)
				}
			}(tcpAddr)
		case *HttpsServerConfig:
			httpsConfig := config.(*HttpsServerConfig)
			tcpAddr := fmt.Sprintf("%s:%d", httpsConfig.Address, httpsConfig.Port)
			certFile := httpsConfig.TlsCert
			keyFile := httpsConfig.TlsKey
			logging.Logger.Sugar().Infof("HTTPS Server is listening on %s\n", tcpAddr)
			go func(addr, cert, key string) {
				if err := fastHttpServer.ListenAndServeTLS(addr, cert, key); err != nil {
					logging.Logger.Sugar().Fatalf("Error in ListenAndServeTLS: %s", err)
				}
			}(tcpAddr, certFile, keyFile)

		case *UnixServerConfig:
			unixConfig := config.(*UnixServerConfig)
			if _, err := os.Stat(unixConfig.Address); os.IsNotExist(err) {
				os.Mkdir(filepath.Dir(unixConfig.Address), 0755)
			}
			logging.Logger.Sugar().Infof("Unix socket server is listening on %s\n", unixConfig.Address)
			go func(path string) {
				if err := fastHttpServer.ListenAndServeUNIX(path, 0666); err != nil {
					logging.Logger.Sugar().Fatalf("Error in ListenAndServeUNIX: %s", err)
				}
			}(unixConfig.Address)

		case *PrometheusServerConfig:
			promConfig := config.(*PrometheusServerConfig)
			go metrics.RunMetricsServer(promConfig.Address, promConfig.Port, promConfig.Path)
		default:
			fmt.Println("Unsupported protocol in server config")
		}
	}
	return fastHttpServer
}
