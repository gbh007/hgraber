package base

import (
	"app/system"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func ParseJSON(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		system.Debug(r.Context(), err)
	}
	return err
}

// GetIP возвращает реальный ip пользователя
func GetIP(r *http.Request) string {
	if ip := r.Header.Get("X-REAL-IP"); ip != "" {
		return ip
	}
	if ips := r.Header.Get("X-FORWARDED-FOR"); ips != "" {
		for _, ip := range strings.Split(ips, ",") {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}
