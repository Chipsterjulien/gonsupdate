package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

func main() {
	logFile := "/var/log/gonsupdate/errors.log"
	confPath := "/etc/gonsupdate"
	confFilename := "gonsupdate"

	fd := initLogging(&logFile)
	defer fd.Close()

	loadConfig(&confPath, &confFilename)
	oldIP := loadOldIP()
	findIt, currentIP := getCurrentIP()
	if findIt {
		updateIPIfNeeded(&oldIP, &currentIP)
	} else {
		log := logging.MustGetLogger("log")
		log.Critical("Unable to get internet IP address !")
		os.Exit(1)
	}
}

func getCurrentIP() (bool, string) {
	log := logging.MustGetLogger("log")
	urlList := viper.GetStringSlice("servertogetip.urlList")

	shuffleList(&urlList)

	var ip string
	findIt := false

	for _, url := range urlList {
		log.Debugf("Finding ip with this url: %s", url)
		client := http.Client{Timeout: 30 * time.Second}

		resp, err := client.Get(url)
		if err != nil {
			log.Warningf("Unable to connect on \"%s\": %s", url, err)
			sendAnEMail("gonsupdate", fmt.Sprintf("Unable to connect on \"%s\": %s", url, err))
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warningf("Unable to read data page: %s", err)
			sendAnEMail("gonsupdate", fmt.Sprintf("Unable to read data on \"%s\": %s", url, err))
			continue
		}

		ip = strings.Trim(string(body), "\n")
		log.Debugf("IP is \"%s\" and it find on \"%s\"", ip, url)
		if ip != "" {
			findIt = true
			break
		}
	}

	log.Debugf("Is IP finding: %s", findIt)
	if findIt {
		log.Debugf("IP is: %s", ip)
	}

	return findIt, ip
}

func loadOldIP() string {
	log := logging.MustGetLogger("log")
	log.Debugf("Try to get old IP in \"%s\" file", viper.GetString("config.oldIpFile"))
	file := viper.GetString("config.oldIpFile")

	ip, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	log.Debugf("IP will be return is: %s", ip)

	return string(ip)
}

func shuffleList(lst *[]string) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	for i := len(*lst) - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		(*lst)[i], (*lst)[j] = (*lst)[j], (*lst)[i]
	}
}

func updateIPIfNeeded(oldIP *string, currentIP *string) {
	log := logging.MustGetLogger("log")

	log.Debugf("Old IP is: \"%s\"\nCurrent IP is: \"%s\"", *oldIP, *currentIP)

	if *oldIP != *currentIP {
		log.Debug("IP are differents")

		client := http.Client{Timeout: 30 * time.Second}

		var url string
		if viper.GetBool("config.ipv4") {
			url = "http://" + viper.GetString("config.user") + ":" + viper.GetString("config.password") + "@ipv4.nsupdate.info/nic/update"
		} else {
			url = "http://" + viper.GetString("config.user") + ":" + viper.GetString("config.password") + "@ipv6.nsupdate.info/nic/update"
		}

		resp, err := client.Get(url)
		if err != nil {
			log.Criticalf("Unable to connect on \"%s\": %s", url, err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Criticalf("Unable to get data page: %s", err)
			os.Exit(1)
		}

		if strings.Contains(string(body), "nochg") || strings.Contains(string(body), "good") {
			log.Debugf("Writing good IP in \"%s\" file", viper.GetString("config.oldIpFile"))
			if err := ioutil.WriteFile(viper.GetString("config.oldIpFile"), []byte(*currentIP), 0644); err != nil {
				log.Criticalf("Unable to write into \"%s\": %s", viper.GetString("config.oldIpFile"), err)
				os.Exit(1)
			}
		} else {
			log.Criticalf("Identification failure when updating IP: %s", string(body))
			os.Exit(1)
		}
	} else {
		log.Debug("Nothing to do")
	}
}
