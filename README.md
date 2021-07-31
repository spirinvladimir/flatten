# Flattening service for S-expressions

There are two POST routes:
 - `/flatten` take a JSON array of S-expressions and response it's flatten representation with a depth.
 - `/history` response statistic for last *100* requests.

There is a memory cache for avoid same calculation again. It doesn't support horizontal scaling. To make it happen there are solutions:
 - replace `cache` variable with Redis like key-value service
 - heat `cache` to make `cache` same between pods
History also doesn't scale horizontally in favour of simplicity.

## Running
```sh
go run server.go
```
Default port is `8080`.

## Usage (example)
Let's calculate S-expressions depth of *Selection Sort* in native *WebAssembly* is *10*.
 ```sh
curl -s https://raw.githubusercontent.com/spirinvladimir/selection-sort/master/main.wat -o ss.wat;
echo '["' + $(cat ss.wat | tr -d "\"") + '"]' > ss.json;
curl -s -d @ss.json -X POST http://localhost:8080/flatten | jq '.[0] .depth ';
10

```

## Test
```sh
go test ./...
```
