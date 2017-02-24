package main

import (
	"os/exec"
	"testing"

	"github.com/drud/build-tools/tests/pkg/version"
	"github.com/stretchr/testify/assert"
	"strings"
	"log"
)

var (
	osname = "" // The operating system.
)


func init() {
	// Operating system - Darwin or Linux
	v, err := exec.Command("uname", "-s").Output()
	if (err != nil) {
		log.Fatalln("Failed to run uname command:", string(v))
	}
	osname = strings.TrimSpace(string(v))
}

// Runs a number of standard make targets and test for basic sanity of result
// Assumes operation in the "testing" directory where the Makefile is
func TestMake(t *testing.T) {
	assert := assert.New(t)

	// Map OS name to output location
	binlocs := map[string]string{
		"Darwin": "bin/darwin/darwin_amd64",
		"Linux":  "bin/linux",
	}


	v, err := exec.Command("which", "make").Output()
	assert.Contains(string(v), "make")

	// Try trivial "make version"
	v, err = exec.Command("make", "version").Output()
	assert.NoError(err)
	assert.Contains(string(v), "VERSION:" + version.VERSION)

	// Run a make clean to start with
	v, err = exec.Command("make", "clean").Output()
	assert.NoError(err)

	// Build darwin and linux cmds
	v, err = exec.Command("make", "darwin").Output()
	assert.NoError(err)
	assert.Contains(string(v), "building darwin")

	v, err = exec.Command("make", "linux").Output()
	assert.NoError(err)
	assert.Contains(string(v), "building linux")

	// Run the native gofmtproblem application to make sure it runs
	v, err = exec.Command(binlocs[osname] + "/gofmtproblem").Output()
	assert.Contains(string(v), "This is gofmtproblem.go")
	assert.Contains(string(v), version.VERSION)

	// Make container
	v, err = exec.Command("make", "container").Output()
	assert.Contains(string(v), "Successfully built")

}

func TestGoFmt(t *testing.T) {
	assert := assert.New(t)

	// Test "make gofmt
	v, _ := exec.Command("make", "gofmt").Output()
	// assert.Error(err) // TODO: This should have an error because we have actual gofmt problems
	assert.Contains(string(v), "gofmtproblem.go")
}