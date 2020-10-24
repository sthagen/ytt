// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package yamlmeta

import (
	"github.com/k14s/ytt/pkg/filepos"
)

type Node interface {
	GetPosition() *filepos.Position

	GetValues() []interface{} // ie children
	SetValue(interface{}) error
	AddValue(interface{}) error
	ResetValue()

	GetMetas() []*Meta
	addMeta(*Meta)

	GetAnnotations() interface{}
	SetAnnotations(interface{})

	DeepCopyAsInterface() interface{}
	DeepCopyAsNode() Node

	Check() TypeCheck

	_private()
}

// Ensure: all types are — in fact — assignable to Node
var _ = []Node{&DocumentSet{}, &Document{}, &Map{}, &MapItem{}, &Array{}, &ArrayItem{}}

type DocumentSet struct {
	Metas    []*Meta
	AllMetas []*Meta

	Items    []*Document
	Position *filepos.Position

	annotations   interface{}
	originalBytes *[]byte
}

type Document struct {
	Metas    []*Meta
	Value    interface{}
	Position *filepos.Position
	Type     interface{}

	annotations interface{}
	injected    bool // indicates that Document was not present in the parsed content
}

type Map struct {
	Metas    []*Meta
	Items    []*MapItem
	Position *filepos.Position
	Type     *MapType

	annotations interface{}
}

type MapItem struct {
	Metas    []*Meta
	Key      interface{}
	Value    interface{}
	Position *filepos.Position
	Type     interface{}

	annotations interface{}
}

type Array struct {
	Metas    []*Meta
	Items    []*ArrayItem
	Position *filepos.Position

	annotations interface{}
}

type ArrayItem struct {
	Metas    []*Meta
	Value    interface{}
	Position *filepos.Position

	annotations interface{}
}

type Meta struct {
	Data     string
	Position *filepos.Position
}
