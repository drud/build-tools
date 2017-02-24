package main

import (
	"os/exec"
	"testing"

	"fmt"
	"github.com/drud/build-tools/tests/pkg/version"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
)

var (
	osname = "" // The operating system.
)

func init() {
	os.Chdir("../..")

	// Operating system - Darwin or Linux
	v, err := exec.Command("uname", "-s").Output()
	if err != nil {
		log.Fatalln("Failed to run uname command:", string(v))
	}
	osname = strings.TrimSpace(string(v))
}

// Runs a number of standard make targets and test for basic sanity of result
// Assumes operation in the "testing" directory where the Makefile is
func TestMake(t *testing.T) {
	a := assert.New(t)

	// Map OS name to output location
	binlocs := map[string]string{
		"Darwin": "bin/darwin/darwin_amd64",
		"Linux":  "bin/linux",
	}

	dir, _ := os.Getwd()
	fmt.Println("Current Directory:", dir)

	v, err := exec.Command("which", "make").Output()
	a.Contains(string(v), "make")

	// Try trivial "make version"
	v, err = exec.Command("make", "version").Output()
	a.NoError(err)
	a.Contains(string(v), "VERSION:"+version.VERSION)
	if err != nil {
		log.Fatalln("make version in", dir, "failed, so exiting. output=", string(v))
	}

	// Run a make clean to start with; linux requires sudo because container left things a mess
	v, err = exec.Command("sudo","make", "clean").Output()
	a.NoError(err)

	// Build darwin and linux cmds
	v, err = exec.Command("make", "darwin").Output()
	a.NoError(err)
	a.Contains(string(v), "building darwin")

	v, err = exec.Command("make", "linux").Output()
	a.NoError(err)
	a.Contains(string(v), "building linux")

	// Run the native gofmtproblem application to make sure it runs
	v, err = exec.Command(binlocs[osname] + "/gofmtproblem").Output()
	a.Contains(string(v), "This is gofmtproblem.go")
	a.Contains(string(v), version.VERSION)

	// Make container
	v, err = exec.Command("make", "container").Output()
	a.Contains(string(v), "Successfully built")

}

// Try gofmt - it should fail with specific gofmtproblem.go complaint
func TestGoFmt(t *testing.T) {
	assert := assert.New(t)

	// Test "make gofmt
	v, err := exec.Command("make", "gofmt").Output()
	assert.Error(err) // We should have an error with gfmtproblem.go
	assert.Contains(string(v), "gofmtproblem.go")

	// Test "make SRC_DIRS=pkg/clean gofmt" - has no errors
	v, err = exec.Command("make", "SRC_DIRS=pkg/clean", "gofmt").Output()
	assert.NoError(err) // We should have an error with gfmtproblem.go

}

// Use govendor with extra and missing items
func TestGovendor(t *testing.T) {
	assert := assert.New(t)

	simpleExtraPackage := "golang.org/x/net/context"
	neededPackage := "github.com/stretchr/testify/assert"

	// Test "make govendor"
	_, err := exec.Command("make", "govendor").Output()
	assert.NoError(err) // Base code should have no errors

	// Add an unused vendor item (net/context) and check that our govendor now fails
	_, err = exec.Command("govendor", "fetch", simpleExtraPackage).Output()
	assert.NoError(err)

	v, err := exec.Command("make", "govendor").Output()
	assert.Error(err) // We should have an error now, with unused item
	assert.Contains(string(v), "u " + simpleExtraPackage)

	// Remove the extra item
	_, err = exec.Command("govendor", "remove", simpleExtraPackage).Output()
	assert.NoError(err)
	// Test "make govendor" - should be back to no errors
	_, err = exec.Command("make", "govendor").Output()
	assert.NoError(err) // Base code should have no errors

	// Remove a necessary package
	_, err = exec.Command("govendor", "remove", neededPackage).Output()
	assert.NoError(err)
	// Test "make govendor" - should show assert as a missing item
	v, err = exec.Command("make", "govendor").Output()
	assert.Error(err)
	assert.Contains(string(v), "m " + neededPackage)

	_, err = exec.Command("govendor", "fetch", neededPackage).Output()
	assert.NoError(err)

}


// Test golint on clean and unclean code
func TestGoLint(t *testing.T) {
	assert := assert.New(t)

	// Test "make golint"
	v, err := exec.Command("make", "golint").Output()
	assert.Error(err) // Should have one complaint about gofmtproblem.go
	assert.Contains(string(v), "exported function SomeExportedFunction should have comment")

	// Test "make SRC_DIRS=pkg golint" to limit to just clean directories
	v, err = exec.Command("make",  "SRC_DIRS=pkg", "golint").Output()
	assert.NoError(err) // Should have one complaint about gofmtproblem.go
}