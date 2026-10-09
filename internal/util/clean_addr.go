package util

import (
	"fmt"
	"net/url"
	"strings"
)

func CleanAddr(addr string, ssl bool) (string, error) {
	addr = strings.TrimSpace(addr)

	if !strings.Contains(addr, "://") {
		if ssl {
			addr = "https://" + addr
		} else {
			addr = "http://" + addr
		}
	}

	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}

	if u.Hostname() == "" {
		return "", fmt.Errorf("invalid address: %s", addr)
	}

	return u.Scheme + "://" + u.Host, nil
}
