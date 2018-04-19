package proxy

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/set"
)

var blockedIPs *set.Set

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
		blockedIPs = set.New()
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

	log.Printf("%d blocked IPs loaded in %s\n", blockedIPs.Size(), time.Now().Sub(t).String())
	return nil
}

func loadPresavedBlockedIPs() error {
	t := time.Now()

	f, err := os.Open("blocked-ips")
	if err != nil {
		return err
	}

	var ips []string

	err = json.NewDecoder(f).Decode(&ips)
	if err != nil {
		return errors.New("failed to decode IPs JSON: " + err.Error())
	}

	if blockedIPs == nil {
		blockedIPs = set.New()
	} else {
		blockedIPs.Clear()
	}

	for _, ip := range ips {
		blockedIPs.Add(ip)
	}

	log.Printf("%d presaved blocked IPs loaded in %s\n", blockedIPs.Size(), time.Now().Sub(t).String())
	return nil
}

func initBlockedIPs() {
	if err := loadBlockedIPs(); err != nil {
		log.Println("[ERR] Fail to load blocked IPs: " + err.Error())
		if err := loadPresavedBlockedIPs(); err != nil {
			log.Fatalln("[ERR] Fail to load presaved blocked IPs: " + err.Error())
		}
	}
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			if err := loadBlockedIPs(); err != nil {
				log.Println("[ERR] Fail to load blocked IPs: " + err.Error())
			}
		}
	}()
}
