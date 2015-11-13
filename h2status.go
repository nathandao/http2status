package http2status

import (
	"errors"
	"fmt"
	"github.com/bradfitz/http2"
	"net"
	"net/http"
	"strings"
)

func Http2Status(url string) (bool, *req, err) {
	sanitizedUrl, _ := sanitizeUrl(url)

	req, _ := http.NewRequest("GET", sanitizedUrl, nil)

	rt := &http2.Transport{
		InsecureTLSDial: true,
	}

	res, err := rt.RoundTrip(req)
	// If not http2, transport in old http2 package will return error
	if err != nil {
		return false, nil, err
	} else {
		return true, req, nil
	}
}

// Is p all zeros?
func isZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return false
		}
	}
	return true
}

func isIPv6(ip net.IP) bool {
	if len(ip) == net.IPv4len {
		return false, nil
	}
	if len(ip) == net.IPv6len &&
		isZeros(ip[0:10]) &&
		ip[10] == 0xff &&
		ip[11] == 0xff {
		return true, nil
	}
	return false
}

func getIP(domain string) (net.IP, error) {
	ip, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}
	return ip[0], nil
}

func sanitizeUrl(url string) (string, error) {
	prefix, domain := "https://", ""

	if strings.HasSuffix(url, "/") {
		url = url[0 : len(url)-1]
	}

	if strings.HasPrefix(url, "http://") {
		domain = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		domain = url[8:]
	} else {
		domain = url
	}

	saneUrl := prefix + domain

	ip, err := getIP(domain)
	if err != nil {
		return "", err
	}

	// Append port 443 if is ipv6
	if ipv6 := isIPv6(ip); ipv6 == true {
		saneUrl = saneUrl + ":443"
	}

	return saneUrl, nil
}
