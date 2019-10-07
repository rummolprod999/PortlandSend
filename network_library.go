package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func DownloadPage(url string) string {
	count := 0
	var st string
	for {
		//fmt.Println("Start download file")
		if count > 5 {
			Logging(fmt.Sprintf("The file was not downloaded in %d attemps %s", count, url))
			return st
		}
		st = GetPageUA(url)
		if st == "" {
			count++
			//Logging("Got empty page", url)
			time.Sleep(time.Second * 5)
			continue
		}
		return st

	}
	return st
}

func GetPageUA(url string) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			Logging(fmt.Sprintf("Was panic, recovered value: %v", r))
			ret = ""
		}
	}()
	var st string
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Logging("Error request", url, err)
		return st
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		Logging("Error download", url, err)
		return st
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logging("Error reading", url, err)
		return st
	}

	return string(body)
}

func DownloadPageComita(url, kw string) string {
	count := 0
	var st string
	for {
		//fmt.Println("Start download file")
		if count > 5 {
			Logging(fmt.Sprintf("The file was not downloaded in %d attemps %s", count, url))
			return st
		}
		st = GetPageUAComita(url, kw)
		if st == "" {
			count++
			Logging("Got empty page", url)
			time.Sleep(time.Second * 5)
			continue
		}
		return st

	}
	return st
}

func GetPageUAComita(url, kw string) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			Logging(fmt.Sprintf("Was panic, recovered value: %v", r))
			ret = ""
		}
	}()
	var st string
	client := &http.Client{}
	bodyReq := []byte(fmt.Sprintf("{\"itemsPerPage\":20,\"page\":0,\"search\":\"%s\",\"orderBy\":\"startDate\",\"orderAsc\":false}", kw))
	rBody := bytes.NewReader(bodyReq)
	request, err := http.NewRequest("POST", url, rBody)
	if err != nil {
		Logging("Error request", url, err)
		return st
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8'")
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		Logging("Error download", url, err)
		return st
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logging("Error reading", url, err)
		return st
	}

	return string(body)
}
