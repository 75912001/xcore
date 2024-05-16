package pool

import (
	xrconstant "dawn-server/impl/xr/lib/constant"
	"testing"
)

// todo menglc 覆盖率测试 [100%]

//go:generate go test -v -gcflags=all=-l -coverprofile=coverage.out
//go:generate go tool cover -html=coverage.out -o coverage.html

func TestBytePool(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	buf := MakeByteSlice(len(data))
	copy(buf[:len(data)], data)
	_ = ReleaseByteSlice(buf)
}

//测试跨协程使用sync.pool
//var dataChan chan []byte
//
//func TestBytePoolMulti(t *testing.T) {
//	dataChan = make(chan []byte, 1000)
//
//	go func() {
//		for d := range dataChan {
//			time.Sleep(time.Millisecond * 150)
//			_, err := util.ReleaseByteSlice(d)
//			fmt.Println("ReleaseByteSlice, get d:", d)
//			if err != nil {
//				fmt.Println("ReleaseByteSlice err :", err)
//			}
//			_ = d
//		}
//	}()
//	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	for i := 0; i < 10; i++ {
//		buf := util.MakeByteSlice(len(data))
//		copy(buf[:len(data)], data)
//		buf = append(buf, byte(i))
//		dataChan <- buf
//		time.Sleep(time.Millisecond * 100)
//	}
//	time.Sleep(time.Second)
//}

/*
*
1024*64:
$ go test -gcflags="-m -m -l" -benchmem -bench=Byte
# dawn-server/impl/xr/tool/script_test [dawn-server/impl/xr/tool/script.test]
.\byte_test.go:8:20: b does not escape
.\byte_test.go:9:16: []byte{...} does not escape
.\byte_test.go:11:14: make([]byte, 1024 * 64) does not escape
.\byte_test.go:16:24: b does not escape
.\byte_test.go:17:16: []byte{...} does not escape
# dawn-server/impl/xr/tool/script.test
C:\Users\k_wang\AppData\Local\Temp\go-build2084400977\b001\_testmain.go:47:42: testdeps.TestDeps{} escapes to heap:
C:\Users\k_wang\AppData\Local\Temp\go-build2084400977\b001\_testmain.go:47:42:   flow: {heap} = &{storage for testdeps.TestDeps{}}:
C:\Users\k_wang\AppData\Local\Temp\go-build2084400977\b001\_testmain.go:47:42:     from testdeps.TestDeps{} (spill) at $WORK\b001\_testmain.go:47:42
C:\Users\k_wang\AppData\Local\Temp\go-build2084400977\b001\_testmain.go:47:42:     from testing.MainStart(testdeps.TestDeps{}, tests, benchmarks, fuzzTargets, examples)
(call parameter) at $WORK\b001\_testmain.go:47:24
C:\Users\k_wang\AppData\Local\Temp\go-build2084400977\b001\_testmain.go:47:42: testdeps.TestDeps{} escapes to heap
goos: windows
goarch: amd64
pkg: dawn-server/impl/xr/tool/script
cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
BenchmarkByte-8         1000000000               0.2300 ns/op          0 B/op          0 allocs/op
BenchmarkBytePool-8     19644660                60.61 ns/op           24 B/op          1 allocs/op
PASS
ok      dawn-server/impl/xr/tool/script 1.542s

1024*64 + 1:
$ go test -gcflags="-m -m -l" -benchmem -bench=Byte
# dawn-server/impl/xr/tool/script_test [dawn-server/impl/xr/tool/script.test]
.\byte_test.go:11:14: make([]byte, 1024 * 64 + 1) escapes to heap:
.\byte_test.go:16:24: b does not escape
.\byte_test.go:17:16: []byte{...} does not escape
# dawn-server/impl/xr/tool/script.test
C:\Users\k_wang\AppData\Local\Temp\go-build1701714591\b001\_testmain.go:47:42: testdeps.TestDeps{} escapes to heap:
C:\Users\k_wang\AppData\Local\Temp\go-build1701714591\b001\_testmain.go:47:42:   flow: {heap} = &{storage for testdeps.TestDeps{}}:
C:\Users\k_wang\AppData\Local\Temp\go-build1701714591\b001\_testmain.go:47:42:     from testdeps.TestDeps{} (spill) at $WORK\b001\_testmain.go:47:42
C:\Users\k_wang\AppData\Local\Temp\go-build1701714591\b001\_testmain.go:47:42:     from testing.MainStart(testdeps.TestDeps{}, tests, benchmarks, fuzzTargets, examples)
(call parameter) at $WORK\b001\_testmain.go:47:24
C:\Users\k_wang\AppData\Local\Temp\go-build1701714591\b001\_testmain.go:47:42: testdeps.TestDeps{} escapes to heap
goos: windows
goarch: amd64
pkg: dawn-server/impl/xr/tool/script
cpu: Intel(R) Core(TM) i7-9700 CPU @ 3.00GHz
BenchmarkByte-8           220268              5664 ns/op           73728 B/op          1 allocs/op
BenchmarkBytePool-8     20125718                60.33 ns/op           24 B/op          1 allocs/op
PASS
ok      dawn-server/impl/xr/tool/script 3.620s
*/
func BenchmarkByte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = make([]byte, 1024*64) // 经测试 64K是临界值 64K+1 会慢很多 从栈逃逸到堆了

	}
}

func BenchmarkByte64k01byte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = make([]byte, 1024*64+1) // 经测试 64K是临界值 64K+1 会慢很多 从栈逃逸到堆了
	}
}

func BenchmarkBytePool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buf := MakeByteSlice(1024 * 64) // 此处注意bytePool类里的maxAreaValue上限 目前设定为64K 超过这个则直接make
		_ = ReleaseByteSlice(buf)
	}
}

func Test_bytePool_getPosByteSize(t *testing.T) {
	type args struct {
		idx  int
		size int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{ // 正常-success
			name: xrconstant.Normal,
			args: args{
				idx:  0,
				size: 0,
			},
			want: 0,
		},
		{ // 正常-success
			name: xrconstant.Normal,
			args: args{
				idx:  0,
				size: 2048,
			},
			want: 31,
		},
		{ // 超出范围-OutOfRange
			name: xrconstant.OutOfRange,
			args: args{
				idx:  0,
				size: 2048 + 1,
			},
			want: 32,
		},

		{ // 正常-success
			name: xrconstant.Normal,
			args: args{
				idx:  1,
				size: 2048 + 1,
			},
			want: 0,
		},
		{ // 正常-success
			name: xrconstant.Normal,
			args: args{
				idx:  1,
				size: 65536,
			},
			want: 30,
		},
		{ // 超出范围-OutOfRange
			name: xrconstant.OutOfRange,
			args: args{
				idx:  1,
				size: 65536 + 1,
			},
			want: 31,
		},
	}
	for _, v := range bytePoolList {
		t.Log(len(v.pool))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bytePoolList[tt.args.idx].getPosByteSize(tt.args.size); got != tt.want {
				t.Errorf("getPosByteSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
