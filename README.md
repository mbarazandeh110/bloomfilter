Uses redisbloom (cuckoo filters) instead of bloomfilter.

## Implementations
- `krakend`: Integration of the `redisbloom` package as a rejecter for KrakenD


## Configuration

Just add the extra config at your backend:

```
"github_com/devopsfaith/bloomfilter": {
  "HashName": "hash-name",
  "Password": "passowrd",
  "Address": "ip:port",
  "TokenKeys": ["tokenkey1"]
},
```