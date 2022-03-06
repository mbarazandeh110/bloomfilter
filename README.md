Uses redisbloom (cuckoo filters) instead of bloomfilter.

## Implementations
- `krakend`: Integration of the `redisbloom` package as a rejecter for KrakenD


## Configuration
set the environment variable REDIS_PASSWORD="YOUR_REDIS_PASS" and REDIS_ADDRESS="ip1:port1,ip2,port2"
and add the extra config at your backend:

```
"github_com/devopsfaith/bloomfilter": {
  "HashName": "hash-name",
  "TokenKeys": ["tokenkey1"]
},
```
