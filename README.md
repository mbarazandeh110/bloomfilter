Uses redisbloom (cuckoo filters) instead of bloomfilter.

## Implementations
- `krakend`: Integration of the `redisbloom` package as a rejecter for KrakenD


## Configuration
set the environment variable REDIS_PASSWORD="YOUR_REDIS_PASS"
and add the extra config at your backend:

```
"github_com/devopsfaith/bloomfilter": {
  "HashName": "hash-name",
  "Address": "ip:port",
  "TokenKeys": ["tokenkey1"]
},
```
