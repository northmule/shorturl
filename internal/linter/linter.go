package linter

import (
	"encoding/json"
	"fmt"

	"github.com/cybozu-go/golang-custom-analyzer/pkg/restrictpkg"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/staticcheck"
)

// StaticlintConfig конфигурация линтера.
type StaticlintConfig struct {
	jsonData []byte
	cfg      Config
}

// Config структура конфигурации
type Config struct {
	Staticcheck []string `json:"staticcheck"`
	Analysis    []string `json:"analysis"`
}

// NewStaticlintConfig конструктор
func NewStaticlintConfig(jsonData []byte) *StaticlintConfig {
	instance := &StaticlintConfig{
		jsonData: jsonData,
	}
	return instance
}

// FilConfig читает конфиг, заполняет структуру
func (s *StaticlintConfig) FilConfig() error {
	var err error
	var cfg Config
	err = json.Unmarshal(s.jsonData, &cfg)
	if err != nil {
		return fmt.Errorf("конфигурация не разобрана :%s", err)
	}
	s.cfg = cfg
	return nil
}

// InitAnalysis Конфигурирование analysis
func (s *StaticlintConfig) InitAnalysis() []*analysis.Analyzer {
	defaultAnalyzer := map[string]*analysis.Analyzer{
		"appends":          appends.Analyzer,
		"asmdecl":          asmdecl.Analyzer,
		"assign":           assign.Analyzer,
		"atomic":           atomic.Analyzer,
		"bools":            bools.Analyzer,
		"buildtag":         buildtag.Analyzer,
		"cgocall":          cgocall.Analyzer,
		"composite":        composite.Analyzer,
		"copylock":         copylock.Analyzer,
		"directive":        directive.Analyzer,
		"errorsas":         errorsas.Analyzer,
		"framepointer":     framepointer.Analyzer,
		"httpresponse":     httpresponse.Analyzer,
		"ifaceassert":      ifaceassert.Analyzer,
		"loopclosure":      loopclosure.Analyzer,
		"lostcancel":       lostcancel.Analyzer,
		"nilfunc":          nilfunc.Analyzer,
		"printf":           printf.Analyzer,
		"shift":            shift.Analyzer,
		"sigchanyzer":      sigchanyzer.Analyzer,
		"stdmethods":       stdmethods.Analyzer,
		"stringintconv":    stringintconv.Analyzer,
		"structtag":        structtag.Analyzer,
		"tests":            tests.Analyzer,
		"testinggoroutine": testinggoroutine.Analyzer,
		"timeformat":       timeformat.Analyzer,
		"unmarshal":        unmarshal.Analyzer,
		"unreachable":      unreachable.Analyzer,
		"unsafeptr":        unsafeptr.Analyzer,
		"unusedresult":     unusedresult.Analyzer,
	}

	var mychecks []*analysis.Analyzer
	for _, cfgValue := range s.cfg.Analysis {
		if analyzer, ok := defaultAnalyzer[cfgValue]; ok {
			mychecks = append(mychecks, analyzer)
		}
	}

	return mychecks
}

// InitStaticCheck конфигурирование staticcheck
func (s *StaticlintConfig) InitStaticCheck() []*analysis.Analyzer {
	var mychecks []*analysis.Analyzer
	checks := make(map[string]bool)
	for _, v := range s.cfg.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	return mychecks
}

// InitOtherCheck прочие анализаторы по ТЗ
func (s *StaticlintConfig) InitOtherCheck() []*analysis.Analyzer {
	mychecks := []*analysis.Analyzer{
		errcheck.Analyzer,
		restrictpkg.RestrictPackageAnalyzer,
	}
	return mychecks
}

// InitOsExitCheck поис вызова os.Exit в main
func (s *StaticlintConfig) InitOsExitCheck() []*analysis.Analyzer {
	mychecks := []*analysis.Analyzer{
		OsExitCheck,
	}
	return mychecks
}
