# gobench2json
json otput for go -bench

    go test -bench=. -run=^$ -benchmem ./... | ./gobench2json


```
  "github.com/caasmo/restinpieces/crypto": [
    {
      "name": "BenchmarkNewJwtSigningKey-8",
      "runs": 722200,
      "ns_per_op": 1415,
      "b_per_op": 657,
      "allocs_per_op": 10
    },
    {
      "name": "BenchmarkNewJwt-8",
      "runs": 313455,
      "ns_per_op": 4252,
      "b_per_op": 2184,
      "allocs_per_op": 35
    },
    {
      "name": "BenchmarkParseJwt_Valid-8",
      "runs": 226909,
      "ns_per_op": 5228,
      "b_per_op": 2584,
      "allocs_per_op": 51
    },
    {
      "name": "BenchmarkParseJwt_InvalidSignature-8",
      "runs": 169788,
      "ns_per_op": 5990,
      "b_per_op": 2657,
      "allocs_per_op": 55
    },
    {
      "name": "BenchmarkParseJwtUnverified-8",
      "runs": 346656,
      "ns_per_op": 3217,
      "b_per_op": 1624,
      "allocs_per_op": 38
    }
  ]
}
```
