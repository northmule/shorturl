package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/northmule/shorturl/internal/linter"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	var analysisList []*analysis.Analyzer
	appFile, err := os.Executable()
	if err != nil {
		log.Fatalf("ошибка %s", err)
		return
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appFile), "config.json"))
	if err != nil {
		log.Fatalf("ошибка чтения файла конфигурации: %s", err)
		return
	}
	lint := linter.NewStaticlintConfig(data)
	err = lint.FilConfig()
	if err != nil {
		log.Fatal(err)
	}

	analysisList = append(analysisList, lint.InitAnalysis()...)
	analysisList = append(analysisList, lint.InitStaticCheck()...)
	analysisList = append(analysisList, lint.InitOtherCheck()...)
	analysisList = append(analysisList, lint.InitOsExitCheck()...)

	multichecker.Main(
		analysisList...,
	)
}
