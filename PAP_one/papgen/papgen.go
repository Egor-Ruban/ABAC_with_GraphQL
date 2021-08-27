package papgen

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
)

func New(filename string, typename string, cfg *config.Config) plugin.Plugin {
	//тут можно добавить какие-то данные, например конфиг
	return &Plugin{filename: filename, typeName: typename, cfg: cfg}
}

type Plugin struct {
	filename string
	typeName string
	cfg      *config.Config
}

var _ plugin.CodeGenerator = &Plugin{}
var _ plugin.ConfigMutator = &Plugin{}

func (m *Plugin) Name() string {
	return "papgen"
}

func (m *Plugin) MutateConfig(cfg *config.Config) error {
	_ = syscall.Unlink("graph/schema.resolvers.go")
	_ = syscall.Unlink("graph/resolver.go")
	return nil
}

type Rules struct {
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	_, err := filepath.Abs(m.filename)
	if err != nil {
		return err
	}
	//pkgName := NameForDir(filepath.Dir(abs))
	for s, d := range data.Schema.Types {
		if s == "Rules" {
			for _, f := range d.Fields {
				fmt.Println(f.Name)
				for _, d := range f.Directives {
					//todo берём название, сплитим по __, 1-я часть - "condition" (иначе ошибка, наверн).
					// 2-я часть - название атрибута, заменяем '_' на '.'
					// 3-я часть как сравнивать. заменяем is на =, over на >, below на <
					// со временем всё ещё не знаю как работать
					fmt.Println("\t", d.Name)
					for _, d := range d.Arguments {
						fmt.Println("\t\t", d.Name, "+", d.Value.Raw)
					}
				}
			}

		}
	}
	return templates.Render(templates.Options{
		PackageName: "graph",
		Filename:    "graph/schema.resolvers.go",
		Data: &ResolverBuild{
			Data:     data,
			TypeName: m.typeName,
		},
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

type ResolverBuild struct {
	*codegen.Data

	TypeName string
}

// NameForDir manually looks for package stanzas in files located in the given directory. This can be
// much faster than having to consult go list, because we already know exactly where to look.
func NameForDir(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return SanitizePackageName(filepath.Base(dir))
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return SanitizePackageName(filepath.Base(dir))
	}
	fset := token.NewFileSet()
	for _, file := range files {
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".go") {
			continue
		}

		filename := filepath.Join(dir, file.Name())
		if src, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly); err == nil {
			return src.Name.Name
		}
	}

	return SanitizePackageName(filepath.Base(dir))
}

var invalidPackageNameChar = regexp.MustCompile(`[^\w]`)

func SanitizePackageName(pkg string) string {
	return invalidPackageNameChar.ReplaceAllLiteralString(filepath.Base(pkg), "_")
}
