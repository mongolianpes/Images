package main

import (
	"encoding/json"
	"testing"
)

var testUser = User{
	ID:    1,
	Name:  "Alice",
	Email: "alice@example.com",
}

var testJSON = []byte(`{"id":1,"name":"Alice","email":"alice@example.com"}`)

func BenchmarkStdMarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(testUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdUnmarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var u User
		err := json.Unmarshal(testJSON, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStdMarshalWrapper(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := UnmarshalUser(testJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterMarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := jsonIter.Marshal(testUser)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterUnmarshal(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var u User
		err := jsonIter.Unmarshal(testJSON, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJsoniterMarshalWrapper(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := JsoniterUnmarshalUser(testJSON)
		if err != nil {
			b.Fatal(err)
		}
	}
}
