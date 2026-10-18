package model

// Simple — простая структура: два поля.
type Simple struct {
	ID   int64  `json:"id" msgpack:"id"`
	Name string `json:"name" msgpack:"name"`
}

// Address — вложенный объект.
type Address struct {
	City   string `json:"city" msgpack:"city"`
	Street string `json:"street" msgpack:"street"`
	Zip    string `json:"zip" msgpack:"zip"`
}

// Item — элемент массива.
type Item struct {
	SKU   string  `json:"sku" msgpack:"sku"`
	Qty   int32   `json:"qty" msgpack:"qty"`
	Price float64 `json:"price" msgpack:"price"`
}

// Nested — структура средней сложности: скаляры, массив строк,
// вложенный объект и массив объектов.
type Nested struct {
	ID      int64    `json:"id" msgpack:"id"`
	Name    string   `json:"name" msgpack:"name"`
	Active  bool     `json:"active" msgpack:"active"`
	Score   float64  `json:"score" msgpack:"score"`
	Tags    []string `json:"tags" msgpack:"tags"`
	Address Address  `json:"address" msgpack:"address"`
	Items   []Item   `json:"items" msgpack:"items"`
}
