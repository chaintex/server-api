package persister

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chaintex/server-api/tomochain"
)

const (
	STEP_SAVE_RATE      = 10 //1 minute
	MAXIMUM_SAVE_RECORD = 60 //60 records

	INTERVAL_UPDATE_CHAINTEX_ENABLE    = 20
	INTERVAL_UPDATE_MAX_GAS            = 70
	INTERVAL_UPDATE_GAS                = 40
	INTERVAL_UPDATE_RATE_USD           = 610
	INTERVAL_UPDATE_GENERAL_TOKEN_INFO = 3600
	INTERVAL_UPDATE_GET_BLOCKNUM       = 20
	INTERVAL_UPDATE_GET_RATE           = 30
	INTERVAL_UPDATE_DATA_TRACKER       = 310
)

type RamPersister struct {
	mu      sync.RWMutex
	timeRun string

	chaintexEnabled       bool
	isNewChainTextEnabled bool

	rates     []tomochain.Rate
	isNewRate bool
	updatedAt int64

	latestBlock      string
	isNewLatestBlock bool

	rateUSD      []RateUSD
	rateTOMO     string
	isNewRateUsd bool

	// rateUSDCG      []RateUSD
	// rateTOMOCG      string
	// isNewRateUsdCG bool

	events     []tomochain.EventHistory
	isNewEvent bool

	maxGasPrice      string
	isNewMaxGasPrice bool

	gasPrice      *tomochain.GasPrice
	isNewGasPrice bool

	tokenInfo map[string]*tomochain.TokenGeneralInfo

	marketInfo              map[string]*tomochain.MarketInfo
	last7D                  map[string][]float64
	rate24H                 []tomochain.RateUSD
	change24H               []tomochain.RateUSD
	changeRate24H           []tomochain.RateUSD
	isNewTrackerData        bool
	numRequestFailedTracker int

	rightMarketInfo map[string]*tomochain.RightMarketInfo

	isNewMarketInfo bool
}

func NewRamPersister() (*RamPersister, error) {
	var mu sync.RWMutex
	location, _ := time.LoadLocation("Asia/Bangkok")
	tNow := time.Now().In(location)
	timeRun := fmt.Sprintf("%02d:%02d:%02d %02d-%02d-%d", tNow.Hour(), tNow.Minute(), tNow.Second(), tNow.Day(), tNow.Month(), tNow.Year())

	chaintexEnabled := true
	isNewChainTextEnabled := true

	rates := []tomochain.Rate{}
	isNewRate := false

	latestBlock := "0"
	isNewLatestBlock := true

	rateUSD := make([]RateUSD, 0)
	rateTOMO := "0"
	isNewRateUsd := true

	events := make([]tomochain.EventHistory, 0)
	isNewEvent := true

	maxGasPrice := "50"
	isNewMaxGasPrice := true

	gasPrice := tomochain.GasPrice{}
	isNewGasPrice := true

	tokenInfo := map[string]*tomochain.TokenGeneralInfo{}

	marketInfo := map[string]*tomochain.MarketInfo{}
	last7D := map[string][]float64{}
	rate24H := []tomochain.RateUSD{}
	change24H := []tomochain.RateUSD{}
	changeRate24H := []tomochain.RateUSD{}
	isNewTrackerData := true

	rightMarketInfo := map[string]*tomochain.RightMarketInfo{}

	isNewMarketInfo := true

	persister := &RamPersister{
		mu:                      mu,
		timeRun:                 timeRun,
		chaintexEnabled:         chaintexEnabled,
		isNewChainTextEnabled:   isNewChainTextEnabled,
		rates:                   rates,
		isNewRate:               isNewRate,
		updatedAt:               0,
		latestBlock:             latestBlock,
		isNewLatestBlock:        isNewLatestBlock,
		rateUSD:                 rateUSD,
		rateTOMO:                rateTOMO,
		isNewRateUsd:            isNewRateUsd,
		events:                  events,
		isNewEvent:              isNewEvent,
		maxGasPrice:             maxGasPrice,
		isNewMaxGasPrice:        isNewMaxGasPrice,
		gasPrice:                &gasPrice,
		isNewGasPrice:           isNewGasPrice,
		tokenInfo:               tokenInfo,
		marketInfo:              marketInfo,
		last7D:                  last7D,
		rate24H:                 rate24H,
		change24H:               change24H,
		changeRate24H:           changeRate24H,
		isNewTrackerData:        isNewTrackerData,
		numRequestFailedTracker: 0,
		rightMarketInfo:         rightMarketInfo,
		isNewMarketInfo:         isNewMarketInfo,
	}
	return persister, nil
}

func (rPersister *RamPersister) SaveGeneralInfoTokens(generalInfo map[string]*tomochain.TokenGeneralInfo) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.tokenInfo = generalInfo
	// rPersister.tokenInfoCG = generalInfoCG
}

