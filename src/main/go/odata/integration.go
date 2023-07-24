package odata

import (
	"github.com/engswee/flashpipe/httpclnt"
)

type Integration struct {
	exe *httpclnt.HTTPExecuter
	typ string
}

// NewIntegration returns an initialised Integration instance.
func NewIntegration(exe *httpclnt.HTTPExecuter) DesigntimeArtifact {
	i := new(Integration)
	i.exe = exe
	i.typ = "Integration"
	return i
}

func (int *Integration) Create(id string, name string, packageId string, artifactDir string) error {
	return create(id, name, packageId, artifactDir, int.typ, int.exe)
}
func (int *Integration) Update(id string, name string, packageId string, artifactDir string) error {
	return update(id, name, packageId, artifactDir, int.typ, int.exe)
}
func (int *Integration) Deploy(id string) error {
	return deploy(id, int.typ, int.exe)
}
func (int *Integration) Delete(id string) error {
	return deleteCall(id, int.typ, int.exe)
}
func (int *Integration) Get(id string, version string) (string, bool, error) {
	return get(id, version, int.typ, int.exe)
}
func (int *Integration) GetContent(id string, version string) ([]byte, error) {
	return getContent(id, version, int.typ, int.exe)
}
func (int *Integration) DiffContent(firstDir string, secondDir string) bool {
	return diffContent(firstDir, secondDir)
}
func (int *Integration) CopyContent(srcDir string, tgtDir string) error {
	return copyContent(srcDir, tgtDir)
}