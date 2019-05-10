package persister

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chaintex/server-api/tomochain"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

const (
	addr      = "http://influx-db:8086"
	database  = "chaintex_db"
	username  = "chaintext"
	password  = "chaintext"
	tablename = "price_token"
)

// InfluxDb storage for cache
type InfluxDb struct {
	dbIns client.Client
}

// NewInfluxDb make bolt instance
func NewInfluxDb() (*InfluxDb, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	return &InfluxDb{
		dbIns: c,
	}, nil
}

// StoreRateInfo store market info
func (bs *InfluxDb) StoreRateInfo(mapInfo []tomochain.RateUSD) error {
	log.Println("+++++++++++++++++StoreRateInfo", mapInfo)
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "s",
	})

	if err != nil {
		log.Println("Error creating a NewBatchPoint", err)
	}

	for _, r := range mapInfo {
		log.Println("+++++++++++++++++start item ", r)
		price, ex := strconv.ParseFloat(r.PriceUsd, 64)
		if ex != nil {
			continue
		}

		tags := map[string]string{
			"symbol": r.Symbol,
		}

		fields := map[string]interface{}{
			"price": price,
		}

		pt, err := client.NewPoint(
			tablename,
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Println("Error creating NewPoint ", err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if errW := bs.dbIns.Write(bp); errW != nil {
			log.Println("Write the batch to influx DB: ", errW)
		}

		log.Println("+++++++++++++++++end item ", r)
	}

	return nil
}

// GetRate24H store market info
func (bs *InfluxDb) GetRate24H(symbols []string) ([]tomochain.RateUSD, error) {
	var rates []tomochain.RateUSD
	wSymbol := strings.Join(symbols, "|")
	cmd := fmt.Sprintf("SELECT MEAN(price) FROM %s WHERE (symbol =~ /%s/) AND  (time >= now() - 24h) GROUP BY symbol", tablename, wSymbol)

	results, err := getQuery(bs.dbIns, cmd)
	if err != nil {
		return rates, nil
	}

	for _, r := range results {
		log.Println("++++++++++++++row: ", r.Tags["symbol"])
		row := r.Values[0][0]

		log.Println("++++++++++++++row: ", row.price)
	}

	log.Println("+++++++++++++++++++++rates: ", rates)
	return rates, nil
}

func getQuery(dbIns client.Client, q string) ([]models.Row, error) {
	var rows []models.Row
	var err error
	query := client.Query{
		Command:  q,
		Database: database,
	}

	if response, er := dbIns.Query(query); er == nil {
		if response.Error() != nil {
			log.Println("++++++++++++++Query error: ", response.Error())
			err = response.Error()
		} else {
			rows = response.Results[0].Series
		}
	}

	return rows, err
}
