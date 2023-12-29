package nftcli

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tak1827/light-nft-indexer/log"
	"github.com/tak1827/light-nft-indexer/util"
)

const (
	EnvPrefix       = "nftcli"
	ConfigType      = "toml"
	ConfigName      = "conf"
	DefaultHomeName = ".nftcli"
)

var (
	ErrInvalidConfig = errors.New("invalid conf")

	homeDir      string
	loggingLevel string
	logWithColor bool
	logger       zerolog.Logger

	Env string

	AwsRegion  string
	AwsProfile string
	// s3
	S3BucketName string
	// blockchain
	ChainEndpoint string
	ChainTimeout  int64
	FaucetPrivKey string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nftcli",
	Short: "The cli tool of light NFT Indexer",
	Long:  `The command line interface of light NFT Indexer. Tool provide 2 kind of indexing. one is subscription of blockchain event, another is fething from logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()
		fmt.Println("use `-h` option to see help")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", fmt.Sprintf("the home directory path (default is $HOME/%s/)", DefaultHomeName))
	rootCmd.PersistentFlags().StringVarP(&loggingLevel, "log-level", "l", "", fmt.Sprintf("the logging level: debug|info|warn|err|fatal (default is /%s/)", log.ToStr(log.DefaultLevel)))
	rootCmd.PersistentFlags().BoolVarP(&logWithColor, "log-with-color", "c", false, "flag of logging with coloer, dafault is false")
}

// initConfig reads in conf file and ENV variables if set.
func initConfig() {
	// Load .env file if exists
	if err := util.LoadDotenv(".env"); err != nil && !util.IsNoEnvErr(err) {
		panic(err)
	}

	// viper.AddConfigPath(".")
	viper.SetConfigType(ConfigType)
	viper.SetConfigName(ConfigName)
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
}

func handleErr(err error) {
	if err != nil {
		logger.Error().Stack().Err(err).Msg("failed executuion")
		logger.Fatal().Msg("stop")
	}
}

func getConfig() {
	if homeDir == "" {
		home, err := os.UserHomeDir()
		handleErr(err)
		homeDir = home + "/" + DefaultHomeName
	}
	viper.AddConfigPath(homeDir)

	if loggingLevel != "" {
		log.SetLevel(log.ToLogLevel(loggingLevel))
	}

	logger = log.CLI()
	if err := viper.ReadInConfig(); err == nil {
		logger.Info().Msgf("using conf file: %s\n", viper.ConfigFileUsed())
	}
}

func getConfigString(key string, val *string) {
	if *val == "" {
		if *val = viper.GetString(key); *val == "" {
			logger.Fatal().Msgf("no `%s` setting", key)
		}
	}
	logger.Info().Msgf("%s: %s", key, *val)
}

func getAwsConfig() {
	// region
	AwsRegion = viper.GetString("aws_region")
	logger.Info().Msgf("aws_region: %s", AwsRegion)
	// profile
	AwsProfile = viper.GetString("aws_profile")
	logger.Info().Msgf("aws_profile: %s", AwsProfile)
}

func getS3Config() {
	// endpoint
	S3BucketName = viper.GetString("s3_bucket_name")
	logger.Info().Msgf("s3_bucket_name: %s", S3BucketName)
}

func getChainConfig() {
	// endpoint
	ChainEndpoint = viper.GetString("blockchain_endpoint")
	logger.Info().Msgf("blockchain_endpoint: %s", ChainEndpoint)
	// timeout
	ChainTimeout = viper.GetInt64("blockchain_timeout")
	logger.Info().Msgf("blockchain_timeout: %d", ChainTimeout)

	if FaucetPrivKey = viper.GetString("faucet_priv_key"); FaucetPrivKey == "" {
		logger.Warn().Msg("please set `STOCLI_FAUCET_PRIV_KEY` as the env variable")
		// dummy key
		FaucetPrivKey = "0000000000000000000000000000000000000000000000000000000000000001"
	}
}

func getEnvConfig() {
	Env = viper.GetString("env")
	logger.Info().Msgf("env: %s", Env)
}
