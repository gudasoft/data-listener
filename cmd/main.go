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

	"buffer-handler/config"
	"buffer-handler/destinations"
	"buffer-handler/logging"
	"buffer-handler/models"
	"buffer-handler/server"
	"buffer-handler/version"
	"buffer-handler/workers"
)

var shutdownChannel = make(chan struct{})
var reloadChan = make(chan bool)
var configFile = "config.toml"

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
	var bufferSize int
	var convertToJSONL bool
	config.LoadConfigs(configFile, &logConfig, &serverConfig, &streamConfigs, &buffererConfigs, &bufferSize, &convertToJSONL)
	_, err := logConfig.InitLogger()
	if err != nil {
		panic(err)
	}

	logConfigs(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)

	var streamerChannel chan models.EntryData
	if len(streamConfigs) > 0 {
		streamerChannel = make(chan models.EntryData)
		go workers.DestinationsNotifier(&streamConfigs, streamerChannel)
	}

	var buffererChannel chan models.EntryData
	if len(buffererConfigs) > 0 {
		buffererChannel = make(chan models.EntryData)
		go workers.DestinationsNotifierBuffered(&bufferSize, &buffererConfigs, &convertToJSONL, buffererChannel, shutdownChannel, reloadChan)
	}

	fastHttpServer := server.StartServer(serverConfig, streamerChannel, buffererChannel)

	printWelcomeBanner(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)

sigloop:
	for {
		select {
		case sigReceived := <-signalChannel:
			if sigReceived == syscall.SIGTERM || sigReceived == syscall.SIGINT {
				close(shutdownChannel)
				logging.Logger.Info("Received shutdown signal. Cleaning up...")
				time.Sleep(100 * time.Millisecond)
				fastHttpServer.Shutdown()
				break sigloop
			} else if sigReceived == syscall.SIGHUP {
				logging.Logger.Info("Received configuration reload signal. Loading up...")
				reloadChan <- true
				config.LoadConfigs(configFile, &logConfig, &serverConfig, &streamConfigs, &buffererConfigs, &bufferSize, &convertToJSONL)
				logConfig.InitLogger()
				time.Sleep(100 * time.Millisecond)
				logConfigs(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)
				printConfigInfo(logConfig, serverConfig, streamConfigs, bufferSize, buffererConfigs, convertToJSONL)
			}
		case <-ctx.Done():
			close(shutdownChannel)
			logging.Logger.Info("Context canceled. Cleaning up...")
			time.Sleep(100 * time.Millisecond)
			fastHttpServer.Shutdown()
			break sigloop
		}
	}
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

func printWelcomeBanner(logConfig logging.LogConfig, serverConfig []server.ServerConfig, streamerConfigs []destinations.StreamConfig, bufferSize int, buffererConfigs []destinations.BufferConfig, convertToJSONL bool) {
	fmt.Println(`    ____        __           __    _      __
   / __ \____ _/ /_____ _   / /   (_)____/ /____  ____  ___  _____
  / / / / __  / __/ __  /  / /   / / ___/ __/ _ \/ __ \/ _ \/ ___/
 / /_/ / /_/ / /_/ /_/ /  / /___/ (__  ) /_/  __/ / / /  __/ /
/_____/\__,_/\__/\__,_/  /_____/_/____/\__/\___/_/ /_/\___/_/
                                                                  `)
	fmt.Println("Gudasoft v" + version.Current)

	go printConfigInfo(logConfig, serverConfig, streamerConfigs, bufferSize, buffererConfigs, convertToJSONL)
}

func printConfigInfo(logConfig logging.LogConfig, serverConfigs []server.ServerConfig, streamerConfigs []destinations.StreamConfig, bufferSize int, buffererConfigs []destinations.BufferConfig, convertToJSONL bool) {
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
		fmt.Printf("Buffer size: %d Megabytes\n", ((bufferSize / 1024) / 1024))
		fmt.Print(config.GetBufferConfigInfo(buffererConfigs))
	}
	fmt.Println(logConfig)
}
