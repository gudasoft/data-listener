package server

import (
	"crypto/tls"
	"crypto/x509"
	"datalistener/src/handlers"
	"datalistener/src/logging"
	"datalistener/src/metrics"
	"datalistener/src/models"
	"datalistener/src/validation"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type ServerConfig interface {
	Start() string
}

func StartServer(serverConfig []ServerConfig, streamerChannel chan models.EntryData, buffererChannel chan models.EntryData, logger *zap.Logger, validationConfigs []validation.ValidationConfig) *fasthttp.Server {
	jsonValidationConfig, whitelistValidationConfig := validation.GetValidationConfig(validationConfigs)
	baseHandler := func(ctx *fasthttp.RequestCtx) {
		handlers.RequestHandler(ctx, buffererChannel, streamerChannel, logger, &jsonValidationConfig.Enabled)
	}
	var requestHandler fasthttp.RequestHandler
	if whitelistValidationConfig.Enabled {
		requestHandler = whitelistMiddleware(whitelistValidationConfig)(baseHandler)
	} else {

		requestHandler = baseHandler
	}

	serverParameters := NewParametersServerConfig()
	for _, config := range serverConfig {
		if p, ok := config.(*ParametersServerConfig); ok {
			serverParameters = p
			break
		}
	}

	fastHttpServer := getServerWithParameters(*serverParameters, requestHandler)

	for _, config := range serverConfig {
		// factory/builder patterns

		switch config := config.(type) {
		case *HttpServerConfig:
			startHttpServer(fastHttpServer, config)

		case *HttpsServerConfig:
			if config.UseMTLS {
				startHttpsMtlsServer(fastHttpServer, config)
			} else {
				startHttpsServer(fastHttpServer, config)
			}

		case *UnixServerConfig:
			startUnixServer(fastHttpServer, config)

		case *PrometheusServerConfig:
			go metrics.RunMetricsServer(config.Address, config.Port, config.Path)

		case *ParametersServerConfig:
			// do nothing
		default:
			logging.Logger.Sugar().Fatalf("Unsupported protocol in server config: %s\n", reflect.TypeOf(config))
		}
	}
	return fastHttpServer
}

func whitelistMiddleware(whitelistConfig validation.WhitelistConfig) func(fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			clientIP := ctx.RemoteIP().String()

			allowed := false
			for _, ip := range whitelistConfig.Networks {
				_, network, _ := net.ParseCIDR(ip)
				if network.Contains(net.ParseIP(clientIP)) {
					allowed = true
					break
				}
			}

			if !allowed {
				ctx.SetStatusCode(fasthttp.StatusForbidden)
				ctx.SetBodyString("Access denied")
				return
			}

			next(ctx)
		}
	}
}

func startHttpServer(fastHttpServer *fasthttp.Server, cfg *HttpServerConfig) {
	tcpAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	logging.Logger.Sugar().Infof("HTTP Server is listening on %s\n", tcpAddr)
	go func(addr string) {
		if err := fastHttpServer.ListenAndServe(addr); err != nil {
			logging.Logger.Sugar().Fatalf("Error in ListenAndServe: %s", err)
		}
	}(tcpAddr)
}

func startHttpsServer(fastHttpServer *fasthttp.Server, cfg *HttpsServerConfig) {
	tcpAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	certFile := cfg.ServerCertFile
	keyFile := cfg.ServerKeyFile
	logging.Logger.Sugar().Infof("HTTPS Server is listening on %s\n", tcpAddr)
	go func(addr, cert, key string) {
		if err := fastHttpServer.ListenAndServeTLS(addr, cert, key); err != nil {
			logging.Logger.Sugar().Fatalf("Error in ListenAndServeTLS: %s", err)
		}
	}(tcpAddr, certFile, keyFile)
}

