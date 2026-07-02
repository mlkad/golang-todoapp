package domain

/*
	field not provided
	field provided: null
	field provided: value
*/

type NullableString struct {
	Value *string
	Set bool
} 

type Nullable[T any] struct {
	Value *T
	Set bool
}
