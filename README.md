# kv

Simple GO K/V database powered by leveldb.

# API
```bash
localhost:${PORT}/${KEY_NAME}

Body: JSON format { "value": "any<interface{}>" }
```

### POST
```bash
curl -X POST -d '{ "value": "[1, 2, 3]" }' localhost:3000/key
# 201 OK | 403 FORBIDDEN (overwrites are not allowed)
```

### GET
```bash
curl -X GET localhost:3000/key 
> [1,2,3]
# 200 OK | 404 NOT FOUND (key doesn't exist)
```

### PATCH
```bash
curl -X PATCH -d '{ "value": "[2, 3, 4]" }' localhost:3000/key
# 201 CREATED | 404 NOT FOUND
```

### DELETE
```bash
curl -X DELETE localhost:3000/key
# 200 OK | 409 CONFLICT (record doesn't exist) 
```


### Docker
TODO: upload image to dockerhub in the meantime...

```
docker build -t kv .
docker run --env PORT=3001 --env LEVEL_DB_PATH=test -P -t kv:latest
```
