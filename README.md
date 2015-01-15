# Go implementation of popcount with SSSE3 optimizations.

Number are based on an Intel Core 2 Duo 2.26ghz:

````
$ go test -bench=".*"
PASS
Benchmark_PopCount16    100000000               13.3 ns/op
Benchmark_PopCount20    100000000               16.1 ns/op
Benchmark_PopCount250   20000000                68.9 ns/op
Benchmark_PopCountData16        50000000                33.5 ns/op
Benchmark_PopCountData20        30000000                41.8 ns/op
Benchmark_PopCountData250        3000000               433 ns/op
...
Benchmark_PopCount32    500000000                3.46 ns/op
Benchmark_PopCount64    500000000                3.75 ns/op
```

In the above, `PopCount32` and `PopCount64` are bithacks-based native go implementations.  `PopCountData` processes
`[]byte` slices using `PopCount32` and `PopCount64`.  `PopCount` uses SSSE3, if such instructions are
supported, to process data 128-bits at a time.

Performance comparisons are shown for 16, 20, and 250 byte slices in the first 6 lines.  Speedup ranges from 2.5x to 6x
in this instance.


