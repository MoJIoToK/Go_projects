package middleware

import (
	"net"
	"net/http"
	"strings"
)

// xForwardedFor - каноничный формат названия заголовка X-Forwarded-For.
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")

// xRealIP - каноничный формат названия заголовка X-Real-IP.
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

// RealIP устанавливает заголовок X-Real-IP с адресом клиента.
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerIP := headerIP(r)
		if headerIP != "" {
			portExc := headerIP
			i := strings.Index(headerIP, ":")
			if i != -1 {
				portExc = headerIP[:i]
			}
			if net.ParseIP(portExc) != nil {
				r.Header.Set(xRealIP, headerIP)
			}
		} else {
			remote := r.RemoteAddr
			r.Header.Set(xRealIP, remote)
		}
		next.ServeHTTP(w, r)
	})
}

// headerIP - проверяет заголовки с адресом клиента.
func headerIP(r *http.Request) string {
	real := r.Header.Get(xRealIP)
	if real != "" {
		return real
	}

	real = r.Header.Get(xForwardedFor)
	if real != "" {
		i := strings.Index(real, ",")
		if i != -1 {
			real = real[:i]
		}
	}
	return real
}
