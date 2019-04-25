package persister

import "github.com/chaintex/server-api/tomochain"

type BoltInterface interface {
	StoreGeneralInfo(map[string]*tomochain.TokenGeneralInfo) error
	GetGeneralInfo(map[string]tomochain.Token) (map[string]*tomochain.TokenGeneralInfo, error)
}
