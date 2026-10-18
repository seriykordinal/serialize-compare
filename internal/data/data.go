package data

import (
	"fmt"
	"math/rand"

	"serialize-compare/internal/model"
)

// Seed фиксирован, чтобы данные были одинаковыми при каждом запуске.
const Seed = 42

var cities = []string{"Moscow", "Kazan", "Omsk", "Tomsk", "Perm"}

// NewRand возвращает генератор с фиксированным seed.
func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(Seed))
}

// Simple возвращает простую структуру.
func Simple() model.Simple {
	return model.Simple{ID: 1, Name: "benchmark"}
}

// Nested генерирует одну запись средней сложности.
func Nested(r *rand.Rand, i int) model.Nested {
	items := make([]model.Item, 2+r.Intn(3)) // от 2 до 4 элементов
	for j := range items {
		items[j] = model.Item{
			SKU:   fmt.Sprintf("SKU-%05d", r.Intn(100000)),
			Qty:   int32(1 + r.Intn(10)),
			Price: float64(r.Intn(1000000)) / 100,
		}
	}

	return model.Nested{
		ID:     int64(i),
		Name:   fmt.Sprintf("user_%d", i),
		Active: i%2 == 0,
		Score:  r.Float64() * 100,
		Tags:   []string{"go", "rpc", fmt.Sprintf("group-%d", i%7)},
		Address: model.Address{
			City:   cities[r.Intn(len(cities))],
			Street: fmt.Sprintf("Lenina %d", r.Intn(200)),
			Zip:    fmt.Sprintf("%06d", 100000+r.Intn(900000)),
		},
		Items: items,
	}
}

// Batch генерирует n записей Nested — «большой» набор данных.
func Batch(n int) []model.Nested {
	r := NewRand()
	out := make([]model.Nested, n)
	for i := range out {
		out[i] = Nested(r, i)
	}
	return out
}
