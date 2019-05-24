package persister

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/chaintex/server-api/tomochain"
	_ "github.com/lib/pq"
)

const (
	host   = "postgres-db"
	port   = "5432"
	db     = "chaintex_db"
	usr    = "chaintext"
	pwd    = "chaintext"
	tbname = "price_token"
)

// PostgresDb storage for cache
type PostgresDb struct {
	psqlInfo string
}

// NewPostgresDb make db instance
func NewPostgresDb() (*PostgresDb, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, usr, pwd, db)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return &PostgresDb{
		psqlInfo: psqlInfo,
	}, nil
}

// StoreRateInfo store market info
func (bs *PostgresDb) StoreRateInfo(mapInfo []tomochain.RateUSD) error {
	timeNow := time.Now().UTC()
	cmd := fmt.Sprintf("INSERT INTO %s (time_at, symbol, price) VALUES ($1, $2, $3)", tbname)
	dbIns, err := sql.Open("postgres", bs.psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbIns.Close()

	for _, r := range mapInfo {
		price, ex := strconv.ParseFloat(r.PriceUsd, 64)
		if ex != nil {
			continue
		}
		insertData(dbIns, cmd, timeNow, r.Symbol, price)
	}

	return nil
}

// GetChange24H store market info
func (bs *PostgresDb) GetChange24H(symbols []string) ([]tomochain.RateUSD, error) {
	var rates []tomochain.RateUSD
	wSymbol := strings.Join(symbols, "','")
	sqlStatement := fmt.Sprintf("SELECT symbol, AVG(price) AS price FROM %s WHERE symbol IN ('%s') AND price > 0 GROUP BY symbol", tbname, wSymbol)

	dbIns, err := sql.Open("postgres", bs.psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbIns.Close()

	rows, err := dbIns.Query(sqlStatement)
	if err != nil {
		return rates, err
	}
	defer rows.Close()

	if err != nil {
		return rates, err
	}
	defer rows.Close()
	for rows.Next() {
		var symbol string
		var price float64
		err = rows.Scan(&symbol, &price)
		if err != nil {
			// handle this error
			panic(err)
		}
		rateUSDItem := tomochain.RateUSD{
			Symbol:   symbol,
			PriceUsd: fmt.Sprint(price),
		}

		rates = append(rates, rateUSDItem)
	}

	err = rows.Err()
	return rates, err
}

// GetRate24H store market info
func (bs *PostgresDb) GetRate24H(symbols []string) ([]tomochain.RateUSD, error) {
	var rates []tomochain.RateUSD
	wSymbol := strings.Join(symbols, "','")
	sqlStatement := fmt.Sprintf("SELECT symbol, AVG(price) AS price FROM %s WHERE symbol IN ('%s') AND price > 0 GROUP BY symbol", tbname, wSymbol)

	dbIns, err := sql.Open("postgres", bs.psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbIns.Close()

	rows, err := dbIns.Query(sqlStatement)
	if err != nil {
		return rates, err
	}
	defer rows.Close()
	for rows.Next() {
		var symbol string
		var price float64
		err = rows.Scan(&symbol, &price)
		if err != nil {
			// handle this error
			panic(err)
		}

		rateUSDItem := tomochain.RateUSD{
			Symbol:   symbol,
			PriceUsd: fmt.Sprint(price),
		}

		rates = append(rates, rateUSDItem)
	}

	err = rows.Err()
	return rates, err
}

func insertData(dbIns *sql.DB, sqlStatement string, args ...interface{}) {
	result, err := dbIns.Exec(sqlStatement, args...)
	if err != nil {
		log.Println("++++++++++++++insertData error: ", err)
	}

	log.Println("++++++++++++++insertData data: ", result)
}
