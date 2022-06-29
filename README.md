# numaflow-http-pull-udf
UDF to periodically pull from an http endpoint to generate message.

Currently we use this to fetch stock information from a given API.

Usage:- 

Add the UDF in the pipeline YAML file

in the -tickers arg, provide a comma seperated list of stocks that you wish to track!

```
- name: sample-udf
  udf:
    container:
      image: quay.io/kohlisid/http-pull-numa:fin1
      args:
      - /http-pull-udf
      - -tickers=AAPL,AMZN
```

We would need to define the edges with conditional forwarding for each stock ticker, the keyIn should be the ticker name

```
edges:
  - from: in
    to: sample-udf
  - from: sample-udf
    to: apple
    conditions:
      keyIn:
      - AAPL
```

The output vertices will be created for each stock, displaying messages like the following:-

```
2022/06/29 16:30:12 (intuit) {"quoteResponse":{"result":[{"language":"en-US","region":"US","quoteType":"EQUITY","triggerable":true,"quoteSourceName":"Nasdaq Real Time Price","currency":"USD","marketState":"REGULAR","regularMarketPrice":386.65,"regularMarketTime":1656520169,"fullExchangeName":"NasdaqGS","market":"us_market","exchangeDataDelayedBy":0,"exchange":"NMS","sourceInterval":15,"exchangeTimezoneName":"America/New_York","exchangeTimezoneShortName":"EDT","gmtOffSetMilliseconds":-14400000,"esgPopulated":false,"tradeable":false,"priceHint":2,"symbol":"INTU"}],"error":null}}
2022/06/29 16:30:13 (intuit) {"quoteResponse":{"result":[{"language":"en-US","region":"US","quoteType":"EQUITY","triggerable":true,"quoteSourceName":"Nasdaq Real Time Price","currency":"USD","fullExchangeName":"NasdaqGS","market":"us_market","exchange":"NMS","regularMarketPrice":386.78,"regularMarketTime":1656520211,"marketState":"REGULAR","priceHint":2,"esgPopulated":false,"tradeable":false,"exchangeDataDelayedBy":0,"sourceInterval":15,"exchangeTimezoneName":"America/New_York","exchangeTimezoneShortName":"EDT","gmtOffSetMilliseconds":-14400000,"symbol":"INTU"}],"error":null}}
2022/06/29 16:30:13 (intuit) {"quoteResponse":{"result":[{"language":"en-US","region":"US","quoteType":"EQUITY","triggerable":true,"quoteSourceName":"Nasdaq Real Time Price","currency":"USD","exchange":"NMS","priceHint":2,"marketState":"REGULAR","fullExchangeName":"NasdaqGS","sourceInterval":15,"exchangeTimezoneName":"America/New_York","exchangeTimezoneShortName":"EDT","gmtOffSetMilliseconds":-14400000,"market":"us_market","regularMarketPrice":386.78,"regularMarketTime":1656520211,"esgPopulated":false,"tradeable":false,"exchangeDataDelayedBy":0,"symbol":"INTU"}],"error":null}}
```
