package typeutils

type CompareRelation int

const (
	Lesser  CompareRelation = -1
	Equal   CompareRelation = 0
	Greater CompareRelation = 1
)

type String[V any] func(V) string
type Compare[V any] func(V) CompareRelation
