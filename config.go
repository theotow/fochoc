package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// FileName of the config
const FileName = ".fochocconfig.json"

// FileMode is the file mode the config will be created in
const FileMode = 0700

// EnvConfig reads the config from ENV variables
type EnvConfig struct{}

// FileConfig reads the config from the config file
type FileConfig struct {
	config ConfigFileStruct
}

// Token describes a erc20 token
type Token struct {
	Comment string
	Address string
}

// ConfigFileStruct describes the json config file structure
type ConfigFileStruct struct {
	Keys        map[string]string `json:"keys"`
	Erc20Tokens []Token           `json:"erc20Tokens"`
}

func (c *ConfigFileStruct) addKey(key string, value string) {
	if c.Keys == nil {
		c.Keys = make(map[string]string)
	}
	c.Keys[key] = value
}

func (c *ConfigFileStruct) addErc20(token Token) {
	if c.Erc20Tokens == nil {
		c.Erc20Tokens = []Token{}
	}
	c.Erc20Tokens = append(c.Erc20Tokens, token)
}

// ConfigInterface describes the interface a config needs to have
type ConfigInterface interface {
	GetKey(name string) string
	GetTokens() []Token
}

// ConfigProviderInterface describes the interface a a provider needs to have
type ConfigProviderInterface interface {
	GetCurrencyValue(name string) float64
	GetAll(keys []string) []BalanceSimple
	AddTestBalance(name string, value float64)
}

// NewEnvConfig generates a new ENV config struct
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

// NewFileConfig generates a new File config struct
func NewFileConfig() *FileConfig {
	config := readFile(getFileString())
	return &FileConfig{config: config}
}

// GetKey returns a secrets / keys for exchange providers
func (c *EnvConfig) GetKey(name string) string {
	return os.Getenv(name)
}

// GetTokens returns an array of Erc20 tokens
func (c *EnvConfig) GetTokens() []Token {
	return []Token{
		{Address: "0xf27d22d64e625c2a34e31369d9b88828146df52b", Comment: "comment"},
	}
}

// GetKey returns a secrets / keys for exchange providers
func (c *FileConfig) GetKey(name string) string {
	if val, ok := c.config.Keys[name]; ok {
		return val
	}
	return ""
}

// Write writes the fileconfig struct to the config file
func (c *FileConfig) Write(data ConfigFileStruct) {
	writeFile(getFileString(), data)
}

// Read reads the fileconfig struct from the config file
func (c *FileConfig) Read() ConfigFileStruct {
	return c.config
}

// GetEmptyConfig returns an empty config struct
func (c *FileConfig) GetEmptyConfig() ConfigFileStruct {
	return ConfigFileStruct{}
}

// GetTokens returns an array of Erc20 tokens
func (c *FileConfig) GetTokens() []Token {
	return c.config.Erc20Tokens
}

// TODO: put in main.go
func getFileString() string {
	usr, err := user.Current()
	if err != nil {
		panic(errors.New("cannot get homedir"))
	}
	return strings.Join([]string{usr.HomeDir, "/", FileName}, "")
}

func check(err error, msg string) {
	if err != nil {
		panic(errors.New(msg))
	}
}

func readFile(path string) ConfigFileStruct {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, _ := os.Create(path)
		if f != nil {
			f.Close()
		}
	}
	fileData, err := ioutil.ReadFile(path)
	check(err, "cannot open file #{path}"+path)
	res := ConfigFileStruct{}
	err = json.Unmarshal(fileData, &res)
	if err != nil {
		return ConfigFileStruct{} // empty
	}
	return res
}

func writeFile(path string, data ConfigFileStruct) {
	dataBytes, err := json.Marshal(&data)
	check(err, "cannot marshal")
	errIo := ioutil.WriteFile(path, dataBytes, FileMode)
	if err != nil {
		panic(errIo)
	}
}