func (rPersister *RamPersister) GetTokenInfo() map[string]*tomochain.TokenGeneralInfo {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.tokenInfo
}

/////------------------------------
func (rPersister *RamPersister) GetRate() []tomochain.Rate {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.rates
}

func (rPersister *RamPersister) GetTimeUpdateRate() int64 {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.updatedAt
}

func (rPersister *RamPersister) SetIsNewRate(isNewRate bool) {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	// return rPersister.rates
	rPersister.isNewRate = isNewRate
}

func (rPersister *RamPersister) GetIsNewRate() bool {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.isNewRate
}

func (rPersister *RamPersister) SaveRate(rates []tomochain.Rate, timestamp int64) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.rates = rates
	if timestamp != 0 {
		rPersister.updatedAt = timestamp
	}
}

//SaveRateUsd24H func
func (rPersister *RamPersister) SaveRateUsd24H(rate24H []tomochain.RateUSD, timestamp int64) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.rate24H = rate24H
	if timestamp != 0 {
		rPersister.updatedAt = timestamp
	}
}

//SaveChangeUsd24H func
func (rPersister *RamPersister) SaveChangeUsd24H(change24H []tomochain.RateUSD, timestamp int64) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.change24H = change24H
	if timestamp != 0 {
		rPersister.updatedAt = timestamp
	}
}

//SaveChangeRate24H func
func (rPersister *RamPersister) SaveChangeRate24H(change24H []tomochain.RateUSD, timestamp int64) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.changeRate24H = change24H
	if timestamp != 0 {
		rPersister.updatedAt = timestamp
	}
}

//SaveChainTextEnabled func
func (rPersister *RamPersister) SaveChainTextEnabled(enabled bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.chaintexEnabled = enabled
	rPersister.isNewChainTextEnabled = true
}

func (rPersister *RamPersister) SetNewChainTextEnabled(isNew bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewChainTextEnabled = isNew
}

func (rPersister *RamPersister) GetChainTextEnabled() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.chaintexEnabled
}

func (rPersister *RamPersister) GetNewChainTextEnabled() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.isNewChainTextEnabled
}

//--------------------------------------------------------

//--------------------------------------------------------

func (rPersister *RamPersister) SetNewMaxGasPrice(isNew bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewMaxGasPrice = isNew
	return
}

func (rPersister *RamPersister) SaveMaxGasPrice(maxGasPrice string) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.maxGasPrice = maxGasPrice
	rPersister.isNewMaxGasPrice = true
	return
}
func (rPersister *RamPersister) GetMaxGasPrice() string {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.maxGasPrice
}
func (rPersister *RamPersister) GetNewMaxGasPrice() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.isNewMaxGasPrice
}

//--------------------------------------------------------

func (rPersister *RamPersister) GetRateUSD() []RateUSD {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.rateUSD
}

func (rPersister *RamPersister) GetRateTOMO() string {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.rateTOMO
}

func (rPersister *RamPersister) GetIsNewRateUSD() bool {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.isNewRateUsd
}

func (rPersister *RamPersister) SaveRateUSD(rateUSDEth string) error {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()

	rates := make([]RateUSD, 0)
	// ratesCG := make([]RateUSD, 0)

	itemRateEth := RateUSD{Symbol: "TOMO", PriceUsd: rateUSDEth}
	// itemRateEthCG := RateUSD{Symbol: "TOMO", PriceUsd: rateUSDEthCG}
	rates = append(rates, itemRateEth)
	// ratesCG = append(ratesCG, itemRateEthCG)
	for _, item := range rPersister.rates {
		if item.Source != "TOMO" {
			priceUsd, err := CalculateRateUSD(item.Rate, rateUSDEth)
			if err != nil {
				log.Print(err)
				rPersister.isNewRateUsd = false
				return nil
			}

			itemRate := RateUSD{Symbol: item.Source, PriceUsd: priceUsd}
			rates = append(rates, itemRate)
		}
	}

	rPersister.rateUSD = rates
	rPersister.rateTOMO = rateUSDEth
	rPersister.isNewRateUsd = true

	return nil
}

func CalculateRateUSD(rateTomo string, rateUSD string) (string, error) {
	bigRateUSD, ok := new(big.Float).SetString(rateUSD)
	if !ok {
		err := errors.New("Cannot convert rate usd of tomochain to big float")
		return "", err
	}
	bigRateEth, ok := new(big.Float).SetString(rateTomo)
	if !ok {
		err := errors.New("Cannot convert rate token-tomo to big float")
		return "", err
	}
	i, e := big.NewInt(10), big.NewInt(18)
	i.Exp(i, e, nil)
	weight := new(big.Float).SetInt(i)

	rateUSDBig := new(big.Float).Mul(bigRateUSD, bigRateEth)
	rateUSDNormal := new(big.Float).Quo(rateUSDBig, weight)
	return rateUSDNormal.String(), nil
}

