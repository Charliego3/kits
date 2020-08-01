package request

import (
	"github.com/whimthen/kits/logger"
	"net/url"
	"testing"
)

func TestHttpRequest_DoGet(t *testing.T) {
	t.Log("Start test doGet......")
	request := &HttpRequest{
		Url: "http://trade.100-123.net/api/queueOrder",
	}

	result, err := request.AddHeader("Content-Type", "application/json;charset=utf8").DoFromPost(url.Values{
		"accesskey": {"fde3310f-8266-4abf-9639-011f0e32d315"},
		"acctType":  {"0"},
		"amount":    {"0.1"},
		"currency":  {"zbusdt"},
		"method":    {"order"},
		"orderType": {"0"},
		"price":     {"6.5"},
		"tradeType": {"1"},
		"sign":      {"f30f72046213f0365ddbee832a582f3b"},
		"reqTime":   {"1594889312000"},
	}).ResponseString()

	if err != nil {
		logger.Error("HttpRequest error, %+v", err)
		return
	}

	logger.Info("Result: %s", result)
	t.Log("Complete......")
}
