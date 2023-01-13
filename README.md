# briq-exporter

a prometheus exporter for briqs

## how to build

```
go build
```

## how to use

- set the `BRIQ_SECRET_KEY` environment variable to your personal briq token
- run the exporter
- setup your prometheus system to scrape the exporter
