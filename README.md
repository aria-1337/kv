# kv
Simple GO K/V database powered by leveldb.

HTTP requests are made via json.
```json
{
    "key":"string",
    "value":"any<interface{}>"
}
```

```bash
POST example
~> curl -X POST -d '{"key":"coolKeyName", "value": "myKeysAwesomeValue"}' localhost:3000
#> 201 created | 403 Forbidden (Post can not overwrite)

GET example
~> curl -X GET -d '{"key": "coolKeyName"}' localhost:3000
#>"myKeysAwesomeValue"
#> 200 success | 404 Record doesn't exist

DELETE example
~> curl -X DELETE -d '{"key": "coolKeyName"}' localhost:3000
#> 204 success | 409 Conflict (key doesn't exist)

PATCH example
~> curl -X PATCH -d '{"key": "existingKey", "value": "newValue"}' localhost:3000
#> 200 success | 404 Record doesn't exist
```

It does:
- GET/SET/DELETE/PATCH

TODO:
- Multiple volumes
- Distrubution
- Docker build