func (rPersister *RamPersister) SetNewRateUSD(isNew bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewRateUsd = isNew
}

func (rPersister *RamPersister) GetLatestBlock() string {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.latestBlock
}

func (rPersister *RamPersister) SaveLatestBlock(blockNumber string) error {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.latestBlock = blockNumber
	rPersister.isNewLatestBlock = true
	return nil
}

func (rPersister *RamPersister) GetIsNewLatestBlock() bool {
	rPersister.mu.RLock()
	defer rPersister.mu.RUnlock()
	return rPersister.isNewLatestBlock
}

func (rPersister *RamPersister) SetNewLatestBlock(isNew bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewLatestBlock = isNew
}

// ----------------------------------------
// return data from chaintex tracker

// use this api for 3 infomations change, marketcap, volume
func (rPersister *RamPersister) GetRightMarketData() map[string]*tomochain.RightMarketInfo {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.rightMarketInfo
}

func (rPersister *RamPersister) GetIsNewTrackerData() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.isNewTrackerData
}

func (rPersister *RamPersister) SetIsNewTrackerData(isNewTrackerData bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewTrackerData = isNewTrackerData
	rPersister.numRequestFailedTracker = 0
}

func (rPersister *RamPersister) GetLast7D(listTokens string) map[string][]float64 {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	tokens := strings.Split(listTokens, "-")
	result := make(map[string][]float64)
	for _, symbol := range tokens {
		if rPersister.last7D[symbol] != nil {
			result[symbol] = rPersister.last7D[symbol]
		}
	}
	return result
}

// GetRate24H func return price of tokens in around 24h
func (rPersister *RamPersister) GetRate24H(listTokens string) map[string]float64 {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	tokens := strings.Split(listTokens, "-")
	result := make(map[string]float64)
	for _, symbol := range tokens {
		for _, r := range rPersister.rate24H {
			if r.Symbol == symbol {
				price, err := strconv.ParseFloat(r.PriceUsd, 64)
				if err == nil {
					result[symbol] = price
				}
				break
			}
		}
	}
	return result
}

// GetChange24H func return price of tokens in around 24h
func (rPersister *RamPersister) GetChange24H(typ string, listTokens string) map[string]float64 {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	tokens := strings.Split(listTokens, "-")
	var change24H []tomochain.RateUSD

	if typ == "usd" {
		change24H = rPersister.change24H
	} else {
		change24H = rPersister.changeRate24H
	}

	result := make(map[string]float64)
	for _, symbol := range tokens {
		for _, r := range change24H {
			if r.Symbol == symbol {
				price, err := strconv.ParseFloat(r.PriceUsd, 64)
				if err == nil {
					result[symbol] = price
				}
				break
			}
		}
	}
	return result
}

func (rPersister *RamPersister) SaveMarketData(marketRate map[string]*tomochain.Rates, mapTokenInfo map[string]*tomochain.TokenGeneralInfo, tokens map[string]tomochain.Token) {
	lastSevenDays := map[string][]float64{}
	newResult := map[string]*tomochain.RightMarketInfo{}
	if len(mapTokenInfo) == 0 {
		rPersister.mu.RLock()
		mapTokenInfo = rPersister.tokenInfo
		rPersister.mu.RUnlock()
	}
	for symbol := range tokens {
		dataSevenDays := []float64{}
		rightMarketInfo := &tomochain.RightMarketInfo{}
		rateInfo := marketRate[symbol]
		if rateInfo != nil {
			dataSevenDays = rateInfo.P
			rightMarketInfo.Rate = &rateInfo.R
		}
		if tokenInfo := mapTokenInfo[symbol]; tokenInfo != nil {
			rightMarketInfo.Quotes = tokenInfo.Quotes
			rightMarketInfo.Change24H = tokenInfo.Change24H
		}

		if rateInfo == nil && rightMarketInfo.Quotes == nil {
			continue
		}

		newResult[symbol] = rightMarketInfo
		lastSevenDays[symbol] = dataSevenDays
	}

	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.last7D = lastSevenDays
	rPersister.rightMarketInfo = newResult
}

func (rPersister *RamPersister) SetIsNewMarketInfo(isNewMarketInfo bool) {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.isNewMarketInfo = isNewMarketInfo
}

func (rPersister *RamPersister) GetIsNewMarketInfo() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.isNewMarketInfo
}

func (rPersister *RamPersister) GetTimeVersion() string {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	return rPersister.timeRun
}

func (rPersister *RamPersister) IsFailedToFetchTracker() bool {
	rPersister.mu.Lock()
	defer rPersister.mu.Unlock()
	rPersister.numRequestFailedTracker++
	if rPersister.numRequestFailedTracker > 12 {
		return true
	}
	return false
}
