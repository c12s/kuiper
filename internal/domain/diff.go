package domain

import "slices"

type Diff interface {
	Type() DiffType
}

type DiffType string

const (
	DiffTypeAddition DiffType = "addition"
	DiffTypeReplace  DiffType = "replace"
	DiffTypeDeletion DiffType = "deletion"
)

func GetDiffTypeValues() []DiffType {
	return []DiffType{
		DiffTypeAddition,
		DiffTypeReplace,
		DiffTypeDeletion,
	}
}

func (dt *DiffType) IsValid() bool {
	if dt != nil && slices.Contains(GetDiffTypeValues(), *dt) {
		return true
	}

	return false
}

type Addition struct {
	Key   string
	Value any
}

func (Addition) Type() DiffType {
	return DiffTypeAddition
}

type Replace struct {
	Key string
	New string
	Old string
}

func (Replace) Type() DiffType {
	return DiffTypeReplace
}

type Deletion struct {
	Key   string
	Value any
}

func (Deletion) Type() DiffType {
	return DiffTypeDeletion
}
