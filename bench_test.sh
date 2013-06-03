#go test -run=none -file bench_test.go -test.bench . -cpuprofile=bench_test.out

#go test -c
#./go-restful.test -test.run=none -test.cpuprofile=tmp.prof -test.bench=BenchmarkMany
#go tool pprof go-restful.test tmp.prof