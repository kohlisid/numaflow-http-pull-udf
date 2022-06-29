# numaflow-http-pull-udf
UDF to periodically pull from an http endpoint to generate message.

Currently we use this to fetch stock information from a given API.

Usage:- 

Add the UDF in the pipeline YAML file

in the -tickers arg, provide a comma seperated list of stocks that you wish to track!


- name: sample-udf
  udf:
    container:
      image: quay.io/kohlisid/http-pull-numa:fin1
      args:
      - /http-pull-udf
      - -tickers=AAPL,AMZN

We would need to define the edges with conditional forwarding for each stock ticker, the keyIn should be the ticker name

edges:
  - from: in
    to: sample-udf
  - from: sample-udf
    to: apple
    conditions:
      keyIn:
      - AAPL
