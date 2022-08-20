package utils

type CompareRelation int

const (
	Lesser  CompareRelation = -1
	Equal   CompareRelation = 0
	Greater CompareRelation = 1
)

type String[V any] func(V) string
type Compare[T any] func(x T) CompareRelation
