package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

// TODO: put in main.go
const FileName = ".fochocconfig.json"
const FileMode = 0700

type EnvConfig struct{}
type FileConfig struct {
	config ConfigFileStruct
}

type Token struct {
	Comment string
	Address string
}

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

type ConfigInterface interface {
	GetKey(name string) string
	GetTokens() []Token
	Initialised() bool
}

type ConfigProviderInterface interface {
	GetCurrencyValue(name string) float64
	GetAll(keys []string) []BalanceSimple
	AddTestBalance(name string, value float64)
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

func NewFileConfig() *FileConfig {
	config := readFile(getFileString())
	return &FileConfig{config: config}
}

func (c *EnvConfig) GetKey(name string) string {
	return os.Getenv(name)
}

func (c *EnvConfig) Initialised() bool {
	return true
}

func (c *EnvConfig) GetTokens() []Token {
	return []Token{
		{Address: "0xf27d22d64e625c2a34e31369d9b88828146df52b", Comment: "comment"},
	}
}

func (c *FileConfig) GetKey(name string) string {
	if val, ok := c.config.Keys[name]; ok {
		return val
	}
	return ""
}

func (c *FileConfig) Initialised() bool {
	return true
}

func (c *FileConfig) Write(data ConfigFileStruct) {
	writeFile(getFileString(), data)
}

func (c *FileConfig) Read() ConfigFileStruct {
	return c.config
}

func (c *FileConfig) GetEmptyConfig() ConfigFileStruct {
	return ConfigFileStruct{}
}

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
