package limiter

import (
	"os"
	"testing"
	"time"
)

var (
	sl       Limiter
	dl       Limiter
	rate     = 100
	size     = 1000
	interval = time.Millisecond * 100
)

func TestMain(m *testing.M) {
	var err error
	sl, err = New(Option{
		Name:     "",
		Mode:     ModeSingle,
		Rate:     rate,
		Size:     size,
		Interval: interval,
	})
	if err != nil {
		panic(err)
	}

	dl, err = New(Option{
		Name:     "test",
		Mode:     ModeDistributed,
		Rate:     rate,
		Size:     size,
		Interval: interval,
		Dsn:      "redis://:@127.0.0.1:6379/0",
	})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestLimiter_SingleGetToken(t *testing.T) {
	for j := 0; j < size; j++ {
		_, ok := sl.GetToken(1)
		if !ok {
			t.Fatalf("got = %v, want = %v", ok, true)
		}
	}
	time.Sleep(time.Second + 2*interval)
	for j := 0; j < size; j++ {
		n, ok := sl.GetToken(1)
		if !ok {
			t.Fatalf("got = %v, n = %d, want = %v", ok, n, true)
		}
	}
}

func TestLimiter_SingleGetToken2(t *testing.T) {
	time.Sleep(time.Second)
	n, ok := sl.GetToken(size + 1)
	if ok {
		t.Fatalf("got = %v, n = %d, want = %v", ok, n, false)
	}
}

func TestLimiter_DistributedGetToken(t *testing.T) {
	for j := 0; j < size; j++ {
		n, ok := dl.GetToken(1)
		if !ok {
			t.Fatalf("got = %v, n = %d, want = %v", ok, n, true)
		}
	}
}

func TestLimiter_DistributedGetToken2(t *testing.T) {
	time.Sleep(time.Second)
	n, ok := dl.GetToken(size + 1)
	if ok {
		t.Fatalf("got = %v, n = %d, want = %v", ok, n, false)
	}
}
