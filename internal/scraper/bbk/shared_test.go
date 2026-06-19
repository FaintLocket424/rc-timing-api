package bbk

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// collectAndSortTestFiles recursively searches for files in the testdata directory
// matching the predicate function and returns them sorted alphabetically.
func collectAndSortTestFiles(t *testing.T, matches func(os.DirEntry) bool) []string {
	var files []string
	err := filepath.WalkDir("testdata", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && matches(d) {
			files = append(files, path)
		}
		return nil
	})
	require.NoError(t, err)

	sort.Strings(files)
	return files
}

// runGoldenTest verifies that the parsed HTML output from parseFunc matches
// the manually verified golden JSON data.
func runGoldenTest[T any](
	t *testing.T,
	htmlPath string,
	parseFunc func(io.Reader) (*T, error),
	normaliseFunc func(t *testing.T, expected, actual *T),
) {
	jsonPath := strings.TrimSuffix(htmlPath, filepath.Ext(htmlPath)) + ".json"

	t.Run(htmlPath, func(t *testing.T) {
		if _, err := os.Stat(filepath.Clean(jsonPath)); os.IsNotExist(err) {
			t.Skipf("No golden file found: %s", jsonPath)
		}

		htmlFile, err := os.Open(filepath.Clean(htmlPath))
		require.NoError(t, err)
		defer func() { _ = htmlFile.Close() }()

		actualData, err := parseFunc(htmlFile)
		require.NoError(t, err)

		expectedJSON, err := os.ReadFile(filepath.Clean(jsonPath))
		require.NoError(t, err)

		var expectedData T
		decoder := json.NewDecoder(bytes.NewReader(expectedJSON))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&expectedData)
		require.NoError(t, err, "Golden JSON is invalid or contains unknown fields")

		if normaliseFunc != nil {
			normaliseFunc(t, &expectedData, actualData)
		}

		assert.Equal(t, &expectedData, actualData, "Parsed data does not match manually verified golden file")
	})
}
