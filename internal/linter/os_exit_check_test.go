package linter

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheckAnalyzer(t *testing.T) {
	t.Run("проверка_в_пакете_main_функции_main", func(t *testing.T) {
		analysistest.Run(t, analysistest.TestData(), OsExitCheck, "./main")
	})
	t.Run("проверка_в_пакете_pkg2_функции_main", func(t *testing.T) {
		analysistest.Run(t, analysistest.TestData(), OsExitCheck, "./pkg2")
	})
}
