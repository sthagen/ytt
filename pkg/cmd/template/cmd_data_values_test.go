package template_test

import (
	"testing"

	cmdcore "github.com/k14s/ytt/pkg/cmd/core"
	cmdtpl "github.com/k14s/ytt/pkg/cmd/template"
	"github.com/k14s/ytt/pkg/files"
)

func TestDataValues(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data_int: #@ data.values.int
data_str: #@ data.values.str`)

	expectedYAMLTplData := `data_int: 123
data_str: str
`

	yamlData := []byte(`
#@data/values
---
int: 123
str: str`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}

func TestDataValuesWithFlags(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data_int: #@ data.values.int
data_str: #@ data.values.str
values: #@ data.values`)

	expectedYAMLTplData := `data_int: 124
data_str: str
values:
  int: 124
  str: str
  boolean: true
  nested:
    value: str
  another:
    nested:
      map: 567
`

	yamlData := []byte(`
#@data/values
---
int: 123
str: str
boolean: false
nested:
  value: not-str
another:
  nested:
    map: {"a": 123}`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	opts.DataValuesFlags = cmdtpl.DataValuesFlags{
		// TODO env and files?
		KVsFromYAML: []string{"int=124", "boolean=true", "nested.value=\"str\"", "another.nested.map=567"},
	}

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<< vs >>>%s<<<", file.Bytes(), expectedYAMLTplData)
	}
}

func TestDataValuesWithFlagsMarkedMissingOk(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
values: #@ data.values`)

	expectedYAMLTplData := `values:
  nested:
    value: str
  another_nested:
    other_value: str2
`

	yamlData := []byte(`
#@data/values
---
nested:
  value: str
`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	opts.DataValuesFlags = cmdtpl.DataValuesFlags{
		// TODO add nested.value2*=str2 since replace with 0 nodes does not do anything
		KVsFromYAML: []string{"another_nested+.other_value=str2"},
	}

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<< vs >>>%s<<<", file.Bytes(), expectedYAMLTplData)
	}
}

func TestDataValuesMultipleFiles(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data_int: #@ data.values.int
data_str: #@ data.values.str`)

	expectedYAMLTplData := `data_int: 123
data_str: str2
`

	yamlData1 := []byte(`
#@data/values
---
int: 123
str: str`)

	yamlData2 := []byte(`
#@data/values
---
str: str2`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data1.yml", yamlData1)),
		files.MustNewFileFromSource(files.NewBytesSource("data2.yml", yamlData2)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}

func TestDataValuesMultipleInOneFile(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data_int: #@ data.values.int
data_str: #@ data.values.str`)

	expectedYAMLTplData := `data_int: 123
data_str: str2
`

	yamlData := []byte(`
#@data/values
---
str: str

#@data/values
---
str: str2
#@overlay/match missing_ok=True
int: 123`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}

func TestDataValuesOverlayNewKey(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data_int: #@ data.values.int
data_str: #@ data.values.str`)

	expectedYAMLTplData := `data_int: 123
data_str: str2
`

	yamlData1 := []byte(`
#@data/values
---
str: str`)

	yamlData2 := []byte(`
#@data/values
---
str: str2
#@overlay/match missing_ok=True
int: 123`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data1.yml", yamlData1)),
		files.MustNewFileFromSource(files.NewBytesSource("data2.yml", yamlData2)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}

func TestDataValuesOverlayRemoveKey(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data: #@ data.values.data`)

	expectedYAMLTplData := `data:
  str: str
`

	yamlData1 := []byte(`
#@data/values
---
data:
  str: str
  int: 123`)

	yamlData2 := []byte(`
#@data/values
---
data:
  #@overlay/remove
  int: null`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data1.yml", yamlData1)),
		files.MustNewFileFromSource(files.NewBytesSource("data2.yml", yamlData2)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}

func TestDataValuesWithNonDataValuesDocsErr(t *testing.T) {
	yamlData := []byte(`
#@data/values
---
str: str
---
non-data-values-doc`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err == nil {
		t.Fatalf("Expected RunWithFiles to fail, but was no error")
	}

	if out.Err.Error() != "Overlaying data values (in following order: data.yml): Templating file 'data.yml': Expected data values file 'data.yml' to only have data values documents" {
		t.Fatalf("Expected RunWithFiles to fail, but was '%s'", out.Err)
	}
}

func TestDataValuesWithNonDocDataValuesErr(t *testing.T) {
	yamlData := []byte(`
---
#@data/values
str: str`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("data.yml", yamlData)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err == nil {
		t.Fatalf("Expected RunWithFiles to fail, but was no error")
	}

	if out.Err.Error() != "Expected YAML document to be annotated with data/values but was *yamlmeta.MapItem" {
		t.Fatalf("Expected RunWithFiles to fail, but was '%s'", out.Err)
	}
}

func TestDataValuesOverlayChildDefaults(t *testing.T) {
	yamlTplData := []byte(`
#@ load("@ytt:data", "data")
data: #@ data.values.data`)

	expectedYAMLTplData := `data:
  str: str
  int: 123
  bool: true
`

	yamlData1 := []byte(`
#@data/values
---
data:
  str: str
  int: 123`)

	yamlData2 := []byte(`
#@data/values
#@overlay/match-child-defaults missing_ok=True
---
data:
  bool: true`)

	filesToProcess := files.NewSortedFiles([]*files.File{
		files.MustNewFileFromSource(files.NewBytesSource("tpl.yml", yamlTplData)),
		files.MustNewFileFromSource(files.NewBytesSource("data1.yml", yamlData1)),
		files.MustNewFileFromSource(files.NewBytesSource("data2.yml", yamlData2)),
	})

	ui := cmdcore.NewPlainUI(false)
	opts := cmdtpl.NewOptions()

	out := opts.RunWithFiles(cmdtpl.TemplateInput{Files: filesToProcess}, ui)
	if out.Err != nil {
		t.Fatalf("Expected RunWithFiles to succeed, but was error: %s", out.Err)
	}

	if len(out.Files) != 1 {
		t.Fatalf("Expected number of output files to be 1, but was %d", len(out.Files))
	}

	file := out.Files[0]

	if file.RelativePath() != "tpl.yml" {
		t.Fatalf("Expected output file to be tpl.yml, but was %#v", file.RelativePath())
	}

	if string(file.Bytes()) != expectedYAMLTplData {
		t.Fatalf("Expected output file to have specific data, but was: >>>%s<<<", file.Bytes())
	}
}
