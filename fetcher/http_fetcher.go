package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/chaintex/server-api/common"
	fCommon "github.com/chaintex/server-api/fetcher/fetcher-common"
)

const (
	TIME_TO_DELETE  = 18000
	API_KEY_TRACKER = "jHGlaMKcGn5cCBxQCGwusS4VcnH0C6tN"
)

type HTTPFetcher struct {
	apiEndpoint string
}

func NewHTTPFetcher(apiEndpoint string) *HTTPFetcher {
	return &HTTPFetcher{
		apiEndpoint: apiEndpoint,
	}
}

func (httpFetcher *HTTPFetcher) GetUserInfo(url string) (*common.UserInfo, error) {
	userInfo := &common.UserInfo{}
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	err = json.Unmarshal(b, userInfo)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return userInfo, nil
}

//TokenPrices
type TokenPrice struct {
	Data []struct {
		Symbol string  `json:"symbol"`
		Price  float64 `json:"price"`
	} `json:"data"`
	Error      bool   `json:"error"`
	TimeUpdate uint64 `json:"timeUpdated"`
}

// GetRateUsdTomo get usd from api
func (httpFetcher *HTTPFetcher) GetRateUsdTomo() (string, error) {
	var tomoPrice string
	url := fmt.Sprintf("%s/token_price?currency=USD", httpFetcher.apiEndpoint)
	b, err := fCommon.HTTPCall(url)
	if err != nil {
		log.Print(err)
		return tomoPrice, err
	}
	var tokenPrice TokenPrice
	err = json.Unmarshal(b, &tokenPrice)
	if err != nil {
		log.Println(err)
		return tomoPrice, err
	}
	if tokenPrice.Error {
		return tomoPrice, errors.New("cannot get token price from api")
	}
	for _, v := range tokenPrice.Data {
		if v.Symbol == common.TOMOSymbol {
			tomoPrice = fmt.Sprintf("%.6f", v.Price)
			break
		}
	}
	return tomoPrice, nil
}
