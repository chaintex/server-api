package persister

import "github.com/chaintex/server-api/tomochain"

// PostgresInterface type
type PostgresInterface interface {
	StoreRateInfo([]tomochain.RateUSD) error
	GetRate24H(symbols []string) ([]tomochain.RateUSD, error)
	GetChange24H(typ string, symbols []string) ([]tomochain.RateUSD, error)
}
