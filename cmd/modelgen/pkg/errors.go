package pkg

import (
	"io/ioutil"
	"sort"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type errors struct {
	errors         []*errorEntry
	tmplFilePath   string
	fieldsFilePath string
	packageName    string
}

type errorEntry struct {
	Name      string
	ErrorCode int
}

// NewErrors constructs a `errors`.
func NewErrors(tmplFilePath, fieldsFilePath, packageName string) *errors {
	return &errors{
		tmplFilePath:   tmplFilePath,
		fieldsFilePath: fieldsFilePath,
		packageName:    packageName,
	}
}

func (e *errors) parseFields() *errors {
	log := logrus.WithField("errors", e)

	rawFields, err := ioutil.ReadFile(e.fieldsFilePath)
	if err != nil {
		log.WithError(err).Fatal("read errors field file failed")
	}

	ef := make(map[int]string)
	err = yaml.Unmarshal(rawFields, ef)
	if err != nil {
		log.WithError(err).Fatal("unmarshal errors failed")
	}
	codes := make([]int, 0)
	for c := range ef {
		codes = append(codes, c)
	}
	sort.Ints(codes)
	tmpErrors := make([]*errorEntry, 0)
	for _, c := range codes {
		e := &errorEntry{
			Name:      ef[c],
			ErrorCode: c,
		}
		tmpErrors = append(tmpErrors, e)
		logrus.Infof("modelgen: error %s parsed", ef[c])
	}
	e.errors = tmpErrors
	return e
}

func (e *errors) GetTemplate() *template.Template {
	return template.Must(template.ParseFiles(e.tmplFilePath))
}

func (e *errors) PrepareData() interface{} {
	e = e.parseFields()
	return struct {
		Timestamp   time.Time
		PackageName string
		Errors      []*errorEntry
	}{
		time.Now().UTC(),
		e.packageName,
		e.errors,
	}
}
