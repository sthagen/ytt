package yttlibrary

import (
	"github.com/k14s/ytt/pkg/template"
	"github.com/k14s/ytt/pkg/template/core"
	"github.com/k14s/ytt/pkg/yamlmeta"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type dataModule struct {
	values starlark.Value
	loader template.CompiledTemplateLoader
}

func NewDataModule(values *yamlmeta.Document, loader template.CompiledTemplateLoader) dataModule {
	val := core.NewGoValueWithOpts(values.AsInterface(), core.GoValueOpts{MapIsStruct: true})
	return dataModule{val.AsStarlarkValue(), loader}
}

func (b dataModule) AsModule() starlark.StringDict {
	return starlark.StringDict{
		"data": &starlarkstruct.Module{
			Name: "data",
			Members: starlark.StringDict{
				"list": starlark.NewBuiltin("data.list", core.ErrWrapper(b.List)),
				"read": starlark.NewBuiltin("data.read", core.ErrWrapper(b.Read)),
				// TODO write?
				"values": b.values,
			},
		},
	}
}

func (b dataModule) List(thread *starlark.Thread, f *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	return b.loader.ListData(thread, f, args, kwargs)
}

func (b dataModule) Read(thread *starlark.Thread, f *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	return b.loader.LoadData(thread, f, args, kwargs)
}
