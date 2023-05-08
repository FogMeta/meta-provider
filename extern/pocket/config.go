package pocket

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"time"
)

type PoktConfig struct {
	SwanApiUrl      string `toml:"swan_api_url"`
	SwanApiKey      string `toml:"swan_api_key"`
	SwanAccessToken string `toml:"swan_access_token"`

	PoktLogLevel          string        `toml:"pokt_log_level"`
	PoktApiUrl            string        `toml:"pokt_api_url"`
	PoktAddress           string        `toml:"pokt_address"`
	PoktDockerImage       string        `toml:"pokt_docker_image"`
	PoktDockerName        string        `toml:"pokt_docker_name"`
	PoktConfigPath        string        `toml:"pokt_path"`
	PoktScanInterval      time.Duration `toml:"pokt_scan_interval"`
	PoktHeartbeatInterval time.Duration `toml:"pokt_heartbeat_interval"`
	PoktServerApiUrl      string        `toml:"pokt_server_api_url"`
	PoktServerApiPort     int           `toml:"pokt_server_api_port"`
	PoktNetworkType       string        `toml:"pokt_network_type"`
}

var config *PoktConfig

func GetConfig() PoktConfig {
	if config == nil {
		initConfig()
	}
	return *config
}

func initConfig() {
	configPath := os.Getenv("SWAN_PATH")
	if configPath == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			GetLog().Fatal("Cannot get home directory.")
		}

		configPath = filepath.Join(homedir, ".swan/")
	}

	initPoktConfig(filepath.Join(configPath, "provider/config-pokt.toml"))

}

func initPoktConfig(configFile string) {
	GetLog().Debug("Your pokt config file is:", configFile)

	if metaData, err := toml.DecodeFile(configFile, &config); err != nil {
		GetLog().Fatal("error:", err)
	} else {
		if !requiredPoktAreGiven(metaData) {
			GetLog().Fatal("required fields not given")
		}
	}
}

func requiredPoktAreGiven(metaData toml.MetaData) bool {
	requiredFields := [][]string{
		{"swan_api_url"},
		{"swan_api_key"},
		{"swan_access_token"},
		{"pokt_log_level"},
		{"pokt_api_url"},
		{"pokt_docker_image"},
		{"pokt_docker_name"},
		{"pokt_path"},
		{"pokt_scan_interval"},
		{"pokt_heartbeat_interval"},
		{"pokt_server_api_url"},
		{"pokt_server_api_port"},
		{"pokt_network_type"},
	}

	for _, v := range requiredFields {
		if !metaData.IsDefined(v...) {
			GetLog().Fatal("required conf fields ", v)
		}
	}

	return true
}

func ReadPocketChains() string {
	configPath := os.Getenv("SWAN_PATH")
	if configPath == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			GetLog().Fatal("Cannot get home directory.")
			return ""
		}

		configPath = filepath.Join(homedir, ".swan/")
	}

	chainsFile, err := os.Open(filepath.Join(configPath, "provider/chains.json"))
	if err != nil {
		GetLog().Error("Open chains.json error:", err)
		return ""
	}
	defer chainsFile.Close()

	fileInfo, err := os.Stat(filepath.Join(configPath, "provider/chains.json"))
	if err != nil {
		GetLog().Error("Stat chains.json error:", err)
		return ""
	}
	chainSize := fileInfo.Size()

	content := make([]byte, chainSize)
	count, err := chainsFile.Read(content)
	if err != nil {
		GetLog().Error("Read chains.json error:", err)
		return ""
	}

	chains := string(content[:count])
	//GetLog().Info("chains.json :\n", chains)
	return chains
}
