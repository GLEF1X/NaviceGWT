package main

import (
	"io"
	"os"
	"testing"
)

func BenchmarkTransformData(b *testing.B) {
	f, err := os.Open("./input.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		TransformData(data)
	}
}
