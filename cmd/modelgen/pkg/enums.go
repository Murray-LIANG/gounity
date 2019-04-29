package pkg

import (
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type enumField map[int][]string

type enums struct {
	enums          []*enum
	tmplFilePath   string
	fieldsFilePath string
	packageName    string
}

type enum struct {
	CapTypeName      string
	CapTypeShortName string
	Entries          []*entry
}

type entry struct {
	Key         string
	Value       int
	Description string
}

// NewEnums constructs a `enums`.
func NewEnums(tmplFilePath, fieldsFilePath, packageName string) *enums {
	return &enums{
		tmplFilePath:   tmplFilePath,
		fieldsFilePath: fieldsFilePath,
		packageName:    packageName,
	}
}

func (e *enums) parseFields() *enums {
	log := logrus.WithField("enums", e)

	rawFields, err := ioutil.ReadFile(e.fieldsFilePath)
	if err != nil {
		log.WithError(err).Fatal("read enums field file failed")
	}

	ef := make(map[string]enumField)
	err = yaml.Unmarshal(rawFields, ef)
	if err != nil {
		log.WithError(err).Fatal("unmarshal enums failed")
	}
	names := make([]string, 0)
	for name := range ef {
		names = append(names, name)
	}
	sort.Strings(names)
	tmpEnums := make([]*enum, 0)
	for _, name := range names {
		e := &enum{
			CapTypeName:      name,
			CapTypeShortName: strings.TrimSuffix(name, "Enum"),
			Entries:          []*entry{},
		}
		values := make([]int, 0)
		for value := range ef[name] {
			values = append(values, value)
		}
		sort.Ints(values)
		for _, value := range values {
			e.Entries = append(e.Entries, &entry{
				Key:         ef[name][value][0],
				Value:       value,
				Description: strings.Join(ef[name][value][1:], ", "),
			})
		}
		tmpEnums = append(tmpEnums, e)
		logrus.Infof("modelgen: enum %s parsed", name)
	}
	e.enums = tmpEnums
	return e
}

func (e *enums) GetTemplate() *template.Template {
	return template.Must(template.ParseFiles(e.tmplFilePath))
}

func (e *enums) PrepareData() interface{} {
	e = e.parseFields()
	return struct {
		Timestamp   time.Time
		PackageName string
		Enums       []*enum
	}{
		time.Now().UTC(),
		e.packageName,
		e.enums,
	}
}
