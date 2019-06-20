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
	metadata       []string
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

	metadata := []string{}
	fieldNames := []string{}
	fields := []field{}

	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		if line == "" && err != nil {
			log.Debug("reach end of file")
			break
		}

		parts := strings.FieldsFunc(
			strings.TrimSpace(line),
			func(c rune) bool { return c == ',' },
		)
		if len(parts) == 0 {
			log.WithField("line", line).Info("empty line")
			continue
		}

		if strings.HasPrefix(parts[0], "__") && strings.HasSuffix(parts[0], "__") {
			metadata = append(metadata, parts[0])
			continue
		}

		fieldNames = append(fieldNames, parts[0])
		if len(parts) == 1 {
			log.WithField("line", line).Info(
				"line has one item, take it as nested property")
			continue
		}
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
	}

	if err != nil && err != io.EOF {
		log.WithError(err).Fatal("read line from file failed")
	}
	r.fieldNames = fieldNames
	r.fields = fields
	r.metadata = metadata

	logrus.Infof("modelgen: resource %s parsed", r.typeName)
	return r
}

func (r *resource) GetTemplate() *template.Template {
	return template.Must(template.ParseFiles(r.tmplFilePath))
}

func (r *resource) PrepareData() interface{} {
	r = r.parseFields()
	return struct {
		Timestamp               time.Time
		PackageName             string
		TypeName                string
		CapTypeName             string
		FieldNames              []string
		Fields                  []field
		IsEmbedded              bool
		HasNameField            bool
		HasStorageResourceField bool
		DeleteDirectly          bool
	}{
		time.Now().UTC(),
		r.packageName,
		r.typeName,
		strings.Title(r.typeName),
		r.fieldNames,
		r.fields,
		r.isEmbedded,
		contains(r.fieldNames, "name"),
		contains(r.fieldNames, "storageResource"),
		contains(r.metadata, "__DeleteDirectly__"),
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
