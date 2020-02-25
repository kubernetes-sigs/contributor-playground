package temp_cce

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/util"
)

var (
	testHTTPServer *httptest.Server
	cceClient      *Client

	logger util.LoggerItf = util.DefaultLogger

	userID string = strings.Replace(util.GetRequestID(), "-", "", -1)
)

type testEnvConfig struct {
	uri          string
	method       string
	statusCode   int
	responseBody []byte
}

type handler func(w http.ResponseWriter, r *http.Request)

func setupTestEnv(configs []*testEnvConfig) {
	credentials := &bce.Credentials{
		AccessKeyID:     strings.Replace(util.GetRequestID(), "-", "", -1),
		SecretAccessKey: strings.Replace(util.GetRequestID(), "-", "", -1),
	}

	var bceConfig = &bce.Config{
		Credentials: credentials,
		Checksum:    true,
	}
	var cceConfig = NewConfig(bceConfig)
	cceClient = NewClient(cceConfig)
	cceClient.SetDebug(true)

	r := mux.NewRouter()
	for _, config := range configs {
		handler := newHandler(config.statusCode, config.responseBody)
		r.HandleFunc(config.uri, handler).Methods(config.method)
	}

	testHTTPServer = httptest.NewServer(r)
	cceClient.Endpoint = testHTTPServer.URL
}

func tearDownTestEnv() {
	testHTTPServer.Close()
}

func newHandler(statusCode int, responseBody []byte) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write(responseBody)
	}
}
