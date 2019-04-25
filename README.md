# Cached server for ChainTeX wallet

## Build

```
docker-compose build cached
```

## Run
docker-compose up -d cached
```


## One command to build and run with docker-compose

docker-compose -f docker-compose-staging.yml up --build
```

## APIs (these APIs will be expired after Jan 20 2019)
 - /getRateUSD: return USD price of token base on it's expectedRate
 - /getRate: return rate of token with tomo (expectedRate and minRate)
 - /getRateTOMO: return USD price of TOMO from Coingecko

## Cache version
 - /cacheVersion: return current cache version
 
### 1. Get Rate USD
`/rateUSD`

(GET) Return USD price of token base on it's expectedRate

Response:
```javascript
{
    "data": [
        {
            "symbol": "TOMO",
            "price_usd": "150.110255"
        }
    ],
    "success": true
}
```

### 2. Get rate
`/rate`

(GET) Return rate of token with tomo (expectedRate and minRate)

Response:
```javascript
{
    "data": [
        {
            "source": "POWR",
            "dest": "TOMO",
            "rate": "580350000000000",
            "minRate": "562939500000000"
        },
        {
            "source": "REQ",
            "dest": "TOMO",
            "rate": "251549999999999",
            "minRate": "244003499999999"
        }
    ],
    "success": true
    }
```