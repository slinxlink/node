package util

import (
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func GetPublicIPv4() string {
	resp, err := httpClient.Get("https://api4.ipify.org")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(ip)
}

func GetPublicIPv6() string {
	resp, err := httpClient.Get("https://api6.ipify.org")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(ip)
}

func GetPublicIPs() (ipv4, ipv6 string) {
	done := make(chan struct{}, 2)

	go func() {
		ipv4 = GetPublicIPv4()
		done <- struct{}{}
	}()
	go func() {
		ipv6 = GetPublicIPv6()
		done <- struct{}{}
	}()
	<-done
	<-done

	return
}
