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
	usr    = "chaintex"
	pwd    = "chaintex"
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
	cmd := fmt.Sprintf("INSERT INTO %s (time_at, symbol, price, rate_with_tomo) VALUES ($1, $2, $3, $4)", tbname)
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

		rateWithTomo, ex2 := strconv.ParseFloat(r.RateWithTomo, 64)

		if ex2 != nil {
			continue
		}

		insertData(dbIns, cmd, timeNow, r.Symbol, price, rateWithTomo)
	}

	return nil
}

// GetChange24H store market info
func (bs *PostgresDb) GetChange24H(typ string, symbols []string) ([]tomochain.RateUSD, error) {
	var rates []tomochain.RateUSD

	col := "price"

	if typ == "usd" {
		col = "price"
	} else {
		col = "rate_with_tomo"
	}

	dbIns, err := sql.Open("postgres", bs.psqlInfo)
	if err != nil {
		panic(err)
	}
	defer dbIns.Close()
	for _, symbol := range symbols {
		sqlStatement := "select COALESCE((select " + col + " from " + tbname + " where symbol = '" + symbol + "' AND " + col + " > 0 and time_at <= now() order by time_at desc limit 1), 0) - COALESCE((select " + col + " from " + tbname + " where symbol = '" + symbol + "' AND " + col + " > 0 and time_at <= (NOW() - INTERVAL '24' HOUR) order by time_at desc limit 1), 0) as change"
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
			var change float64
			err = rows.Scan(&change)
			if err != nil {
				// handle this error
				panic(err)
			}

			rateUSDItem := tomochain.RateUSD{
				Symbol:   symbol,
				PriceUsd: fmt.Sprint(change),
			}

			rates = append(rates, rateUSDItem)
		}

		err = rows.Err()
	}

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
