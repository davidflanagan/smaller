package smaller

import (
	"math/rand"
	"strconv"
	"testing"
)

func randomString(rng *rand.Rand, length int) string {
	bytes := make([]byte, length)
	rng.Read(bytes)
	return string(bytes)
}

func TestStringRoundtrip(t *testing.T) {
	var randomStrings []string

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 10240; i++ {
		s := randomString(rng, i)
		randomStrings = append(randomStrings, s)
	}

	for i, s := range randomStrings {
		S := NewString(s)
		l := S.Len()
		if l != i {
			t.Fatalf("Expected length %d but got %d", i, l)
		}

		s2 := S.String()
		if s2 != s {
			t.Fatalf("Expected %s but got %s", s, s2)
		}
	}
}

func BenchmarkSmallerString(b *testing.B) {
	lengths := []int{0, 1, 3, 7, 15, 31, 63, 127, 256}
	strings := []string{}
	rng := rand.New(rand.NewSource(1))
	for _, length := range lengths {
		strings = append(strings, randomString(rng, length))
	}

	var storage []*String = make([]*String, 1000)

	for n, length := range lengths {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			s := strings[n]
			for i := 0; i < b.N; i++ {
				for j := 0; j < 1000; j++ {
					ss := NewString(s)
					storage[j] = &ss
				}
			}
		})
	}
}

func BenchmarkRegularString(b *testing.B) {
	lengths := []int{0, 1, 3, 7, 15, 31, 63, 127, 256}
	strings := [][]byte{}
	rng := rand.New(rand.NewSource(1))
	for _, length := range lengths {
		strings = append(strings, []byte(randomString(rng, length)))
	}

	var storage []*string = make([]*string, 1000)

	for n, length := range lengths {
		b.Run(strconv.Itoa(length), func(b *testing.B) {
			s := strings[n]
			for i := 0; i < b.N; i++ {
				for j := 0; j < 1000; j++ {
					ss := string(s)
					storage[j] = &ss
				}
			}
		})
	}
}
