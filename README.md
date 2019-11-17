# Exchange

![Build](https://img.shields.io/docker/cloud/build/pmdcosta/exchange)

A golang application that provides the latest exchange rates for the base currencies of GBP and USD.
It uses https://exchangeratesapi.io/ for the exchange rates. It checks the latest value against the historic
rate for the last week and makes a naive recommendation as to whether this is a good time to exchange money or not.
