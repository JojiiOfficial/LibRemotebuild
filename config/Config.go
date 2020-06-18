package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	libremotebuild "github.com/JojiiOfficial/LibRemotebuild"
	"github.com/JojiiOfficial/configService"
	"github.com/JojiiOfficial/gaw"
	"github.com/denisbrodbeck/machineid"
	"github.com/fatih/color"
	"github.com/zalando/go-keyring"
	"gopkg.in/yaml.v2"
)

// ...
const (
	// File constants
	DataDir           = ".remotebuild"
	DefaultConfigFile = "config.yaml"

	// Keyring constants
	DefaultKeyring     = "login"
	KeyringServiceName = "RemoteBuild"
)

var (
	// ErrUnlockingKeyring error if keyring is available but can't be unlocked
	ErrUnlockingKeyring = errors.New("Error unlocking keyring")
)

// Config Configuration structure
type Config struct {
	File      string
	MachineID string
	User      userConfig

	Server serverConfig

	DataManager dataManager
}

type userConfig struct {
	Username       string
	SessionToken   string
	DisableKeyring bool
	Keyring        string
	ForceVerify    bool
}

type serverConfig struct {
	URL        string `required:"true"`
	IgnoreCert bool
}

type dataManager struct {
	Namespaces map[string]string
}

// GetDefaultConfigFile return path of default config
func GetDefaultConfigFile() string {
	return filepath.Join(getDataPath(), DefaultConfigFile)
}

func getDefaultConfig() Config {
	return Config{
		MachineID: GenMachineID(),
		Server: serverConfig{
			URL:        "http://localhost:9999",
			IgnoreCert: false,
		},
		User: userConfig{
			DisableKeyring: false,
			Keyring:        DefaultKeyring,
			ForceVerify:    false,
		},
		DataManager: dataManager{
			Namespaces: map[string]string{
				libremotebuild.JobAUR.String(): "AURbuild",
			},
		},
	}
}

// InitConfig inits the configfile
func InitConfig(defaultFile, file string) (*Config, error) {
	var needCreate bool
	var config Config

	if len(file) == 0 {
		file = defaultFile
		needCreate = true
	}

	// Check if config already exists
	_, err := os.Stat(file)
	needCreate = err != nil

	if needCreate {
		// Autocreate folder
		path, _ := filepath.Split(file)
		_, err := os.Stat(path)
		if err != nil {
			err = os.MkdirAll(path, 0700)
			if err != nil {
				return nil, err
			}
		}

		// Set config to default config
		config = getDefaultConfig()
		config.File = file
	}

	// Create config file if not exists and fill it with the default values
	isDefault, err := configService.SetupConfig(&config, file, configService.NoChange)
	if err != nil {
		return nil, err
	}

	// Return if created but further steps are required
	if isDefault {
		if needCreate {
			return nil, nil
		}
	}

	// Load configuration
	if err = configService.Load(&config, file); err != nil {
		return nil, err
	}

	config.File = file
	config.SetMachineID()

	return &config, nil
}

// SetMachineID sets machineID if empty
func (config *Config) SetMachineID() {
	if len(config.MachineID) == 0 {
		config.MachineID = GenMachineID()
		config.Save()
	}
}

// Validate check the config
func (config *Config) Validate() error {
	// Put in your validation logic here
	return nil
}

// GetMachineID returns the machineID
func (config *Config) GetMachineID() string {
	// Gen new MachineID if empty
	if len(config.MachineID) == 0 {
		config.SetMachineID()
	}

	// Check length of machineID
	if len(config.MachineID) > 100 {
		fmt.Println("Warning: MachineID too big")
		return ""
	}

	return config.MachineID
}

// IsLoggedIn return true if sessiondata is available
func (config *Config) IsLoggedIn() bool {
	if len(config.User.Username) == 0 {
		return false
	}

	var token string
	var err error

	if !config.User.DisableKeyring {
		token, err = keyring.Get(KeyringServiceName, config.User.Username)
	}

	// If no keyring was found, use unencrypted token
	if config.User.DisableKeyring || err != nil {
		token = config.User.SessionToken
	}

	return IsTokenValid(token)
}

// IsTokenValid return true if given token is
// a vaild session token
func IsTokenValid(token string) bool {
	return len(token) == 64
}

// GetKeyring returns the keyring to use
func (config *Config) GetKeyring() string {
	if len(config.User.Keyring) == 0 {
		return DefaultKeyring
	}

	return config.User.Keyring
}

// View view config
func (config Config) View(redactSecrets bool) string {
	// React secrets if desired
	if redactSecrets {
		config.User.SessionToken = "<redacted>"
	}

	// Create yaml
	ymlB, err := yaml.Marshal(config)
	if err != nil {
		return err.Error()
	}

	return string(ymlB)
}

