package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"os"
	"sort"
)

func combine() error {
	inputMaps := make([]map[int]exercice, len(inputFiles))

	for i, f := range inputFiles {
		m := map[int]exercice{}
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(file)
		err = dec.Decode(&m)
		if err != nil {
			return err
		}
		inputMaps[i] = m
	}

	targetMap := make(map[int]*exercice)

	hashMap := make(map[int]map[string]bool)

	count := make([]int, len(inputMaps))

	for mi, m := range inputMaps {
		for ei, e := range m {
			var target *exercice
			if targetMap[ei] == nil {
				target = &exercice{
					LaTeX:  e.LaTeX,
					Number: e.Number,
					Tasks:  []task{},
				}
			} else {
				target = targetMap[ei]
			}
			if target.LaTeX != e.LaTeX {
				return errors.New("Exercise mismatch")
			}
			for _, t := range e.Tasks {
				h := md5.Sum([]byte(t.LaTeX))
				hStr := string(h[:])
				if hashMap[ei] == nil {
					hashMap[ei] = make(map[string]bool)
				}
				if !hashMap[ei][hStr] {
					count[mi]++
					target.Tasks = append(target.Tasks, t)
					hashMap[ei][hStr] = true
				}
			}
			targetMap[ei] = target
		}
	}

	targetArr := []exercice{}
	for _, e := range targetMap {
		if e == nil {
			continue
		}
		targetArr = append(targetArr, *e)
	}

	sort.Sort(exerciceArray(targetArr))

	sheets := []sheet{
		{Exercices: targetArr},
	}

	var makeFunc outputFunction
	if outputFormat == "tex" {
		makeFunc = outputFunction(makeLaTeX)
	} else if outputFormat == "json" {
		makeFunc = outputFunction(makeJSON)
	}

	err := makeFunc(sheets, os.Stdout)

	if err != nil {
		return err
	}

	return nil
}

type exerciceArray []exercice

func (e exerciceArray) Len() int {
	return len(e)
}
func (e exerciceArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
func (e exerciceArray) Less(i, j int) bool {
	return e[i].Number < e[j].Number
}
