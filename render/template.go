package render

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// HTML built-in renderer.
type TemplateRender struct {
	ContentType string
	Charset     string
	Factory     TemplateFactory
	// If IsDevelopment is set to true, this will recompile the templates on every request. Default if false.
	IsDevelopment bool

	templates *template.Template
}

// NewXMLRender creates a new JSON render with default values.
func NewTemplateRender() *TemplateRender {
	return &TemplateRender{
		ContentType:   "text/html",
		Charset:       "UTF-8",
		IsDevelopment: false,
	}
}

func (render *TemplateRender) SetContentType(value string) *TemplateRender {
	render.ContentType = value
	return render
}

func (render *TemplateRender) SetCharset(value string) *TemplateRender {
	render.Charset = value
	return render
}

func (render *TemplateRender) SetFactory(value TemplateFactory) *TemplateRender {
	render.Factory = value
	return render
}

func (render *TemplateRender) CompileTemplates() error {
	if render.IsDevelopment || render.templates == nil {
		if tmpl, err := render.Factory.Create(); err != nil {
			return err
		} else {
			render.templates = tmpl
		}
	}
	return nil
}

// Render a template response.
func (render *TemplateRender) Render(writer http.ResponseWriter, status int, name string, binding interface{}) error {
	if err := render.CompileTemplates(); err != nil {
		return err
	}

	out := new(bytes.Buffer)
	if err := render.templates.ExecuteTemplate(out, name, binding); err != nil {
		return err
	}

	writeHeader(writer, status, render.ContentType, render.Charset)
	writer.Write(out.Bytes())
	return nil
}

type TemplateOptions struct {
	// Left delimiter, defaults to {{.
	LeftDelimiter string
	// Right delimiter, defaults to }}.
	RightDelimiter string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Functions []template.FuncMap
}

type TemplateFactory interface {
	Create() (*template.Template, error)
}

type DiskTemplateFactory struct {
	// Options of templates
	Options TemplateOptions
	// Directory to load templates. Default is "templates".
	Directory string
	// Extensions to parse template files from. Defaults to [".tmpl"].
	Extensions []string
}

func NewDiskTemplateFactory() *DiskTemplateFactory {
	return &DiskTemplateFactory{
		Directory:  "templates",
		Extensions: []string{".tmpl"},
	}
}

func (factory *DiskTemplateFactory) SetDirectory(value string) *DiskTemplateFactory {
	factory.Directory = value
	return factory
}

func (factory *DiskTemplateFactory) SetExtensions(values ...string) *DiskTemplateFactory {
	factory.Extensions = values
	return factory
}

func (factory *DiskTemplateFactory) SetDelimiters(left string, right string) *DiskTemplateFactory {
	factory.Options.LeftDelimiter = left
	factory.Options.RightDelimiter = right
	return factory
}

func (factory *DiskTemplateFactory) SetFunctions(functions ...template.FuncMap) *DiskTemplateFactory {
	factory.Options.Functions = functions
	return factory
}

func (factory *DiskTemplateFactory) Create() (*template.Template, error) {
	result := template.New(factory.Directory)
	result.Delims(factory.Options.LeftDelimiter, factory.Options.RightDelimiter)

	return result, filepath.Walk(factory.Directory, func(path string, info os.FileInfo, err error) error {
		relative, err := filepath.Rel(factory.Directory, path)
		if err != nil {
			return err
		}

		ext := filepath.Ext(relative)

		for _, extension := range factory.Extensions {
			if ext == extension {
				buf, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				tmpl := result.New(filepath.ToSlash((relative[0 : len(relative)-len(ext)])))
				for _, funcs := range factory.Options.Functions {
					tmpl.Funcs(funcs)
				}

				if _, err := tmpl.Parse(string(buf)); err != nil {
					return err
				}
				break
			}
		}
		return nil
	})
}

type SubTemplate func() (string, string)

type StaticTemplateFactory struct {
	// Options of templates
	Options TemplateOptions
	// Name of the template
	Name string
	// Subtemplates facotry
	SubTemplates []SubTemplate
}

func NewStaticTemplateFactory(name string) *StaticTemplateFactory {
	return &StaticTemplateFactory{
		Name: name,
	}
}

func (factory *StaticTemplateFactory) SetSubTemplates(subs ...SubTemplate) *StaticTemplateFactory {
	factory.SubTemplates = subs
	return factory
}

func (factory *StaticTemplateFactory) SetDelimiters(left string, right string) *StaticTemplateFactory {
	factory.Options.LeftDelimiter = left
	factory.Options.RightDelimiter = right
	return factory
}

func (factory *StaticTemplateFactory) SetFunctions(functions ...template.FuncMap) *StaticTemplateFactory {
	factory.Options.Functions = functions
	return factory
}

func (factory *StaticTemplateFactory) Create() (*template.Template, error) {
	result := template.New(factory.Name)
	result.Delims(factory.Options.LeftDelimiter, factory.Options.RightDelimiter)

	for _, sub := range factory.SubTemplates {
		name, content := sub()
		tmpl := result.New(name)
		for _, funcs := range factory.Options.Functions {
			tmpl.Funcs(funcs)
		}

		if _, err := tmpl.Parse(content); err != nil {
			return nil, err
		}
	}
	return result, nil
}
