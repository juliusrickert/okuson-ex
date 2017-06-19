package main

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"text/template"
)

type templateData struct {
	Sheets []sheet
}

func (templateData) Escape(s string) string {
	s = strings.Replace(s, "{", "\\{", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	//s = strings.Replace(s, "\\},\\{", "\\},\\\\\\{", -1)
	return s
}

func makeLaTeX(sheets []sheet, out io.Writer) error {
	t := template.New("tpl")
	t.Delims("<<", ">>")
	t, err := t.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	fileTemplates := t.Templates()

	if len(fileTemplates) == 0 {
		return errors.New("An unknown error occured while parsing the template")
	}

	err = t.ExecuteTemplate(out, fileTemplates[0].Name(), templateData{Sheets: sheets})
	if err != nil {
		return err
	}

	return nil
}

func makeJSON(sheets []sheet, out io.Writer) error {
	var exerciceMap = make(map[int]exercice)
	for _, s := range sheets {
		for _, e := range s.Exercices {
			exerciceMap[e.Number] = e
		}
	}

	enc := json.NewEncoder(out)
	err := enc.Encode(exerciceMap)
	if err != nil {
		return err
	}

	return nil
}
