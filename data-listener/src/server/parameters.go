package server

type ParametersServerConfig struct {
	Enabled                                    bool
	NameIsActive                               bool
	Name                                       string
	ConcurrencyIsActive                        bool
	Concurrency                                int
	ReadBufferSizeKilobyteIsActive             bool
	ReadBufferSizeKilobyte                     int
	WriteBufferSizeKilobyteIsActive            bool
	WriteBufferSizeKilobyte                    int
	WriteTimeoutSecondsIsActive                bool
	WriteTimeoutSeconds                        int
	IdleTimeoutSecondsIsActive                 bool
	IdleTimeoutSeconds                         int
	MaxConnsPerIPIsActive                      bool
	MaxConnsPerIP                              int
	MaxRequestsPerConnIsActive                 bool
	MaxRequestsPerConn                         int
	MaxKeepAliveDurationSecondsIsActive        bool
	MaxKeepAliveDurationSeconds                int
	MaxRequestBodySizeKilobyteIsActive         bool
	MaxRequestBodySizeKilobyte                 int
	DisableKeepAliveIsActive                   bool
	DisableKeepAlive                           bool
	TCPKeepAliveIsActive                       bool
	TCPKeepAlive                               bool
	ReduceMemoryUsageIsActive                  bool
	ReduceMemoryUsage                          bool
	GetOnlyIsActive                            bool
	GetOnly                                    bool
	DisablePreParseMultipartFormIsActive       bool
	DisablePreParseMultipartForm               bool
	LogAllErrorsIsActive                       bool
	LogAllErrors                               bool
	SecureErrorLogMessageIsActive              bool
	SecureErrorLogMessage                      bool
	DisableHeaderNamesNormalizingIsActive      bool
	DisableHeaderNamesNormalizing              bool
	SleepWhenConcurrencyLimitsExceededIsActive bool
	SleepWhenConcurrencyLimitsExceededSeconds  int
	NoDefaultServerHeaderIsActive              bool
	NoDefaultServerHeader                      bool
	NoDefaultDateIsActive                      bool
	NoDefaultDate                              bool
	NoDefaultContentTypeIsActive               bool
	NoDefaultContentType                       bool
	KeepHijackedConnsIsActive                  bool
	KeepHijackedConns                          bool
	CloseOnShutdownIsActive                    bool
	CloseOnShutdown                            bool
	StreamRequestBodyIsActive                  bool
	StreamRequestBody                          bool
}

func NewParametersServerConfig() *ParametersServerConfig {
	//notdefaultValues
	return &ParametersServerConfig{
		Enabled:                       true,
		Name:                          "AdvancedServer",
		Concurrency:                   10,
		ReadBufferSizeKilobyte:        1024,
		WriteBufferSizeKilobyte:       1024,
		WriteTimeoutSeconds:           10,
		IdleTimeoutSeconds:            60,
		MaxConnsPerIP:                 10,
		MaxRequestsPerConn:            1000,
		MaxKeepAliveDurationSeconds:   300,
		MaxRequestBodySizeKilobyte:    1024,
		DisableKeepAlive:              false,
		TCPKeepAlive:                  true,
		ReduceMemoryUsage:             false,
		GetOnly:                       false,
		DisablePreParseMultipartForm:  false,
		LogAllErrors:                  false,
		SecureErrorLogMessage:         true,
		DisableHeaderNamesNormalizing: false,
		SleepWhenConcurrencyLimitsExceededSeconds: 10,
		NoDefaultServerHeader:                     false,
		NoDefaultDate:                             false,
		NoDefaultContentType:                      false,
		KeepHijackedConns:                         false,
		CloseOnShutdown:                           true,
		StreamRequestBody:                         false,
	}
}

func (cfg ParametersServerConfig) Start() string {
	return "KEEPIT"
}

func (cfg ParametersServerConfig) String() string {
	if cfg.Enabled {
		return "Server Parameters enabled"
	}
	return "Server Parameters disabled"
}
