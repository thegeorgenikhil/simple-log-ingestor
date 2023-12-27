package main

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
)

const (
	ErrorReportingTemplateFile = "error.html.tmpl"

	ErrorReportingSubject = "Error Reported!"
)

func ParseTemplate(fileName string, info any) (string, error) {
	t := template.New(fileName)

	templateFilePath := filepath.Join("alert_service", "templates", fileName)

	var err error
	t, err = t.ParseFiles(templateFilePath)
	if err != nil {
		log.Println("Error while parsing file:" + err.Error())
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, info); err != nil {
		log.Println("Error while executing template execution:" + err.Error())
		return "", err
	}

	result := tpl.String()

	return result, nil
}
