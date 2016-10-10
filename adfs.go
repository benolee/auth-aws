package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"

	"gopkg.in/ini.v1"
)

type ADFSConfig struct {
	Username string `ini:"user"`
	Password string `ini:"pass"`
	Hostname string `ini:"host"`
}

var (
	settingsPath string = os.Getenv("HOME") + "/.config/auth-aws/config.ini"
)

func loadSettingsFile(adfsConfig *ADFSConfig, settingsFile io.Reader) {
	b, err := ioutil.ReadAll(settingsFile)
	checkError(err)

	cfg, err := ini.Load(b)
	if err == nil {
		err = cfg.Section("adfs").MapTo(adfsConfig)
		checkError(err)
	}
}

func loadEnvVars(adfsConfig *ADFSConfig) {
	if val, ok := os.LookupEnv("ADFS_USER"); ok {
		adfsConfig.Username = val
	}
	if val, ok := os.LookupEnv("ADFS_PASS"); ok {
		adfsConfig.Password = val
	}
	if val, ok := os.LookupEnv("ADFS_HOST"); ok {
		adfsConfig.Hostname = val
	}
}

func loadAskVars(adfsConfig *ADFSConfig) {
	reader := bufio.NewReader(os.Stdin)

	if adfsConfig.Username == "" {
		fmt.Printf("Username: ")
		user, err := reader.ReadString('\n')
		checkError(err)
		adfsConfig.Username = strings.Trim(user, "\n")
	}
	if adfsConfig.Password == "" {
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		checkError(err)
		adfsConfig.Password = string(pass[:])
	}
	if adfsConfig.Hostname == "" {
		fmt.Printf("Hostname: ")
		host, err := reader.ReadString('\n')
		checkError(err)
		adfsConfig.Hostname = strings.Trim(host, "\n")
	}
}

func newADFSConfig() *ADFSConfig {

	adfsConfig := new(ADFSConfig)

	if settingsPath != "" {
		if settingsFile, err := os.Open(settingsPath); err == nil {
			defer settingsFile.Close()
			loadSettingsFile(adfsConfig, settingsFile)
		}
	}

	loadEnvVars(adfsConfig)
	loadAskVars(adfsConfig)

	return adfsConfig
}