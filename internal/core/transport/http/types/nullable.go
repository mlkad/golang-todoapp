package core_http_types

import (
	"encoding/json"

	"github.com/mlkad/golang-todoapp/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}
/*
мой UnmarshalJSON оборачивает json.Unmarshal, 
чтобы вокруг добавить Set = true и проверку на null
*/

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T 
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set: n.Set,
	}
}