package sink

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/donkeywon/gtil/httpc"
	"github.com/donkeywon/gtil/util"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	HttpSinkName = "http"
)

func init() {
	err := zap.RegisterSink(HttpSinkName, httpSink)
	if err != nil {
		fmt.Println("Register http sink fail", err)
	}
}

func httpSink(url *url.URL) (zap.Sink, error) {
	return &Http{
		httpc: httpc.New(httpc.NewConfig(), context.Background()),
		url:   url,
	}, nil
}

type Http struct {
	httpc *httpc.HttpC
	url   *url.URL
}

func (h *Http) Write(p []byte) (n int, err error) {
	req, err := http.NewRequest("POST", h.url.String(), bytes.NewReader(p))
	if err != nil {
		return 0, err
	}

	resp, err := h.httpc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return 0, errors.New(util.Bytes2String(respBody))
	}

	return len(p), nil
}

func (h *Http) Sync() error {
	return nil
}

func (h *Http) Close() error {
	return h.httpc.Close()
}
