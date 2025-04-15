package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func skipIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") != "1" {
		t.Skip("Skipping integration tests")
	}
}

func TestCLI(t *testing.T) {
	skipIntegration(t)

	tempDir, err := os.MkdirTemp("", "")
	require.NoErrorf(t, err, "error creating temporary directory: %s", err)

	binPath := filepath.Join(tempDir, "cueitup")
	buildArgs := []string{"build", "-o", binPath, "../.."}

	c := exec.Command("go", buildArgs...)
	err = c.Run()
	require.NoErrorf(t, err, "error building binary: %s", err)

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			fmt.Printf("couldn't clean up temporary directory (%s): %s", binPath, err)
		}
	}()

	t.Run("Showing help works", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "-h")
		b, err := c.CombinedOutput()

		// THEN
		assert.NoError(t, err, "output:\n%s", b)
	})

	t.Run("Validate config", func(t *testing.T) {
		// GIVEN
		// WHEN
		c := exec.Command(binPath, "config", "validate", "-c", "static/config-bad.yml")
		outputBytes, err := c.CombinedOutput()

		// THEN
		require.Error(t, err)
		require.Equal(t, "exit status 1", err.Error())
		expected := `config has some errors:
- profile config is invalid at index 1
  - profile name is empty
  - queue URL is incorrect ("sqs.eu-central-1.amazonaws.com/000000000000/queue-a"): needs to be a proper URL
- profile config is invalid at index 2
  - encoding format is incorrect: "unknown"; possible values: [json, none]
  - incorrect config source provided; possible values: "env", "profile:<aws-shared-config-profile-name>", "assume:<arn-of-role-to-assume>"
- profile config is invalid at index 3
  - context key is empty
  - subset key is empty
`
		assert.Equal(t, expected, string(outputBytes))
	})
}
