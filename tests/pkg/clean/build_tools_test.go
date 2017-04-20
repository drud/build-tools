package clean

import (
	"os/exec"
	"testing"

	"fmt"
	"log"
	"os"
	"strings"

	"github.com/drud/build-tools/tests/pkg/version"
	"github.com/stretchr/testify/assert"
)

var (
	osname = "" // The operating system.
)

func init() {
	// Default directory is the directory of the test file, but we need to run make from the directory of the Makefile
	err := os.Chdir("../..")
	if err != nil {
		log.Fatalln("Failed to chdir to ../..", err)
	}
	// Operating system - Darwin or Linux
	v, err := exec.Command("uname", "-s").Output()
	if err != nil {
		log.Fatalln("Failed to run uname command:", string(v))
	}
	osname = strings.TrimSpace(string(v))
}

// Runs a number of standard make targets and test for basic sanity of result
// Assumes operation in the "testing" directory where the Makefile is
func TestBuild(t *testing.T) {
	a := assert.New(t)

	// Map OS name to output location
	binlocs := map[string]string{
		"Darwin": "bin/darwin/darwin_amd64",
		"Linux":  "bin/linux",
	}

	dir, _ := os.Getwd()
	fmt.Println("Current Directory:", dir)

	v, err := exec.Command("which", "make").Output()
	a.NoError(err)
	a.Contains(string(v), "make")

	// Try trivial "make version". This does use local system's make and git commands
	v, err = exec.Command("make", "version").Output()
	a.NoError(err)
	a.Contains(string(v), "VERSION:"+version.VERSION)
	if err != nil {
		log.Fatalln("make version in", dir, "failed, so exiting. output=", string(v))
	}

	// Run a make clean to start with; linux requires sudo because container left things a mess
	v, err = exec.Command("make", "clean").Output()
	a.NoError(err, "make clean failed. output="+string(v))

	// Build darwin and linux cmds
	v, err = exec.Command("make", "darwin").Output()
	a.NoError(err, "Failed to 'make darwin'")
	a.Contains(string(v), "building darwin")

	v, err = exec.Command("make", "linux").Output()
	a.NoError(err)
	a.Contains(string(v), "building linux")

	// Run the native gofmtproblem application to make sure it runs
	v, err = exec.Command(binlocs[osname] + "/build_tools_dummy").Output()
	a.NoError(err)
	a.Contains(string(v), "This is build_tools_dummy.go")
	a.Contains(string(v), version.VERSION)

	// Make container
	v, err = exec.Command("make", "container").Output()
	a.NoError(err)
	a.Contains(string(v), "Successfully built")

}

// Try gofmt - it should fail with specific gofmtproblem.go complaint
func TestGoFmt(t *testing.T) {
	assert := assert.New(t)

	// Test "make gofmt
	v, err := exec.Command("make", "gofmt").Output()
	assert.Error(err) // We should have an error with bad_gofmt_code.go
	assert.Contains(string(v), "pkg/dirtyComplex/bad_gofmt_code.go")

	// Test "make SRC_DIRS=pkg/clean gofmt" - has no errors
	v, err = exec.Command("make", "SRC_DIRS=pkg/clean", "gofmt").Output()
	assert.NoError(err, "make SRC_DIRS=pkg/clean gofmt failed: output="+string(v)) // No error on the clean directory

}

