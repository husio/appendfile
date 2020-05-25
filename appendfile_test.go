package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func BenchmarkAppendNoSync(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendSync(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}
		if err := fd.Sync(); err != nil {
			b.Fatalf("sync: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendFdatasyncOpenDsync(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY|unix.O_DSYNC, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}

		if err := unix.Fdatasync(int(fd.Fd())); err != nil {
			b.Fatalf("fdatasync: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendFdatasync(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}

		if err := unix.Fdatasync(int(fd.Fd())); err != nil {
			b.Fatalf("fdatasync: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendFdatasyncAndFallocateDefaultMode(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}

	if err := unix.Fallocate(int(fd.Fd()), 0, 0, 1e7); err != nil {
		b.Fatalf("fallocate: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}

		if err := unix.Fdatasync(int(fd.Fd())); err != nil {
			b.Fatalf("fdatasync: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendFdatasyncAndFallocateDefaultModeOpenDsync(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY|unix.O_DSYNC, 0644)
	if err != nil {
		b.Fatal(err)
	}

	if err := unix.Fallocate(int(fd.Fd()), 0, 0, 1e7); err != nil {
		b.Fatalf("fallocate: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}

		if err := unix.Fdatasync(int(fd.Fd())); err != nil {
			b.Fatalf("fdatasync: %s", err)
		}
	}
	b.StopTimer()
}

func BenchmarkAppendFdatasyncAndFallocateZero(b *testing.B) {
	fd, err := os.OpenFile(tempfile(b), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}

	if err := unix.Fallocate(int(fd.Fd()), unix.FALLOC_FL_ZERO_RANGE|unix.FALLOC_FL_KEEP_SIZE, 0, 1e7); err != nil {
		b.Fatalf("fallocate: %s", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := fd.Write([]byte("1234567890qwertyuiop")); err != nil {
			b.Fatalf("write: %s", err)
		}

		if err := unix.Fdatasync(int(fd.Fd())); err != nil {
			b.Fatalf("fdatasync: %s", err)
		}
	}
	b.StopTimer()
}

func tempfile(t testing.TB) string {
	t.Helper()

	fpath := fmt.Sprintf("%s/benchmark-append.%d", os.TempDir(), time.Now().UnixNano())
	t.Cleanup(func() { _ = os.Remove(fpath) })
	return fpath
}
