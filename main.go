package main

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func getCurrentIp() (bool, string) {
	log := logging.MustGetLogger("log")
	urlList := viper.GetStringSlice("serverToGetIp.urlList")

	ip := ""
	findIt := false

	for _, url := range urlList {
		client := http.Client{Timeout: 30 * time.Second}

		resp, err := client.Get(url)
		if err != nil {
			log.Warning("Unable to connect on \"%s\":", url, err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warning("Unable to read data page:", err)
			continue
		}

		ip = strings.Trim(string(body), "\n")
		if ip != "" {
			break
		}
	}

	return findIt, ip
}

func loadOldIp() string {
	file := viper.GetString("config.oldIpFile")

	ip, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return string(ip)
}

//func update_hopper_if_needed(cfg *Config, ip_file *string, old_ip *string, current_ip *string) {
func updateIpIfNeeded(oldIp *string, currentIp *string) {
	log := logging.MustGetLogger("log")
	if *oldIp != *currentIp {
		client := http.Client{Timeout: 30 * time.Second}
		url := ""
		if viper.GetBool("config.ipv4") {
			url = "http://" + viper.GetString("config.user") + ":" + viper.GetString("config.password") + "@ipv4.nsupdate.info/nic/update"
		} else {
			url = "http://" + viper.GetString("config.user") + ":" + viper.GetString("config.password") + "@ipv6.nsupdate.info/nic/update"
		}

		resp, err := client.Get(url)
		if err != nil {
			log.Critical("Unable to connect on \"%s\":", url, err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Critical("Unable to get data page:", err)
			os.Exit(1)
		}

		//if strings.Contains(string(body), "nochg") || strings.Contains(string(body), "good") {
		if strings.Contains(string(body), "good") {
			if err := ioutil.WriteFile(viper.GetString("config.oldIpFile"), []byte(*currentIp), 0644); err != nil {
				log.Critical("Unable to write into \"%s\":", viper.GetString("config.oldIpFile"), err)
				os.Exit(1)
			}
		} else {
			log.Critical("Identification failure when updating IP:", string(body))
			os.Exit(1)
		}
	} else {
		log.Info("Nothing to do")
	}
}

func main() {
	logFile := "/var/log/gonsupdate/errors.log"
	confPath := "/etc/gonsupdate"
	confFilename := "gonsupdate"

	/*
		logFile := "errors.log"
		confPath := "cfg/"
		confFilename := "gonsupdate"
	*/

	fd := initLogging(&logFile)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	oldIp := loadOldIp()
	findIt, currentIp := getCurrentIp()
	if findIt {
		updateIpIfNeeded(&oldIp, &currentIp)
	} else {
		log := logging.MustGetLogger("log")
		log.Critical("Unable to get internet IP address !")
		os.Exit(1)
	}
}
