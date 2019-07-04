package main

import (
	"flag"
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/afero"
	"log"
	"strings"
	"text/template"
)

type Metadata struct {
	TypeName        string
	TypeNamePlural  string
	PackageName     string
	PackageReceiver string
}

func main() {
	var typeName string
	var fs afero.Fs

	templates := []string{
		"init",
		"requests",
		"response",
		"resource",
	}

	p := pluralize.NewClient()
	fs = afero.NewOsFs()

	box := packr.New("templates", "./templates")

	flag.StringVar(&typeName, "type", "", "name of the type for which the code is to be generated")

	flag.Parse()

	meta := Metadata{
		TypeName:        strings.Title(typeName),
		TypeNamePlural:  p.Plural(strings.Title(typeName)),
		PackageName:     p.Plural(strings.ToLower(typeName)),
		PackageReceiver: strings.ToLower(typeName[:1]),
	}

	funcs := template.FuncMap{
		"lc": func(s string) string {
			return strings.ToLower(s)
		},
	}

	for _, tpl := range templates {
		tStr, err := box.FindString(fmt.Sprintf("%v.gotpl", tpl))
		if err != nil {
			log.Fatal(err)
		}

		t := template.Must(template.New(fmt.Sprintf("%v_template", tpl)).Funcs(funcs).Parse(tStr))

		file, err := fs.Create(fmt.Sprintf("%v_%v.go", meta.PackageName, tpl))
		if err != nil {
			log.Fatal(err)
		}

		err = t.Execute(file, meta)
		if err != nil {
			log.Fatal(err)
		}
	}
}
