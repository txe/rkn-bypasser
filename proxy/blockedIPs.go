package proxy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/set"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var blockedIPs set.Interface

func loadBlockedIPs() error {
	t := time.Now()

	req, err := http.NewRequest("GET", "https://reestr.rublacklist.net/api/v2/ips/json", nil)
	if err != nil {
		return err
	}

	c := http.Client{
		Timeout: time.Minute,
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("not OK response status: " + res.Status)
	}

	if blockedIPs == nil {
		blockedIPs = set.New(set.NonThreadSafe)
	} else {
		blockedIPs.Clear()
	}

	var ips []string

	err = json.NewDecoder(res.Body).Decode(&ips)
	if err != nil {
		return errors.New("failed to decode IPs JSON: " + err.Error())
	}

	for _, ip := range ips {
		blockedIPs.Add(ip)
	}

	logrus.Infof("%d blocked IPs loaded in %s", blockedIPs.Size(),
		time.Now().Sub(t).String())

	return nil
}

func loadPresavedBlockedIPs() error {
	t := time.Now()

	f, err := os.Open("blocked-ips.json")
	if err != nil {
		return err
	}

	var ips []string

	err = json.NewDecoder(f).Decode(&ips)
	if err != nil {
		return errors.New("failed to decode IPs JSON: " + err.Error())
	}

	if blockedIPs == nil {
		blockedIPs = set.New(set.NonThreadSafe)
	} else {
		blockedIPs.Clear()
	}

	for _, ip := range ips {
		blockedIPs.Add(ip)
	}

	logrus.Infof("%d presaved blocked IPs loaded in %s", len(ips),
		time.Now().Sub(t).String())

	return nil
}

func loadAdditionalBlockedIPs() error {
	t := time.Now()

	ipsYAML, err := ioutil.ReadFile("./additional-ips.yml")
	if err != nil {
		return errors.New("failed to read additional-ips file: " + err.Error())
	}

	var ips []string

	err = yaml.Unmarshal(ipsYAML, &ips)
	if err != nil {
		return errors.New("failed to unmarshal additional ip YAML: " + err.Error())
	}

	for _, ip := range ips {
		blockedIPs.Add(ip)
	}

	logrus.Infof("%d additional blocked IPs loaded in %s", len(ips),
		time.Now().Sub(t).String())

	return nil
}

func initBlockedIPs(withAdditionalIPs bool) {
	if err := loadBlockedIPs(); err != nil {
		logrus.WithError(err).Error("Failed to load blocked IPs")
		if err := loadPresavedBlockedIPs(); err != nil {
			logrus.WithError(err).Error(
				"Failed to load presaved blocked IPs")
		}
	}
	if withAdditionalIPs {
		if err := loadAdditionalBlockedIPs(); err != nil {
			logrus.WithError(err).Error("Failed to load additional blocked IPs")
		}
	}
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			if err := loadBlockedIPs(); err != nil {
				logrus.WithError(err).Error("Failed to load blocked IPs")
			}
		}
	}()
}
