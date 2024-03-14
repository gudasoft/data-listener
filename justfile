default:
  @just --choose

start:
   go run cmd/main.go

benchmark:
  echo "Benchmark signle threaded"

  go-wrk  -c=450 -n=40000 -m="POST" -p=benchmarks/json/data1.json  -i http://127.0.0.1:8080/
  echo "Benchmark with 8 threads"
  go-wrk  -c=450 -n=40000 -t=8 -m="POST" -p=benchmarks/json/data1.json  -i http://127.0.0.1:8080/
