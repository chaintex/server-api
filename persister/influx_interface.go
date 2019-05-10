package persister

import "github.com/chaintex/server-api/tomochain"

// InfluxInterface type
type InfluxInterface interface {
	StoreRateInfo([]tomochain.RateUSD) error
	GetRate24H(symbols []string) ([]tomochain.RateUSD, error)
}
