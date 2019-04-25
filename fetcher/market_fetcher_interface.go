package fetcher

import (
	mFetcher "github.com/chaintex/server-api/fetcher/market-fetcher"
	"github.com/chaintex/server-api/tomochain"
)

type MarketFetcherInterface interface {
	GetRateUsdTomo() (string, error)
	GetGeneralInfo(string) (*tomochain.TokenGeneralInfo, error)
}

func NewMarketFetcherInterface() MarketFetcherInterface {
	return mFetcher.NewCGFetcher()
}
