package odata

import (
	"github.com/engswee/flashpipe/httpclnt"
)

type ScriptCollection struct {
	exe *httpclnt.HTTPExecuter
	typ string
}

// NewScriptCollection returns an initialised ScriptCollection instance.
func NewScriptCollection(exe *httpclnt.HTTPExecuter) DesigntimeArtifact {
	sc := new(ScriptCollection)
	sc.exe = exe
	sc.typ = "ScriptCollection"
	return sc
}

func (sc *ScriptCollection) Create(id string, name string, packageId string, artifactDir string) error {
	return create(id, name, packageId, artifactDir, sc.typ, sc.exe)
}
func (sc *ScriptCollection) Update(id string, name string, packageId string, artifactDir string) (err error) {
	return update(id, name, packageId, artifactDir, sc.typ, sc.exe)
}
func (sc *ScriptCollection) Deploy(id string) (err error) {
	return deploy(id, sc.typ, sc.exe)
}
func (sc *ScriptCollection) Delete(id string) (err error) {
	return deleteCall(id, sc.typ, sc.exe)
}
func (sc *ScriptCollection) Get(id string, version string) (string, bool, error) {
	return get(id, version, sc.typ, sc.exe)
}
func (sc *ScriptCollection) GetContent(id string, version string) ([]byte, error) {
	return getContent(id, version, sc.typ, sc.exe)
}
func (sc *ScriptCollection) DiffContent(firstDir string, secondDir string) bool {
	return diffContent(firstDir, secondDir)
}
func (sc *ScriptCollection) CopyContent(srcDir string, tgtDir string) error {
	return copyContent(srcDir, tgtDir)
}