package yttlibrary

import (
	"fmt"

	"github.com/k14s/ytt/pkg/template/core"
	"github.com/k14s/ytt/pkg/yamlmeta"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	YAMLAPI = starlark.StringDict{
		"yaml": &starlarkstruct.Module{
			Name: "yaml",
			Members: starlark.StringDict{
				"encode": starlark.NewBuiltin("yaml.encode", core.ErrWrapper(yamlModule{}.Encode)),
				"decode": starlark.NewBuiltin("yaml.decode", core.ErrWrapper(yamlModule{}.Decode)),
			},
		},
	}
)

type yamlModule struct{}

func (b yamlModule) Encode(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	val := core.NewStarlarkValue(args.Index(0)).AsInterface()
	docSet := yamlmeta.NewDocumentSetFromInterface(val)

	valBs, err := docSet.AsBytes()
	if err != nil {
		return starlark.None, err
	}

	return starlark.String(string(valBs)), nil
}

func (b yamlModule) Decode(thread *starlark.Thread, f *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if args.Len() != 1 {
		return starlark.None, fmt.Errorf("expected exactly one argument")
	}

	valEncoded, err := core.NewStarlarkValue(args.Index(0)).AsString()
	if err != nil {
		return starlark.None, err
	}

	var valDecoded interface{}

	err = yamlmeta.PlainUnmarshal([]byte(valEncoded), &valDecoded)
	if err != nil {
		return starlark.None, err
	}

	return core.NewGoValue(valDecoded, false).AsStarlarkValue(), nil
}
