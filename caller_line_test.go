package zentests

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

const callerLineProbeEnv = "ZENTESTS_CALLER_LINE_PROBE"

// TestCallerLineReporting runs a failing assertion in a subprocess and checks that its reported
// location is the test call site.
func TestCallerLineReporting(t *testing.T) {
	if os.Getenv(callerLineProbeEnv) == "1" {
		app := fiber.New()
		New(t).Describe("request", func(ctx *SuiteContext) {
			ctx.Describe("status", func(ctx *SuiteContext) {
				ctx.It("reports the call site", func(zt *T) {
					zt.Get(app, "/missing").OK()
				})
			})
		})
		return
	}

	//nolint:gosec // G204: runs the current test binary with a fixed selector.
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestCallerLineReporting$")
	cmd.Env = append(os.Environ(), callerLineProbeEnv+"=1")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected the probe assertion to fail")
	}

	failureLocation := regexp.MustCompile(`(?m)^\s*caller_line_test\.go:\d+:`)
	if !failureLocation.Match(output) || !strings.Contains(string(output), "status code mismatch") {
		t.Fatalf("failure did not report the caller's test line:\n%s", output)
	}
}
