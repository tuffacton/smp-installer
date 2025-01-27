package render

import (
	"context"
	"os"
	"path"
	"strings"
	"text/template"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

type TemplateRenderer struct {
	configStore store.DataStore
	outputStore store.DataStore
}

func NewTemplateRenderer(configStore store.DataStore, outputStore store.DataStore) *TemplateRenderer {
	return &TemplateRenderer{
		configStore: configStore,
		outputStore: outputStore,
	}
}

func (t *TemplateRenderer) Render(ctx context.Context,
	data map[string]interface{},
	templateDir string) error {

	allTpls := t.readTemplates(templateDir)
	vals := make(map[string]interface{})
	vals["Self"] = data
	vals["Config"] = t.configStore.DataMap(ctx)
	vals["Output"] = t.outputStore.DataMap(ctx)

	tpl := template.New("gotpl")
	tpl.Option("missingkey=zero")
	tpl.Funcs(sprig.FuncMap())

	for filename, filedata := range allTpls {
		_, err := tpl.New(filename).Parse(filedata)
		if err != nil {
			log.Err(err).Msgf("failed to parse template file %s/%s", templateDir, filename)
			return err
		}
	}

	rendered := make(map[string]string, len(allTpls))
	for filename, _ := range allTpls {
		var buf strings.Builder
		if err := tpl.ExecuteTemplate(&buf, filename, vals); err != nil {
			log.Err(err).Msgf("failed to execute template for file %s/%s", templateDir, filename)
			return err
		}
		rendered[filename] = strings.ReplaceAll(buf.String(), "<no value>", "")
	}

	return t.writeTemplates(templateDir, rendered)
}

func (t *TemplateRenderer) writeTemplates(templateDir string, rendered map[string]string) error {
	for filename, filedata := range rendered {
		err := os.WriteFile(path.Join(templateDir, filename), []byte(filedata), 0666)
		if err != nil {
			log.Err(err).Msgf("failed to write template file %s to directory %s", filename, templateDir)
			return err
		}
	}
	return nil
}

func (t *TemplateRenderer) readTemplates(templateDir string) map[string]string {
	files, err := os.ReadDir(templateDir)
	if err != nil {
		log.Err(err).Msgf("could not read template directory: %s", templateDir)
		return nil
	}
	templates := make(map[string]string)
	for _, f := range files {
		if f.Type().IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), ".terraform") {
			continue
		}
		if strings.HasPrefix(f.Name(), "terraform") {
			continue
		}
		data, err := os.ReadFile(path.Join(templateDir, f.Name()))
		if err != nil {
			log.Err(err).Msgf("cannot read file %s/%s", templateDir, f.Name())
			return nil
		}
		templates[f.Name()] = string(data)
	}
	return templates
}
