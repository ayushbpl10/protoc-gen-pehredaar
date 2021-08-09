package main_test

import (
	"bytes"
	peh "github.com/ayushbpl10/protoc-gen-pehredaar"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/spf13/afero"
	"os"
	"testing"
)

func TestModule(t *testing.T) {
	req, err := os.Open("./code_generator_request.pb.bin")
	if err != nil {
		t.Fatal(err)
	}

	fs := afero.NewMemMapFs()
	res := &bytes.Buffer{}

	pgs.Init(
		pgs.ProtocInput(req),  // use the pre-generated request
		pgs.ProtocOutput(res), // capture CodeGeneratorResponse
		pgs.FileSystem(fs),    // capture any custom files written directly to disk
	).RegisterModule(&peh.RightsGen{ModuleBase: pgs.ModuleBase{}}).Render()

	// check res and the fs for output
}
