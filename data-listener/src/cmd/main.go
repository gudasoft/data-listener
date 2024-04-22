package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"datalistener/src/config"
	"datalistener/src/destinations"
	"datalistener/src/logging"
	"datalistener/src/models"
	"datalistener/src/server"
	"datalistener/src/validation"
	"datalistener/src/version"
	"datalistener/src/workers"
)

var (
	shutdownChannelBuf = make(chan struct{})
	shutdownChannelStr = make(chan struct{})
	reloadChannelBuf   = make(chan bool)
	readyChannel       = make(chan bool)
	configFile         = "config.toml"
)

func main() {
	showVersion, configFileParameter := parseFlags()

	if showVersion {
		fmt.Printf("Data-Listener: v%s\n", version.Current)
		return
	}
	if configFileParameter != "" {
		configFile = configFileParameter
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	var logConfig logging.LogConfig
	var serverConfig []server.ServerConfig
	var streamConfigs []destinations.StreamConfig
	var buffererConfigs []destinations.BufferConfig
	var validationCongfigs []validation.ValidationConfig
	var bufferSize int
	var convertToJSONL bool
	// var validateJSON bool
	config.LoadConfigs(configFile, &logConfig, &serverConfig, &streamConfigs, &buffererConfigs, &bufferSize, &convertToJSONL, &validationCongfigs)
	logger, err := logConfig.InitLogger()
	if err != nil {
		panic(err)
	}

	logConfigs(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)

	var streamerChannel chan models.EntryData
	if len(streamConfigs) > 0 {
		streamerChannel = make(chan models.EntryData)
		go workers.DestinationsNotifier(&streamConfigs, streamerChannel, shutdownChannelStr, readyChannel)
	}

	var buffererChannel chan models.EntryData
	if len(buffererConfigs) > 0 {
		buffererChannel = make(chan models.EntryData)
		go workers.DestinationsNotifierBuffered(&bufferSize, &buffererConfigs, &convertToJSONL, buffererChannel, shutdownChannelBuf, reloadChannelBuf, readyChannel)
	}
	fmt.Println("Starting server...")
	fmt.Println(validationCongfigs)

	fastHttpServer := server.StartServer(serverConfig, streamerChannel, buffererChannel, logger, validationCongfigs)

	printWelcomeBanner(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL, validationCongfigs)

sigloop:
	for {
		select {
		case sigReceived := <-signalChannel:
			if sigReceived == syscall.SIGTERM || sigReceived == syscall.SIGINT {
				shutdownSignalMessage := "Received shutdown signal. Flushing up..."
				fmt.Println(shutdownSignalMessage)
				logging.Logger.Info(shutdownSignalMessage)
				fastHttpServer.Shutdown()
				if len(buffererConfigs) > 0 {
					close(shutdownChannelBuf)
					<-readyChannel
				}
				if len(streamConfigs) > 0 {
					close(shutdownChannelStr)
					<-readyChannel
				}
				break sigloop
			} else if sigReceived == syscall.SIGHUP {
				reloadSignalMessage := "Received configuration reload signal. Loading up..."
				fmt.Println(reloadSignalMessage)
				logging.Logger.Info(reloadSignalMessage)
				if len(buffererConfigs) > 0 {
					reloadChannelBuf <- true
					<-readyChannel
				}
				config.LoadConfigs(configFile, &logConfig, &serverConfig, &streamConfigs, &buffererConfigs, &bufferSize, &convertToJSONL, &validationCongfigs)
				logConfig.InitLogger()
				logConfigs(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)
				printConfigInfo(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL, validationCongfigs)
			}
		case <-ctx.Done():
			close(shutdownChannelBuf)
			logging.Logger.Info("Context canceled. Cleaning up...")
			time.Sleep(100 * time.Millisecond)
			fastHttpServer.Shutdown()
			break sigloop
		}
	}
	fmt.Println("Exit")
}

func parseFlags() (showVersion bool, configFileParameter string) {
	versionFlag := flag.Bool("version", false, "Print the version of the program")
	versionFlagShort := flag.Bool("v", false, "Print the version of the program (short)")
	configFlag := flag.String("config", "", "Path to configuration file")
	configFlagShort := flag.String("c", "", "Path to configuration file (short)")

	flag.Parse()

	showVersion = *versionFlag || *versionFlagShort

	if *configFlag != "" {
		configFileParameter = *configFlag
	} else if *configFlagShort != "" {
		configFileParameter = *configFlagShort
	}
	return
}

func logConfigs(logConfig logging.LogConfig, serverConfigs []server.ServerConfig, streamerConfigs []destinations.StreamConfig, bufferSize int, buffererConfigs []destinations.BufferConfig, convertToJSONL bool) {
	logging.Logger.Debug("Logger Config:")
	logging.Logger.Debug(logConfig.String())
	logging.Logger.Debug("Server Config:")
	logging.Logger.Debug(config.GetServerConfigInfo(serverConfigs))
	logging.Logger.Debug("Streamer Configs:")
	logging.Logger.Debug(config.GetStreamConfigInfo(streamerConfigs))
	logging.Logger.Debug(("Bufferer Configs:"))
	logging.Logger.Sugar().Debugf("Buffer size: %d Bytes", bufferSize)
	logging.Logger.Sugar().Debugf("Convert to jsonl: %t", convertToJSONL)
	logging.Logger.Debug(config.GetBufferConfigInfo(buffererConfigs))
}

func printWelcomeBanner(logConfig logging.LogConfig, serverConfig []server.ServerConfig, streamerConfigs []destinations.StreamConfig, bufferSize int, buffererConfigs []destinations.BufferConfig, convertToJSONL bool, validationCongfigs []validation.ValidationConfig) {
	fmt.Println(`    ____        __           __    _      __
   / __ \____ _/ /_____ _   / /   (_)____/ /____  ____  ___  _____
  / / / / __  / __/ __  /  / /   / / ___/ __/ _ \/ __ \/ _ \/ ___/
 / /_/ / /_/ / /_/ /_/ /  / /___/ (__  ) /_/  __/ / / /  __/ /
/_____/\__,_/\__/\__,_/  /_____/_/____/\__/\___/_/ /_/\___/_/
                                                                  `)
	fmt.Println("Gudasoft v" + version.Current)

	go printConfigInfo(logConfig, serverConfig, streamerConfigs, bufferSize, buffererConfigs, convertToJSONL, validationCongfigs)
}

func printConfigInfo(logConfig logging.LogConfig, serverConfigs []server.ServerConfig, streamerConfigs []destinations.StreamConfig, bufferSize int, buffererConfigs []destinations.BufferConfig, convertToJSONL bool, validationCongfigs []validation.ValidationConfig) {
	time.Sleep(1 * time.Second)

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	fmt.Printf("Configuration file: %s/%s\n", exeDir, configFile)

	if len(serverConfigs) > 0 {
		fmt.Print(config.GetServerConfigInfo(serverConfigs))
	}
	if len(streamerConfigs) > 0 {
		fmt.Print(config.GetStreamConfigInfo(streamerConfigs))
	}
	if len(buffererConfigs) > 0 {
		if convertToJSONL {
			fmt.Println("Converting input to jsonl")
		}
		fmt.Printf("Buffer size: %d Kilobytes\n", (bufferSize / 1024))
		fmt.Print(config.GetBufferConfigInfo(buffererConfigs))
	}
	if len(validationCongfigs) > 0 {
		fmt.Print(config.GetValidationConfigInfo(validationCongfigs))

	}
	fmt.Println(logConfig)
}
