# cmd/staticlint

Статический анализатор кода

### config.json

 - Секция staticcheck содержит коды проверок https://staticcheck.dev/docs/

 - Секция analysis содержит набор проверок https://pkg.go.dev/golang.org/x/tools/go/analysis/passes

### Дополнительные проверки
 - github.com/cybozu-go/golang-custom-analyzer/pkg/restrictpkg
 - github.com/kisielk/errcheck/errcheck

### Проверка на наличие вызова os.Exit в функции main пакета main