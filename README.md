# gobench2json
json otput for go -bench

    go test -bench=. -run=^$ -benchmem ./... | ./gobenchtojson
