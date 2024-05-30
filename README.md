# Golang Developer Assigment

## How to run this

docker build --tag go-exercise .
docker run -p 8080:8080 go-exercise

go test (to run the tests)

## One pair
![image](https://github.com/LucasRosello/go-exercise/assets/55340118/0308689b-c9ee-4a9b-99c8-dde7c313f878)

## Multiple Pairs

![image](https://github.com/LucasRosello/go-exercise/assets/55340118/acb3485c-15c8-49c3-abf2-58820fe2750c)

## Exercise
Develop in Go language a service that will provide an API for retrieval of the Last Traded Price of Bitcoin for the following currency pairs:

1. BTC/USD
2. BTC/CHF
3. BTC/EUR


The request path is:
/api/v1/ltp

The response shall constitute JSON of the following structure:
```json
{
  "ltp": [
    {
      "pair": "BTC/CHF",
      "amount": 49000.12
    },
    {
      "pair": "BTC/EUR",
      "amount": 50000.12
    },
    {
      "pair": "BTC/USD",
      "amount": 52000.12
    }
  ]
}

```

# Requirements:

1. The incoming request can done for as for a single pair as well for a list of them
2. You shall provide time accuracy of the data up to the last minute.
3. Code shall be hosted in a remote public repository
4. readme.md includes clear steps to build and run the app
5. Integration tests
6. Dockerized application

# Docs
The public Kraken API might be used to retrieve the above LTP information
[API Documentation](https://docs.kraken.com/rest/#tag/Spot-Market-Data/operation/getTickerInformation)
(The values of the last traded price is called “last trade closed”)
