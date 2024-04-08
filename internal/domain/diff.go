package domain

import (
	"encoding/json"
	"log"
	"slices"
)

type Diff interface {
	Type() DiffType
}

type DiffType string

const (
	DiffTypeAddition DiffType = "addition"
	DiffTypeReplace  DiffType = "replacement"
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

func (a Addition) String() string {
	str := struct {
		Type, Key string
		Value     any
	}{
		Type:  string(a.Type()),
		Key:   a.Key,
		Value: a.Value,
	}
	jsonBytes, err := json.Marshal(str)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(jsonBytes)
}

type Replace struct {
	Key string
	New string
	Old string
}

func (Replace) Type() DiffType {
	return DiffTypeReplace
}

func (r Replace) String() string {
	str := struct {
		Type, Key          string
		OldValue, NewValue string
	}{
		Type:     string(r.Type()),
		Key:      r.Key,
		OldValue: r.Old,
		NewValue: r.New,
	}
	jsonBytes, err := json.Marshal(str)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(jsonBytes)
}

type Deletion struct {
	Key   string
	Value any
}

func (Deletion) Type() DiffType {
	return DiffTypeDeletion
}

func (d Deletion) String() string {
	str := struct {
		Type, Key string
		Value     any
	}{
		Type:  string(d.Type()),
		Key:   d.Key,
		Value: d.Value,
	}
	jsonBytes, err := json.Marshal(str)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(jsonBytes)
}
