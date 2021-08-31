package papgen

import (
	"errors"
	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/99designs/gqlgen/plugin"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
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
	// todo можно добавить название ресолвера, который будет содержать правила
	Rules []Rule
}

type Rule struct {
	Mode       string
	Conditions []Condition
}

type Condition struct {
	Compare string
	To      string
	With    string
}

func (m *Plugin) GenerateCode(data *codegen.Data) error {
	_, err := filepath.Abs(m.filename)
	if err != nil {
		return err
	}
	//pkgName := NameForDir(filepath.Dir(abs))
	r := Rules{Rules: []Rule{}}
	for s, d := range data.Schema.Types {
		if s == "Rules" {
			for _, f := range d.Fields {
				rule := Rule{f.Type.Name(), []Condition{}}
				for _, d := range f.Directives {
					c := Condition{}
					cond := strings.Split(d.Name, "__")
					if cond[0] != "condition" {
						return errors.New("wrong rule annotation: " + d.Name)
					}
					compare := strings.Split(cond[1], "_")
					for c := range compare {
						compare[c] = strings.Title(compare[c])
						if strings.ToLower(compare[c]) == "id" {
							compare[c] = "ID"
						}
					}
					cond[1] = strings.Join(compare, ".")
					tr := map[string]string{"over": " > ", "below": " < ", "is": " == ", "not": "!="}

					for k, v := range tr {
						cond[2] = strings.ReplaceAll(cond[2], k, v)
					}
					c.Compare = cond[1]
					c.With = cond[2]
					// берём название, сплитим по __, 1-я часть - "condition" (иначе ошибка, наверн).
					// 2-я часть - название атрибута, заменяем '_' на '.'
					// 3-я часть как сравнивать. заменяем is на =, over на >, below на <
					// со временем всё ещё не знаю как работать
					// f.Value.Raw сплитим по __, Если первая часть - "attr", то сравнивать надо с атрибутом
					if len(d.Arguments) != 1 {
						return errors.New("too much arguments for rule annotation: " + d.Name)
					}
					arg := strings.Split(d.Arguments[0].Value.Raw, "__")
					if strings.Contains(arg[0], "attr") {
						compare := strings.Split(arg[1], "_")
						for c := range compare {
							compare[c] = strings.Title(compare[c])
							if strings.ToLower(compare[c]) == "id" {
								compare[c] = "ID"
							}
						}
						c.To = "attrs." + strings.Join(compare, ".")
					} else {
						if strings.Contains(d.Arguments[0].Name, "int") {
							c.To = arg[0]
						} else if strings.Contains(d.Arguments[0].Name, "id") {
							c.To = "\"" + arg[0] + "\""
						} else if strings.Contains(d.Arguments[0].Name, "str") {
							c.To = "\"" + arg[0] + "\""
						} else if strings.Contains(d.Arguments[0].Name, "time") {
							c.Compare = c.Compare + ".Weekday().String()"
							c.To = "\"" + arg[0] + "\""
						}

					}
					rule.Conditions = append(rule.Conditions, c)
					//fmt.Println("\t", d.Name + "(" + d.Arguments[0].Value.Raw + ")" + " is " + c.Compare +c.With + c.To)
				}
				r.Rules = append(r.Rules, rule)
			}
		}
	}

	return templates.Render(templates.Options{
		PackageName: "graph",
		Filename:    "graph/schema.resolvers.go",
		Data: &ResolverBuild{
			Data:     data,
			TypeName: m.typeName,
			Rules:    r,
		},
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

type ResolverBuild struct {
	*codegen.Data
	Rules    Rules
	TypeName string
}

// NameForDir manually looks for package stanzas in files located in the given directory. This can be
// much faster than having To consult go list, because we already know exactly where To look.
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
