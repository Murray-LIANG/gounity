package main

import (
	"bytes"
	"flag"
	"go/build"
	"io/ioutil"
	"path"
	"strings"
	"text/template"

	"github.com/factioninc/gounity/cmd/modelgen/pkg"
	"github.com/sirupsen/logrus"
)

type generator interface {
	GetTemplate() *template.Template
	PrepareData() interface{}
}

func generate(g generator, outputPath string) {
	t := g.GetTemplate()
	data := g.PrepareData()
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		logrus.WithError(err).WithFields(map[string]interface{}{
			"template": t,
			"data":     data,
		}).Fatal("fill data to template failed")
	}
	if err := ioutil.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		logrus.WithField("outputPath", outputPath).WithError(err).Fatal(
			"write to file for enums failed")
	}
}

func main() {
	var (
		workingDir = flag.String("w", ".", "the working directory.")
		outputDir  = flag.String("o", ".", "the directory to put the generated go files.")
	)

	flag.Parse()

	log := logrus.WithField("workingDir", *workingDir).WithField("outputDir", *outputDir)

	pkgInfo, err := build.Default.ImportDir(".", 0)
	if err != nil {
		log.WithError(err).Fatal("get package info failed")
	}

	fieldsRootDir := path.Join(*workingDir, "fields")
	templatesRootDir := path.Join(*workingDir, "templates")

	var g generator
	// Generate enums
	g = pkg.NewEnums(
		path.Join(templatesRootDir, "enums.tmpl"),
		path.Join(fieldsRootDir, "enums.yml"),
		pkgInfo.Name,
	)
	outputPath := path.Join(*outputDir, "enums_gen.go")
	generate(g, outputPath)
	logrus.Info("modelgen: enums generated")

	// Generate errors
	g = pkg.NewErrors(
		path.Join(templatesRootDir, "errors.tmpl"),
		path.Join(fieldsRootDir, "errors.yml"),
		pkgInfo.Name,
	)
	outputPath = path.Join(*outputDir, "errors_gen.go")
	generate(g, outputPath)
	logrus.Info("modelgen: errors generated")

	// Generate resources one by one
	resourceFieldsDir := path.Join(fieldsRootDir, "resource")

	files, err := ioutil.ReadDir(resourceFieldsDir)
	if err != nil {
		log.WithError(err).Fatal("list field files of resources failed")
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".csv") {
			continue
		}
		resourceName := strings.TrimSuffix(f.Name(), ".csv")
		isEmbedded := false
		if strings.HasSuffix(resourceName, "_embed") {
			isEmbedded = true
			resourceName = strings.TrimSuffix(resourceName, "_embed")
		}

		g = pkg.NewResource(
			path.Join(templatesRootDir, "resource.tmpl"),
			path.Join(resourceFieldsDir, f.Name()),
			pkgInfo.Name, resourceName, isEmbedded,
		)

		outputPath := path.Join(*outputDir, strings.ToLower(resourceName+"_gen.go"))
		generate(g, outputPath)
		logrus.Infof("modelgen: resource %s generated", resourceName)
	}

}
