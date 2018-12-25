package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

type resource struct {
	packageName    string
	typeName       string
	tmplFilePath   string
	fieldsFilePath string
	fieldNames     []string
	fields         []field
	isEmbedded     bool
}

type field struct {
	CapName     string
	Type        string
	Description string
	JsonSrc     string
}

// NewResource constructs a `resource`.
func NewResource(
	tmplFilePath, fieldsFilePath, packageName, typeName string,
	isEmbedded bool,
) *resource {

	return &resource{
		tmplFilePath:   tmplFilePath,
		fieldsFilePath: fieldsFilePath,
		packageName:    packageName,
		typeName:       typeName,
		isEmbedded:     isEmbedded,
	}
}

func (r *resource) parseFields() *resource {
	log := logrus.WithField("resource", r)

	f, err := os.Open(r.fieldsFilePath)
	if err != nil {
		log.WithError(err).Fatal("open fields file failed")
	}
	defer f.Close()

	fieldNames := []string{}
	fields := []field{}

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		parts := strings.FieldsFunc(
			strings.TrimSpace(line),
			func(c rune) bool { return c == ',' },
		)
		if len(parts) < 2 {
			log.WithField("line", line).Fatal(
				"each line should at least have 2 items split by comma")
		}
		fieldNames = append(fieldNames, parts[0])
		description := ""
		if len(parts) >= 3 {
			description = strings.Join(parts[2:], ", ")
		}
		fields = append(fields, field{
			CapName:     strings.Title(parts[0]),
			Type:        parts[1],
			Description: description,
			JsonSrc:     fmt.Sprintf("`json:\"%s\"`", parts[0]),
		})

		if err != nil {
			break
		}
	}

	if err != nil && err != io.EOF {
		log.WithError(err).Fatal("read line from file failed")
	}
	r.fieldNames = fieldNames
	r.fields = fields

	logrus.Infof("modelgen: resource %s parsed", r.typeName)
	return r
}

func (r *resource) GetTemplate() *template.Template {
	return template.Must(template.ParseFiles(r.tmplFilePath))
}

func (r *resource) PrepareData() interface{} {
	r = r.parseFields()
	return struct {
		Timestamp   time.Time
		PackageName string
		TypeName    string
		CapTypeName string
		FieldNames  []string
		Fields      []field
		IsEmbedded  bool
		HasNameField bool
	}{
		time.Now().UTC(),
		r.packageName,
		r.typeName,
		strings.Title(r.typeName),
		r.fieldNames,
		r.fields,
		r.isEmbedded,
		contains(r.fieldNames, "name"),
	}
}

func contains(ss []string, e string) bool {
	for _, s := range ss {
		if s == e {
			return true
		}
	}
	return false
}
