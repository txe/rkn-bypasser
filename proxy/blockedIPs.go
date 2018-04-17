package proxy

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fatih/set"
)

var blockedIPs *set.Set

func loadBlockedIPs() error {
	t := time.Now()
	res, err := http.Get("https://reestr.rublacklist.net/api/v2/ips/json")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("not OK response status: " + err.Error())
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

func initBlockedIPs() {
	if err := loadBlockedIPs(); err != nil {
		log.Fatalln("[ERR] Fail to load blocked IPs: " + err.Error())
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