// InsertUser insert a new user
func (config *Config) InsertUser(user, token string) {
	config.User.Username = user
	config.MustSetToken(token)
}

// SetToken sets token for client
// Tries to save token in a keyring, if not supported
// save it unencrypted
func (config *Config) SetToken(token string) error {
	var err error
	if !config.User.DisableKeyring {
		// Save to keyring. Exit return on success
		if err = keyring.Set(KeyringServiceName, config.User.Username, token); err == nil {
			return nil
		}
	}

	fmt.Printf("Your platform doesn't have support for a keyring. Refer to https://github.com/JojiiOfficial/RemoteBuildClient#keyring\n--> !!! Your token will be saved %s !!! <--\n", color.HiRedString("UNENCRYPTED"))

	// Save sessiontoken in config unencrypted
	config.User.SessionToken = token
	return config.Save()
}

// MustSetToken fatals on error
func (config *Config) MustSetToken(token string) {
	if err := config.SetToken(token); err != nil {
		log.Fatal(err)
	}
}

// GetToken returns user token
func (config *Config) GetToken() (string, error) {
	var token string
	var err error

	if !config.User.DisableKeyring {
		token, err = keyring.Get(KeyringServiceName, config.User.Username)
	}

	if config.User.DisableKeyring || err != nil {
		// Return unlock error if sessiontoken is empty,
		// to allow using the unencrypted version
		if IsUnlockError(err) && len(config.User.SessionToken) == 0 {
			return "", ErrUnlockingKeyring
		}

		// If keyring can be opened, but key was not found
		// Return error
		if err == keyring.ErrNotFound {
			return "", err
		}

		// Otherwise return the error and sessiontoken
		return config.User.SessionToken, nil
	}

	return token, nil
}

// ClearKeyring removes session from keyring
func (config *Config) ClearKeyring(username string) error {
	if config.User.DisableKeyring {
		return nil
	}

	if len(username) == 0 {
		username = config.User.Username
	}

	return keyring.Delete(KeyringServiceName, username)
}

// IsUnlockError return true if err is unlock error
func IsUnlockError(err error) bool {
	if err == nil {
		return false
	}

	return strings.HasPrefix(err.Error(), "failed to unlock correct collection") || err == ErrUnlockingKeyring
}

// IsDefault returns true if config is equal to the default config
func (config Config) IsDefault() bool {
	defaultConfig := getDefaultConfig()
	return config.User == defaultConfig.User &&
		config.Server.IgnoreCert == defaultConfig.Server.IgnoreCert
}

// MustGetRequestConfig create a libdm requestconfig from given cli client config and fatal on error
func (config Config) MustGetRequestConfig() *libremotebuild.RequestConfig {
	token, err := config.GetToken()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &libremotebuild.RequestConfig{
		MachineID:    config.GetMachineID(),
		URL:          config.Server.URL,
		IgnoreCert:   config.Server.IgnoreCert,
		SessionToken: token,
		Username:     config.User.Username,
	}
}

// ToRequestConfig create a libdm requestconfig from given cli client config
// If token is not set, error has a value and token is equal to an empty string
func (config Config) ToRequestConfig() (*libremotebuild.RequestConfig, error) {
	token, err := config.GetToken()
	return &libremotebuild.RequestConfig{
		MachineID:    config.GetMachineID(),
		URL:          config.Server.URL,
		IgnoreCert:   config.Server.IgnoreCert,
		SessionToken: token,
		Username:     config.User.Username,
	}, err
}

// Save saves the config
func (config *Config) Save() error {
	return configService.Save(config, config.File)
}

// GenMachineID detect the machineID.
// If not detected return random string
func GenMachineID() string {
	username := getPseudoUsername()

	// Protect with username to allow multiple user
	// on a system using the same manager username
	id, err := machineid.ProtectedID(username)
	if err == nil {
		return id
	}

	// If not detected reaturn random string
	return gaw.RandString(60)
}

func getPseudoUsername() string {
	var username string
	user, err := user.Current()
	if err != nil {
		username = gaw.RandString(10)
	} else {
		username = user.Username
	}

	return username
}

func getDataPath() string {
	path := filepath.Join(gaw.GetHome(), DataDir)
	s, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0700)
		if err != nil {
			log.Fatalln(err.Error())
		}
	} else if s != nil && !s.IsDir() {
		log.Fatalln("DataPath-name already taken by a file!")
	}
	return path
}

// GetNamspace return namespace to use for a given job
func (config *Config) GetNamspace(jobType libremotebuild.JobType) string {
	return config.DataManager.Namespaces[jobType.String()]
}
