package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	httptransport "github.com/go-kit/kit/transport/http"
)

func NewJsonClient(method string, target *url.URL, decode httptransport.DecodeResponseFunc, authorization string, requestFunc httptransport.RequestFunc) *httptransport.Client {
	return httptransport.NewClient(
		method,
		target,
		httptransport.EncodeJSONRequest,
		decode,
		httptransport.ClientBefore(httptransport.SetRequestHeader("Authorization", authorization)),
	)
}

func DefaultHttpRequestEncoder(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	r.Close = true
	return nil
}

func QueryParamBeforeFunc(queryParamMap map[string]string) httptransport.RequestFunc {
	f := func(ctx context.Context, r *http.Request) context.Context {
		q := r.URL.Query()
		for key, value := range queryParamMap {
			q.Add(key, value)
		}
		r.URL.RawQuery = q.Encode()
		return ctx
	}

	return f
}
