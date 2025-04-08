package api

import "net/http"

// copyHeader копирует все заголовки из одного http.Header в другой.
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// request подготавливает запрос и выполняет его через переданный http.Client.
func request(host string, r *http.Request, cl *http.Client) (*http.Response, error) {
	r.URL.Host = host
	r.Host = host
	r.URL.Scheme = "http"
	r.RequestURI = ""

	r.Header.Add("X-Forwarded-For", r.URL.Hostname())

	resp, err := cl.Do(r)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
