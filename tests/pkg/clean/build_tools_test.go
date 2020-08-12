package clean

import (
	"os/exec"
	"testing"

	"log"
	"os"

	"runtime"

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
	// Operating system - Darwin or Linux or Windows
	osname = runtime.GOOS
}

// Runs a number of standard make targets and test for basic sanity of result
// Assumes operation in the "testing" directory where the Makefile is
func TestBuild(t *testing.T) {
	a := assert.New(t)

	// Map OS name to output location
	binlocs := map[string]string{
		"darwin":  ".gotmp/bin/darwin_amd64/build_tools_dummy",
		"linux":   ".gotmp/bin/build_tools_dummy",
		"windows": ".gotmp/bin/windows_amd64/build_tools_dummy.exe",
	}

	dir, _ := os.Getwd()

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

	// Make sure it builds from scratch, but don't delete our pkg cache
	v, err = exec.Command("rm", "-f", "darwin", "linux", "windows").Output()
	a.NoError(err, "error output from rm: %v", string(v))

	// Build darwin and linux cmds
	v, err = exec.Command("bash", "-c", "pwd && make darwin").Output()
	a.NoError(err, "Failed to 'make darwin'")
	a.Contains(string(v), "building darwin")

	v, err = exec.Command("bash", "-c", "pwd && make linux").Output()
	a.NoError(err, "failed 'make linux', err=%v, output='%v'", err, string(v))
	a.Contains(string(v), "building linux")

	v, err = exec.Command("bash", "-c", "pwd && make windows").Output()
	a.NoError(err)
	a.Contains(string(v), "building windows")

	// Run the native gofmtproblem application to make sure it runs
	v, err = exec.Command(binlocs[osname]).Output()
	a.NoError(err)
	a.Contains(string(v), "This is build_tools_dummy.go")
	a.Contains(string(v), version.VERSION)
	a.Contains(string(v), "Built ")
	a.NotContains(string(v), "COMMIT should be overridden")
	a.NotContains(string(v), "BUILDINFO should have new info")

	// Make container
	v, err = exec.Command("make", "container").Output()
	a.NoError(err)
	a.Contains(string(v), "Successfully built")

}

// Try gofmt - it should fail with specific gofmtproblem.go complaint
func TestGoFmt(t *testing.T) {
	a := assert.New(t)

	// Test "make gofmt
	v, err := exec.Command("bash", "-c", "pwd && make gofmt").Output()
	a.Error(err) // We should have an error with bad_gofmt_code.go
	a.Contains(string(v), "pkg/dirtyComplex/bad_gofmt_code.go")

	// Test "make SRC_DIRS=pkg/clean gofmt" - has no errors
	v, err = exec.Command("make", "SRC_DIRS=pkg/clean", "gofmt").Output()
	a.NoError(err, "make SRC_DIRS=pkg/clean gofmt failed: output="+string(v)) // No error on the clean directory

}

// Test golint on clean and unclean code
func TestGoLint(t *testing.T) {
	a := assert.New(t)

	// Test "make golint"
	v, err := exec.Command("bash", "-c", "pwd && make golint").Output()
	a.Error(err) // Should have one complaint about gofmtproblem.go
	a.Contains(string(v), "exported function DummyExported_function should have comment")

	// Test "make SRC_DIRS=pkg golint" to limit to just clean directories
	_, err = exec.Command("make", "golint", "SRC_DIRS=pkg/clean").Output()
	a.NoError(err) // Should have one complaint about gofmtproblem.go
}

// Test govet for simple problems
func TestGoVet(t *testing.T) {
	a := assert.New(t)

	// cmd/gofmtproblem/gofmtproblem.go
	// Test "make govet"
	v, err := exec.Command("bash", "-c", "pwd && make govet").Output()
	a.Error(err) // Should have one complaint about gofmtproblem.go
	a.Contains(string(v), "pkg/dirtyComplex/bad_govet_code.go")

	// Test "make SRC_DIRS=pkg govet" to limit to just clean directories
	v, err = exec.Command("bash", "-c", "pwd && make govet SRC_DIRS=pkg/clean").Output()
	a.NoError(err, "error output from make govet: %v", string(v)) // Should have no complaints in clean package
}

// Test errcheck.
func TestErrCheck(t *testing.T) {
	a := assert.New(t)

	// pkg/dirtycomplex/bad_errcheck_code.go
	// Test "make errcheck"
	v, err := exec.Command("bash", "-c", "pwd && make errcheck").Output()
	a.Error(err) // Should have one complaint about bad_errcheck_code.go
	a.Contains(string(v), "pkg/dirtyComplex/bad_errcheck_code.go")

	// Test "make SRC_DIRS=pkg errcheck" to limit to just clean directories
	_, err = exec.Command("bash", "-c", "pwd && make errcheck SRC_DIRS=pkg/clean").Output()
	a.NoError(err) // Should have no complaints in clean package
}

// Test misspell.
func TestMisspell(t *testing.T) {
	a := assert.New(t)

	// Test "make codecoroner"
	v, err := exec.Command("bash", "-c", "make --no-print-directory misspell").Output()
	a.NoError(err)                                               // This one doesn't make an error return
	a.Contains(string(v), " is a misspelling of \"misspelled\"") // Check an exported function

	// Test "make SRC_DIRS=pkg/clean misspell" to limit to just clean directories
	v, err = exec.Command("bash", "-c", "make --no-print-directory misspell SRC_DIRS=pkg/clean").Output()
	a.NoError(err) // Should have no complaints in clean package
	a.Equal("Checking for misspellings: \n", string(v))
}

// Test golangci-lint.
func TestGolangciLint(t *testing.T) {
	a := assert.New(t)

	// Test "make golangci-lint"
	v, err := exec.Command("bash", "-c", "make golangci-lint").Output()
	a.Error(err) // Should complain about pretty much everything in the dirtyComplex package.
	execResult := string(v)
	a.Contains(execResult, "don't use MixedCaps in package name; dirtyComplex should be dirtycomplex")
	a.Contains(execResult, "don't use underscores in Go names; func DummyExported_function should be DummyExportedFunction (golint)")
	a.Contains(execResult, "File is not `gofmt`-ed with `-s`")
	a.Contains(execResult, "ineffectual assignment to `num` (ineffassign)")
	a.Contains(execResult, "yetAnotherExportedFunction` is unused (deadcode)")
	a.Contains(execResult, "Error return value of `os.Chown` is not checked (errcheck)")

	out, err := exec.Command("bash", "-c", "make golangci-lint SRC_DIRS=pkg/clean").Output()
	a.NoError(err, "Failed to get clean result for golangci-lint: %v (output=%s)", err, out) // Should have no complaints in clean package
}
