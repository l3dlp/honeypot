package main

import (
	"crypto/sha256"
	"sync"
	"testing"
)

// 1. CPU Pur : Récursion lourde (Fibonacci)
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func Benchmark1_CPUPur(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fib(25) // Étape 25 pour pousser le CPU un peu plus loin
	}
}

// 2. Concurrence & Threads : Cryptographie parallèle (Goroutines)
func Benchmark2_Concurrence(b *testing.B) {
	data := []byte("intelligences-agency-secret-string-to-hash-12345")
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		// On lance 50 Goroutines en parallèle qui calculent des SHA-256
		for g := 0; g < 50; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				h := sha256.New()
				for j := 0; j < 10; j++ {
					h.Write(data)
					_ = h.Sum(nil)
				}
			}()
		}
		wg.Wait()
	}
}

// 3. Gestion RAM & Garbage Collector : Allocations massives de Slices
func Benchmark3_AllocationRAM(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// On force Go à allouer et libérer de la mémoire à chaque itération
		_ = make([]byte, 10*1024*1024) // 10 Mo
	}
}

