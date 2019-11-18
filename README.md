# Exchange

![Build](https://img.shields.io/docker/cloud/build/pmdcosta/exchange)

A golang application that provides the latest exchange rates and recommendations for currency conversion based on
https://exchangeratesapi.io/ data.
 
It compares the latest exchange value against the historic rate for the last week and makes a naive recommendation
as to whether this is a good time to exchange money or not.


## Endpoints

The rates endpoint returns the current exchange rate in `euros` for the given currency.
It also returns whether now is a good time to exchange money or not.

`GET http://{host}/v1/rates/{currency}`
```json
{
  "rate": {
    "EUR": 0.9062896502,
    "USD": 1
  },
  "exchange": {
    "recommended": false,
    "reason": "the current exchange value is below the weekly average of 0.907558"
  }
}
```


## Deployment

The repository in integrating with [hub.docker.com](https://hub.docker.com/repository/docker/pmdcosta/exchange). 
Updates to master trigger docker image releases with tag `latest` and git tags trigger image build with the corresponding tag.   

The docker image `pmdcosta/exchange` can be used to easily deploy the service.

Additionally, the following `Make` commands are available:
- `make test` runs unit tests.
- `make build` builds the docker image.
- `make start` starts the server locally.
- `make stop` stops the server.
- `make integration-tests` starts the server in a container and runs integration tests against it.


## Possible Improvements

- Instead of keeping the cache in memory, a Redis instance could be used.
Adding a new cache backend is fairly trivial, just requires injecting another 
client that satisfies the interface in the exchange client.
- The current cache stores the mapping between requests and responses at the http level in the exchangeratesapi.io client.
A different strategy would be to store the exchange rates for each day and evaluate if the http requests is needed at all.
- The exchangeratesapi.io client is quite flexible, but the current http server is very specific. More parameters could be added to
the http handler, or more endpoints created.

