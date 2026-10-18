package convert

import (
	"serialize-compare/internal/model"
	"serialize-compare/internal/pb"
)

// Simple переводит model.Simple в protobuf-сообщение.
func Simple(m model.Simple) *pb.Simple {
	return &pb.Simple{Id: m.ID, Name: m.Name}
}

// Nested переводит model.Nested в protobuf-сообщение.
func Nested(m model.Nested) *pb.Nested {
	items := make([]*pb.Item, len(m.Items))
	for i, it := range m.Items {
		items[i] = &pb.Item{Sku: it.SKU, Qty: it.Qty, Price: it.Price}
	}
	return &pb.Nested{
		Id:      m.ID,
		Name:    m.Name,
		Active:  m.Active,
		Score:   m.Score,
		Tags:    m.Tags,
		Address: &pb.Address{City: m.Address.City, Street: m.Address.Street, Zip: m.Address.Zip},
		Items:   items,
	}
}

// Batch переводит срез записей в protobuf-сообщение.
func Batch(ms []model.Nested) *pb.Batch {
	recs := make([]*pb.Nested, len(ms))
	for i := range ms {
		recs[i] = Nested(ms[i])
	}
	return &pb.Batch{Records: recs}
}
