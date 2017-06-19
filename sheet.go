package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const sheetListUri = "exquery.html"
const sheetUri = "QuerySheet"

func getSheetList() ([]string, error) {
	doc, err := goquery.NewDocument(okusonURL + sheetListUri)
	if err != nil {
		return []string{}, err
	}

	buttons := doc.Find("input[type='submit'][name='sheet']")
	sheets := make([]string, buttons.Length())

	buttons.Each(func(i int, btn *goquery.Selection) {
		sheets[i], _ = btn.Attr("value")
	})

	return sheets, nil
}

type exercice struct {
	Number int
	LaTeX  string
	Tasks  []task
}

const (
	multipleChoiceTask = "multipleChoice"
	choiceTask         = "choice"
	textTask           = "text"
)

type task struct {
	Type           string
	LaTeX          string
	Answers        []string `json:",omitempty"`
	CorrectAnswers string   `json:",omitempty"`
}

type sheet struct {
	Exercices []exercice
}

type outputFunction func([]sheet, io.Writer) error

func getSheet(s, username, password string) (*sheet, error) {
	form := url.Values{}
	form.Add("id", username)
	form.Add("passwd", password)
	form.Add("format", "HTML")
	form.Add("sheet", s)

	resp, err := http.Post(okusonURL+sheetUri, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	decoder := charmap.ISO8859_1.NewDecoder().Reader(resp.Body)

	doc, err := goquery.NewDocumentFromReader(decoder)
	if err != nil {
		return nil, err
	}

	h1 := doc.Find("h1").Text()
	if h1 == "Fehler: Falsches Passwort" || h1 == "Fehler: Ungültige Matrikelnummer" {
		return nil, errors.New("Wrong login credentials")
	}

	currentSheet := sheet{Exercices: []exercice{}}

	var currentExercice *exercice

	doc.Find("tr").Each(func(_ int, tr *goquery.Selection) {
		numChildren := tr.Children().Length()
		if numChildren == 2 {
			if currentExercice != nil {
				currentSheet.Exercices = append(currentSheet.Exercices, *currentExercice)
			}
			currentExercice = &exercice{}
			currentExercice.Number, _ = strconv.Atoi(tr.Find("td.exnr").Text())
			currentExercice.LaTeX, _ = tr.Find("td > img").Attr("alt")
			currentExercice.Tasks = []task{}
		} else if numChildren == 3 {
			if currentExercice == nil {
				return
			}
			currentTask := task{}
			currentTask.LaTeX, _ = tr.Find("td.question img").Attr("alt")
			currentTask.CorrectAnswers = tr.Find("span.erg").Text()
			inputs := tr.Find("input")
			inputType, _ := inputs.Attr("type")
			if inputType == "checkbox" || inputType == "radio" {
				if inputType == "checkbox" {
					currentTask.Type = multipleChoiceTask
				} else {
					currentTask.Type = choiceTask
				}
				currentTask.Answers = make([]string, inputs.Length())
				inputs.Each(func(i int, input *goquery.Selection) {
					currentTask.Answers[i] = strings.TrimSpace(input.Parent().Text())
				})
			} else {
				currentTask.Type = textTask
			}
			currentExercice.Tasks = append(currentExercice.Tasks, currentTask)
		}
	})
	if currentExercice != nil && currentExercice.LaTeX != "" {
		currentSheet.Exercices = append(currentSheet.Exercices, *currentExercice)
	}

	return &currentSheet, nil
}
