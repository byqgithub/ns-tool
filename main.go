package main

import (
    "fmt"
	"time"

	command "github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel        = 5
	defaultLogPath         = "/var/log/ns-tools"
	defaultLogName         = "ns-tools.log"
	defaultLogMaxAge       = time.Hour * 1200
	defaultLogRotationTime = time.Hour * 24
)

type logConfig struct {
	LogLevel        int
	LogPath         string
	LogName         string
	LogMaxAge       int64
	LogRotationTime int64
}

type namespaceType struct {
	uts     bool
	ipc     bool
	network bool
	pid     bool
	user    bool
	mount   bool
	cgroup  bool
	all     bool
}

var (
	version = "1.0.0"

	logConf  logConfig
	rootCmd  *command.Command
	nsType   namespaceType
)

// ParseOpt parse opt
func parseOpt()  {
	rootCmd = &command.Command{
		Use:     "ns-tools",
		Short:   "Linux namespace tools",
		Version: version,
	}

	var createCmd = &command.Command{
		Use:     "create",
		Short:   "create linux namespace",
		Args:  command.MaximumNArgs(10),
		RunE: func(cmd *command.Command, args []string) error {
			err := createNS()
			if err != nil {
				fmt.Printf("Collect start error: %v\n", err)
				return err
			}
			return nil
		},
	}

    createCmd.PersistentFlags().BoolVarP(&nsType.uts, "uts", "u", false, "create UTS namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.ipc, "ipc", "i", false, "create IPC namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.network, "net", "n", false, "create network namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.pid, "pid", "p", false, "create PID namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.user, "user", "U", false, "create user namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.mount, "mount", "m", false, "create mount namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.cgroup, "cgroup", "C", false, "create cgroup namespace")
	createCmd.PersistentFlags().BoolVarP(&nsType.all, "all", "a", false, "create all namespace")

	rootCmd.PersistentFlags().IntVar(
		&logConf.LogLevel,
		"log-level",
		defaultLogLevel,
		"log level, 0 panic, 1 fatal, 2 error, 3 warn, 4 info, 5 debug, 6 trace")

	rootCmd.PersistentFlags().StringVar(
		&logConf.LogPath,
		"log-path",
		defaultLogPath,
		"log path")

	rootCmd.PersistentFlags().StringVar(
		&logConf.LogName,
		"log-name",
		defaultLogName,
		"log name")

	rootCmd.PersistentFlags().Int64Var(
		&logConf.LogMaxAge,
		"log-age",
		int64(defaultLogMaxAge),
		"log max age")

	rootCmd.PersistentFlags().Int64Var(
		&logConf.LogRotationTime,
		"log-rotation",
		int64(defaultLogRotationTime),
		"log rotation time")

	var versionCmd = &command.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *command.Command, args []string) {
			fmt.Printf("pi-collect version %s\n", version)
		},
	}

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.SetVersionTemplate(version)
}

func initEnv() error {
	// var err error
	fmt.Printf("Log Path %v\n", logConf.LogPath)
	initLog(
		logConf.LogPath,
		logConf.LogName,
		logConf.LogMaxAge,
		logConf.LogRotationTime,
		logConf.LogLevel)
	log.Infof("version: %v", version)

	// ctx, cancel := context.WithCancel(context.Background())
	// var wg sync.WaitGroup

	// quit
	// sigquit := make(chan os.Signal)
	// signal.Notify(sigquit, syscall.SIGTERM, syscall.SIGINT)
	// log.Infoln("Start normally")
	// select {
	// case <-sigquit:
	// 	log.Infoln("Stop sign")
	// 	executor.Stop(cancel)
	// 	wg.Wait()
	// 	log.Infoln("Quit normally")
	// 	return nil
	// }
	return nil
}

func createNS() error {
	// if err := initEnv(); err != nil {
	// 	return err
	// }

	return createNamespace(&nsType)
}

func main() {
	parseOpt()

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		log.Errorf("Main panic %v", err)
	// 	}
	// }()

	if err := initEnv(); err != nil {
		fmt.Printf("Init env error: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Collect running error: %v", err)
	}
}
