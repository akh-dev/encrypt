# encryption server

## quick start
1) Checkout, check config folders in each service 

2) build both services with
```bash
go build github.com/akh-dev/encrypt/encryption-service
go build github.com/akh-dev/encrypt/storage-service
```

3) Run both services from command line

## usage
To store encrypted text on the server, follow this example:
```curl
curl -X POST -d '{"id":"my-1st-text","payload":"some very long text version one"}' -H "Content-Type:application/json" localhost:8080/store
```

sample result:
```json
{"status_code":0,"status_message":"Success","result":{"id":"my-1st-text","key":"JAvDBuhM8yB4iKymW3mHOO8JpQ7nDN/dg+mgebuSIRs="}}
```


To retrieve and decrypt stored text, follow this example:
```curl
curl -X GET -d '{"id":"my-1st-text","key":"JAvDBuhM8yB4iKymW3mHOO8JpQ7nDN/dg+mgebuSIRs="}' -H "Content-Type:application/json" localhost:8080/retrieve
```

sample result:
```json
{"status_code":0,"status_message":"Success","result":{"id":"my-1st-text","payload":"some very long text version one"}}
```
