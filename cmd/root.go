package cmd

import (
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yiplee/blockquiz/config"
	"github.com/yiplee/blockquiz/version"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:           "blockquiz",
	Short:         "a quiz dapp running on mixin messenger",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		initLogger(debug)
	},
}

func Execute() {
	rootCmd.Version = version.String()
	rand.Seed(time.Now().UnixNano())

	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blockquiz.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	rootCmd.PersistentFlags().String("dialect", "mysql", "db dialect")
	_ = viper.BindPFlag("db.dialect", rootCmd.PersistentFlags().Lookup("dialect"))
	rootCmd.PersistentFlags().String("dbhost", "", "db host")
	_ = viper.BindPFlag("db.host", rootCmd.PersistentFlags().Lookup("dbhost"))
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			logrus.Fatal(err)
		}

		cfgFile = path.Join(home, ".blockquiz.yaml")
	}

	c, err := config.Load(cfgFile)
	if err != nil {
		logrus.Fatal(err)
	}

	cfg = c
}

func initLogger(enableDebug bool) {
	if enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Info("log level: debug")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
