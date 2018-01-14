package proxy

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/set"
)

var blockedIPs *set.Set

func loadBlockedIPs() error {
	t := time.Now()
	res, err := http.Get("https://reestr.rublacklist.net/api/ips")
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("not OK response status: " + err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("fail read reponse body: " + err.Error())
	}
	body := string(bodyBytes)
	body = strings.TrimPrefix(body, `"`)
	body = strings.TrimSuffix(body, `"`)
	if blockedIPs == nil {
		blockedIPs = set.New()
	} else {
		blockedIPs.Clear()
	}
	for _, ip := range strings.Split(body, ";") {
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
