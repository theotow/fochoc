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
const FileName = ".crypto.json"
const FileMode = 0700

type EnvConfig struct{}
type FileConfig struct {
	config map[string]string
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

func (c *FileConfig) GetKey(name string) string {
	if val, ok := c.config[name]; ok {
		return val
	}
	return ""
}

func (c *FileConfig) Initialised() bool {
	return len(c.config) > 0
}

func (c *FileConfig) Write(data map[string]string) {
	writeFile(getFileString(), data)
}

func (c *FileConfig) Read() map[string]string {
	return c.config
}

type ConfigInterface interface {
	GetKey(name string) string
	Initialised() bool
}

type ConfigProviderInterface interface {
	GetCurrencyValue(name string) float64
	GetAll(keys []string) map[string]float64
	AddTestBalance(name string, value float64)
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

func readFile(path string) map[string]string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, _ := os.Create(path)
		if f != nil {
			f.Close()
		}
	}
	fileData, err := ioutil.ReadFile(path)
	check(err, "cannot open file #{path}"+path)
	jsonMap := make(map[string]string)
	err = json.Unmarshal(fileData, &jsonMap)
	if err != nil {
		return make(map[string]string)
	}
	return jsonMap
}

func writeFile(path string, data map[string]string) {
	dataBytes, err := json.Marshal(&data)
	check(err, "cannot marshal")
	ioutil.WriteFile(path, dataBytes, FileMode)
}