func startHttpsMtlsServer(fastHttpServer *fasthttp.Server, cfg *HttpsServerConfig) {
	tcpAddr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	serverCert, err := tls.LoadX509KeyPair(cfg.ServerCertFile, cfg.ServerKeyFile)
	if err != nil {
		logging.Logger.Sugar().Panicf("Error loading server certificate and key for MTLS server: %s\n", err)
		return
	}

	caCert, err := os.ReadFile(cfg.CACertFile)
	if err != nil {
		logging.Logger.Sugar().Panicf("Error loading CA certificate for MTLS server: %s\n", err)
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	fastHttpServer.TLSConfig = tlsConfig

	go func(addr string) {
		if err := fastHttpServer.ListenAndServeTLS(tcpAddr, "", ""); err != nil {
			logging.Logger.Sugar().Fatalf("Error in ListenAndServeTLS: %s", err)
		}
	}(tcpAddr)
}

func startUnixServer(fastHttpServer *fasthttp.Server, cfg *UnixServerConfig) {
	if _, err := os.Stat(cfg.Address); os.IsNotExist(err) {
		os.Mkdir(filepath.Dir(cfg.Address), 0755)
	}
	logging.Logger.Sugar().Infof("Unix socket server is listening on %s\n", cfg.Address)
	go func(addr string) {
		if err := fastHttpServer.ListenAndServeUNIX(addr, 0666); err != nil {
			logging.Logger.Sugar().Fatalf("Error in ListenAndServeUNIX: %s", err)
		}
	}(cfg.Address)
}

func getServerWithParameters(params ParametersServerConfig, requestHandler fasthttp.RequestHandler) *fasthttp.Server {
	server := &fasthttp.Server{
		Handler: requestHandler,
	}

	if params.Enabled {
		if params.NameIsActive {
			server.Name = params.Name
		}

		if params.ConcurrencyIsActive {
			server.Concurrency = params.Concurrency
		}

		if params.ReadBufferSizeKilobyteIsActive {
			server.ReadBufferSize = params.ReadBufferSizeKilobyte * 1024 // Convert Kilobytes to bytes
		}

		if params.WriteBufferSizeKilobyteIsActive {
			server.WriteBufferSize = params.WriteBufferSizeKilobyte * 1024 // Convert Kilobytes to bytes
		}

		if params.WriteTimeoutSecondsIsActive {
			server.WriteTimeout = time.Duration(params.WriteTimeoutSeconds) * time.Second
		}

		if params.IdleTimeoutSecondsIsActive {
			server.IdleTimeout = time.Duration(params.IdleTimeoutSeconds) * time.Second
		}

		if params.MaxConnsPerIPIsActive {
			server.MaxConnsPerIP = params.MaxConnsPerIP
		}

		if params.MaxRequestsPerConnIsActive {
			server.MaxRequestsPerConn = params.MaxRequestsPerConn
		}

		if params.MaxKeepAliveDurationSecondsIsActive {
			server.MaxKeepaliveDuration = time.Duration(params.MaxKeepAliveDurationSeconds) * time.Second
		}

		if params.MaxRequestBodySizeKilobyteIsActive {
			server.MaxRequestBodySize = params.MaxRequestBodySizeKilobyte * 1024 // Convert Kilobytes to bytes
		}

		if params.DisableKeepAliveIsActive {
			server.DisableKeepalive = params.DisableKeepAlive
		}

		if params.TCPKeepAliveIsActive {
			server.TCPKeepalive = params.TCPKeepAlive
		}

		if params.ReduceMemoryUsageIsActive {
			server.ReduceMemoryUsage = params.ReduceMemoryUsage
		}

		if params.GetOnlyIsActive {
			server.GetOnly = params.GetOnly
		}

		if params.DisablePreParseMultipartFormIsActive {
			server.DisablePreParseMultipartForm = params.DisablePreParseMultipartForm
		}

		if params.LogAllErrorsIsActive {
			server.LogAllErrors = params.LogAllErrors
		}

		if params.SecureErrorLogMessageIsActive {
			server.SecureErrorLogMessage = params.SecureErrorLogMessage
		}

		if params.DisableHeaderNamesNormalizingIsActive {
			server.DisableHeaderNamesNormalizing = params.DisableHeaderNamesNormalizing
		}

		if params.SleepWhenConcurrencyLimitsExceededIsActive {
			server.SleepWhenConcurrencyLimitsExceeded = time.Duration(params.SleepWhenConcurrencyLimitsExceededSeconds) * time.Second
		}

		if params.NoDefaultServerHeaderIsActive {
			server.NoDefaultServerHeader = params.NoDefaultServerHeader
		}

		if params.NoDefaultDateIsActive {
			server.NoDefaultDate = params.NoDefaultDate
		}

		if params.NoDefaultContentTypeIsActive {
			server.NoDefaultContentType = params.NoDefaultContentType
		}

		if params.KeepHijackedConnsIsActive {
			server.KeepHijackedConns = params.KeepHijackedConns
		}

		if params.CloseOnShutdownIsActive {
			server.CloseOnShutdown = params.CloseOnShutdown
		}

		if params.StreamRequestBodyIsActive {
			server.StreamRequestBody = params.StreamRequestBody
		}
	}

	return server
}
