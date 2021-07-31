# Flattening service

There are two POST routes:
 - /flatten
 - /history

/flatten take JSON array of S-expressions and response flatten representation with it's depth.
/history response statistic for last 100 requests.

Features:
 - Memory cache for avoid same calculation again. It does not support horizontal scale. To make it happen:
 - replace `cache` variable with Redis like key-value service
 - heat `cache` to make `cache` same between pods

 # Running
```sh
go run server.go
```


 # Usage
Depth of Selection Sort in native WebAssembly is 10.
 ```sh
curl -s https://raw.githubusercontent.com/spirinvladimir/selection-sort/master/main.wat -o ss.wat;
sed -i 's/"//g' ss.wat;
echo '["' + $(cat ss.wat) + '"]' > ss.json;
curl -s -d @ss.json -X POST http://localhost:8080/flatten | jq '.[0] .depth ';
 10
```

# Test
```sh
go test ./...
```
