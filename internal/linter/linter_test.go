package linter

import (
	"testing"

	"github.com/cybozu-go/golang-custom-analyzer/pkg/restrictpkg"
	"github.com/kisielk/errcheck/errcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
)

func TestNewStaticlintConfig(t *testing.T) {
	jsonData := []byte(`{"staticcheck": ["SA1000"], "analysis": ["appends"]}`)
	config := NewStaticlintConfig(jsonData)
	require.NotNil(t, config)
	assert.Equal(t, jsonData, config.jsonData)
}

func TestFilConfig(t *testing.T) {
	jsonData := []byte(`{"staticcheck": ["SA1000"], "analysis": ["appends"]}`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.NoError(t, err)
	assert.Equal(t, Config{Staticcheck: []string{"SA1000"}, Analysis: []string{"appends"}}, config.cfg)
}

func TestFilConfig_InvalidJSON(t *testing.T) {
	jsonData := []byte(`invalid json`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "конфигурация не разобрана")
}

func TestInitAnalysis(t *testing.T) {
	jsonData := []byte(`{"staticcheck": [], "analysis": ["appends", "asmdecl"]}`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.NoError(t, err)

	analyzers := config.InitAnalysis()
	require.Len(t, analyzers, 2)
	assert.Equal(t, appends.Analyzer, analyzers[0])
	assert.Equal(t, asmdecl.Analyzer, analyzers[1])
}

func TestInitStaticCheck(t *testing.T) {
	jsonData := []byte(`{"staticcheck": ["SA1000", "SA1001"], "analysis": []}`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.NoError(t, err)

	analyzers := config.InitStaticCheck()
	require.Len(t, analyzers, 2)
	// Assuming the staticcheck analyzers have the names "SA1000" and "SA1001"
	assert.Equal(t, "SA1000", analyzers[0].Name)
	assert.Equal(t, "SA1001", analyzers[1].Name)
}

func TestInitOtherCheck(t *testing.T) {
	jsonData := []byte(`{"staticcheck": [], "analysis": []}`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.NoError(t, err)

	analyzers := config.InitOtherCheck()
	require.Len(t, analyzers, 2)
	assert.Equal(t, errcheck.Analyzer, analyzers[0])
	assert.Equal(t, restrictpkg.RestrictPackageAnalyzer, analyzers[1])
}

func TestInitOsExitCheck(t *testing.T) {
	jsonData := []byte(`{"staticcheck": [], "analysis": []}`)
	config := NewStaticlintConfig(jsonData)
	err := config.FilConfig()
	require.NoError(t, err)

	analyzers := config.InitOsExitCheck()
	require.Len(t, analyzers, 1)
	// Assuming the OsExitCheck analyzer is registered
	assert.Equal(t, OsExitCheck, analyzers[0])
}
