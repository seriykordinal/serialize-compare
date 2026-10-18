package bench_test

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"

	"serialize-compare/internal/convert"
	"serialize-compare/internal/data"
	"serialize-compare/internal/model"
	"serialize-compare/internal/pb"
)

// format описывает один формат сериализации.
type format struct {
	name   string
	encode func(v any) ([]byte, error)
	decode func(d []byte, v any) error
}

var formats = []format{
	{
		name:   "json",
		encode: func(v any) ([]byte, error) { return json.Marshal(v) },
		decode: func(d []byte, v any) error { return json.Unmarshal(d, v) },
	},
	{
		name:   "msgpack",
		encode: func(v any) ([]byte, error) { return msgpack.Marshal(v) },
		decode: func(d []byte, v any) error { return msgpack.Unmarshal(d, v) },
	},
	{
		name:   "protobuf",
		encode: func(v any) ([]byte, error) { return proto.Marshal(v.(proto.Message)) },
		decode: func(d []byte, v any) error { return proto.Unmarshal(d, v.(proto.Message)) },
	},
}

// benchDataset запускает бенчмарки encode и decode для каждого формата.
//
//	goVal, goNew  — значение и конструктор пустого объекта для JSON и MessagePack;
//	pbVal, pbNew  — то же для Protobuf.
func benchDataset(b *testing.B, goVal any, goNew func() any, pbVal any, pbNew func() any) {
	for _, f := range formats {
		val, newFn := goVal, goNew
		if f.name == "protobuf" {
			val, newFn = pbVal, pbNew
		}

		// Заранее кодируем данные для декодирования (вне замера).
		encoded, err := f.encode(val)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(f.name+"/encode", func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(encoded)))
			for i := 0; i < b.N; i++ {
				if _, err := f.encode(val); err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(f.name+"/decode", func(b *testing.B) {
			// Проверка корректности до замера.
			if err := f.decode(encoded, newFn()); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(encoded)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := f.decode(encoded, newFn()); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkSimple — простая структура из двух полей.
func BenchmarkSimple(b *testing.B) {
	v := data.Simple()
	benchDataset(b,
		v, func() any { return new(model.Simple) },
		convert.Simple(v), func() any { return new(pb.Simple) },
	)
}

// BenchmarkNested — вложенная структура с массивами и объектами.
func BenchmarkNested(b *testing.B) {
	v := data.Nested(data.NewRand(), 0)
	benchDataset(b,
		v, func() any { return new(model.Nested) },
		convert.Nested(v), func() any { return new(pb.Nested) },
	)
}

// BenchmarkBatch — массив из 1000 вложенных записей.
func BenchmarkBatch(b *testing.B) {
	v := data.Batch(1000)
	benchDataset(b,
		v, func() any { return new([]model.Nested) },
		convert.Batch(v), func() any { return new(pb.Batch) },
	)
}

// TestPayloadSizes выводит размер сериализованных данных в байтах.
// Запуск: go test -run TestPayloadSizes -v ./bench/
func TestPayloadSizes(t *testing.T) {
	s := data.Simple()
	n := data.Nested(data.NewRand(), 0)
	bt := data.Batch(1000)

	cases := []struct {
		name  string
		goVal any
		pbVal any
	}{
		{"simple", s, convert.Simple(s)},
		{"nested", n, convert.Nested(n)},
		{"batch", bt, convert.Batch(bt)},
	}

	for _, c := range cases {
		for _, f := range formats {
			val := c.goVal
			if f.name == "protobuf" {
				val = c.pbVal
			}
			out, err := f.encode(val)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("%-8s %-9s %8d bytes", c.name, f.name, len(out))
		}
	}
}
