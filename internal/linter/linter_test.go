package linter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
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

func TestInitAnalyzers(t *testing.T) {
	tests := []struct {
		name          string
		jsonData      []byte
		initFunc      func(*StaticlintConfig) []*analysis.Analyzer
		expectedNames []string
	}{
		{
			name:     "InitAnalysis",
			jsonData: []byte(`{"staticcheck": [], "analysis": ["appends", "asmdecl"]}`),
			initFunc: func(config *StaticlintConfig) []*analysis.Analyzer {
				return config.InitAnalysis()
			},
			expectedNames: []string{"appends", "asmdecl"},
		},
		{
			name:     "InitStaticCheck",
			jsonData: []byte(`{"staticcheck": ["SA1000", "SA1001"], "analysis": []}`),
			initFunc: func(config *StaticlintConfig) []*analysis.Analyzer {
				return config.InitStaticCheck()
			},
			expectedNames: []string{"SA1000", "SA1001"},
		},
		{
			name:     "InitOtherCheck",
			jsonData: []byte(`{"staticcheck": [], "analysis": [], "other": ["errcheck", "restrictpkg"]}`),
			initFunc: func(config *StaticlintConfig) []*analysis.Analyzer {
				return config.InitOtherCheck()
			},
			expectedNames: []string{"errcheck", "restrictpkg"},
		},
		{
			name:     "InitOsExitCheck",
			jsonData: []byte(`{"staticcheck": [], "analysis": [], "other": ["osexit"]}`),
			initFunc: func(config *StaticlintConfig) []*analysis.Analyzer {
				return config.InitOsExitCheck()
			},
			expectedNames: []string{"osexit"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewStaticlintConfig(tt.jsonData)
			err := config.FilConfig()
			require.NoError(t, err)
			analyzers := tt.initFunc(config)
			require.Len(t, analyzers, len(tt.expectedNames))
			for i, expectedName := range tt.expectedNames {
				assert.Equal(t, expectedName, analyzers[i].Name)
			}
		})
	}
}