// Use govendor with extra and missing items
// There is some danger here that failure can leave the vendor directory with changed items.
// Note that the spare make target "make container_cmd" is used as a generic way to execute govendor in the container.
// The COMMAND=some command argument to exec.Command() is an oddity - due to the way this argument is processed,
// it must not be escaped.
func TestGovendor(t *testing.T) {
	assert := assert.New(t)

	simpleExtraPackage := "golang.org/x/net/context"
	neededPackage := "github.com/stretchr/testify/assert"

	// Test "make govendor"
	_, err := exec.Command("make", "govendor").Output()
	assert.NoError(err) // Base code should have no errors

	// Add an unused vendor item (net/context) and check that our govendor now fails
	v, err := exec.Command("make", "COMMAND=govendor fetch "+simpleExtraPackage, "container_cmd").Output()
	assert.NoError(err, "Failed 'govendor fetch %v', result=%v", simpleExtraPackage, string(v))

	v, err = exec.Command("make", "govendor").Output()
	assert.Error(err) // We should have an error now, with unused item
	assert.Contains(string(v), "u "+simpleExtraPackage)

	// Remove the extra item
	_, err = exec.Command("make", "COMMAND=govendor remove "+simpleExtraPackage, "container_cmd").Output()
	assert.NoError(err)
	// Test "make govendor" - should be back to no errors
	_, err = exec.Command("make", "govendor").Output()
	assert.NoError(err) // Base code should have no errors

	// Remove a necessary package
	_, err = exec.Command("make", "COMMAND=govendor remove "+neededPackage, "container_cmd").Output()
	assert.NoError(err)
	// Test "make govendor" - should show assert as a missing item
	v, err = exec.Command("make", "govendor").Output()
	assert.Error(err)
	assert.Contains(string(v), "m "+neededPackage)

	_, err = exec.Command("make", "COMMAND=govendor fetch "+neededPackage, "container_cmd").Output()
	assert.NoError(err)

}

// Test golint on clean and unclean code
func TestGoLint(t *testing.T) {
	assert := assert.New(t)

	// Test "make golint"
	v, err := exec.Command("make", "golint").Output()
	assert.Error(err) // Should have one complaint about gofmtproblem.go
	assert.Contains(string(v), "exported function DummyExported_function should have comment")

	// Test "make SRC_DIRS=pkg golint" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "golint").Output()
	assert.NoError(err) // Should have one complaint about gofmtproblem.go
}

// Test govet for simple problems
func TestGoVet(t *testing.T) {
	assert := assert.New(t)

	// cmd/gofmtproblem/gofmtproblem.go
	// Test "make govet"
	v, err := exec.Command("make", "govet").Output()
	assert.Error(err) // Should have one complaint about gofmtproblem.go
	assert.Contains(string(v), "pkg/dirtyComplex/bad_govet_code.go")

	// Test "make SRC_DIRS=pkg govet" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "govet").Output()
	assert.NoError(err) // Should have no complaints in clean package
}

// Test errcheck.
func TestErrCheck(t *testing.T) {
	assert := assert.New(t)

	// pkg/dirtycomplex/bad_errcheck_code.go
	// Test "make errcheck"
	v, err := exec.Command("make", "errcheck").Output()
	assert.Error(err) // Should have one complaint about bad_errcheck_code.go
	assert.Contains(string(v), "pkg/dirtyComplex/bad_errcheck_code.go")

	// Test "make SRC_DIRS=pkg errcheck" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "errcheck").Output()
	assert.NoError(err) // Should have no complaints in clean package
}

// Test staticcheck.
func TestStaticcheck(t *testing.T) {
	assert := assert.New(t)

	// Test "make staticcheck"
	v, err := exec.Command("make", "staticcheck").Output()
	assert.Error(err) // Should have one complaint about bad_staticcheck_code.go
	assert.Contains(string(v), "pkg/dirtyComplex/bad_staticcheck_code.go")

	// Test "make SRC_DIRS=pkg/clean staticcheck" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "staticcheck").Output()
	assert.NoError(err) // Should have no complaints in clean package
}

// Test unused.
func TestUnused(t *testing.T) {
	assert := assert.New(t)

	// Test "make unused"
	v, err := exec.Command("make", "unused").Output()
	assert.Error(err) // Should have one complaint about bad_unused_code.go
	assert.Contains(string(v), "pkg/dirtyComplex/bad_unused_code.go")

	// Test "make SRC_DIRS=pkg/clean unused" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "unused").Output()
	assert.NoError(err) // Should have no complaints in clean package
}

// Test codecoroner.
func TestCodeCoroner(t *testing.T) {
	assert := assert.New(t)

	// Test "make unused"
	v, err := exec.Command("make", "codecoroner").Output()
	assert.Error(err)                                        // Should complain about pretty much everything in the dirtyComplex package.
	assert.Contains(string(v), "AnotherExportedFunction")    // Check an exported function
	assert.Contains(string(v), "yetAnotherExportedFunction") // Check an unexported function.

	// Test "make SRC_DIRS=pkg/clean unused" to limit to just clean directories
	_, err = exec.Command("make", "SRC_DIRS=pkg/clean", "unused").Output()
	assert.NoError(err) // Should have no complaints in clean package
}
