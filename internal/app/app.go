package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"github.com/worldiety/gdoc/internal/generator/asciidoc"
	"github.com/worldiety/gdoc/internal/parser/golang"
	"gopkg.in/yaml.v3"
	"log"
	"strings"
)

type OutputFormat string

func (f OutputFormat) Category() string {
	switch f {
	case Yaml, Json:
		return "stdlib"
	case Adoc, Pdf:
		return "adoc"
	default:
		return ""
	}
}

const (
	Yaml = "yaml"
	Json = "json"
	Adoc = "adoc"
	Pdf  = "pdf"
	Html = "html"
)

type Config struct {
	ModPath      string
	OutputFormat string
	Packages     string
	PkgSep       string
}

func (c *Config) Reset() {
	wd, err := golang.ModWdRoot()
	if err != nil {
		log.Fatal(fmt.Errorf("could not walk to mod root directory: %w", err))
	}

	c.ModPath = wd
	c.OutputFormat = Adoc
	c.PkgSep = "/"
}

func (c *Config) Flags(flags *flag.FlagSet) {
	flags.StringVar(&c.ModPath, "modPath", c.ModPath, "the modules path to use")
	flags.StringVar(&c.OutputFormat, "format", c.OutputFormat, "default is adoc. yaml|json are available as well. "+
		"html is available, if ascidoctor is installed. "+
		"pdf is available, if asciidoctor-pdf is installed")
	flags.StringVar(&c.Packages, "packages", c.Packages, "if not empty, only scan the listed packages separated by ;")
	flags.StringVar(&c.PkgSep, "pkgSep", c.PkgSep, "sets the path separator between packages. Default is / which is not json-pointer friendly")
}

// Apply takes a Config and uses the contained instructions to generate documentation.
func Apply(cfg Config) ([]byte, error) {
	pkgs := strings.Split(cfg.Packages, ";")
	if len(pkgs) == 1 && pkgs[0] == "" {
		pkgs = nil
	}

	node, err := golang.Parse(cfg.ModPath, pkgs...)
	if err != nil {
		return nil, fmt.Errorf("cannot parse from %s: %w", cfg.ModPath, err)
	}
	// add information not available in the ast package to the module
	err = golang.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve %s: %w", node.Name, err)
	}

	tmp := map[api.ImportPath]*api.Package{}
	if cfg.PkgSep != "/" {
		for path, p := range node.Packages {
			qualifier := strings.ReplaceAll(path, "/", cfg.PkgSep)
			tmp[qualifier] = p

			var impTmp []api.Import
			for _, s := range p.Imports {
				qualifier := strings.ReplaceAll(string(s), "/", cfg.PkgSep)
				impTmp = append(impTmp, api.Import(qualifier))
			}
			p.Imports = impTmp
		}

		node.Packages = tmp
	}

	switch cfg.OutputFormat {
	case Json:
		buf, err := json.Marshal(node)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal json: %w", err)
		}

		return buf, nil

	case Yaml:
		buf, err := yaml.Marshal(node)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal yaml: %w", err)
		}

		return buf, nil
	case Adoc, Pdf, Html:
		output, _ := asciidoc.CreateModuleTemplate(golang.NewAModule(*node))

		return output.Bytes(), nil
	default:
		return nil, fmt.Errorf("invalid output format: %s", cfg.OutputFormat)
	}
}
