package main

import (
	"errors"
	"fmt"
	"strings"
)

type Takon struct {
	Content   []string
	State     TakonState
	Locations []Location
}

type TakonState struct {
	Stack []string
	State string
}

type Location struct {
	Keyword     string
	EqualSymbol string
	Path        string
	EndSymbol   string
	Content     []string
	TryFiles    string
	Lastline    string
}

type Output struct {
	node []Node
}

func NewTakon() Takon {
	return Takon{}
}

func (T *Takon) SetContent(content []string) {
	for _, v := range content {
		if !strings.HasPrefix(strings.TrimSpace(v), "#") {
			T.Content = append(T.Content, v)
		}
	}
}

func (T *Takon) Start() error {
	for _, c := range T.Content {
		c := strings.TrimSpace(c)
		if strings.HasPrefix(c, "location") {
			if len(T.State.Stack) != 0 {
				return errors.New("too much {")
			}
			location := Location{}
			takons := strings.Split(c, " ")
			if len(takons) < 3 {
				return errors.New("find Keyword LOCATION,but Unrecognized grammar")
			}
			switch len(takons) {
			case 3:
				location.Keyword = strings.Replace(takons[0], " ", "", -1)
				location.EqualSymbol = ""
				location.Path = strings.Replace(takons[1], " ", "", -1)
				location.EndSymbol = strings.Replace(takons[2], " ", "", -1)
			case 4:
				location.Keyword = strings.Replace(takons[0], " ", "", -1)
				location.EqualSymbol = strings.Replace(takons[1], " ", "", -1)
				location.Path = strings.Replace(takons[2], " ", "", -1)
				location.EndSymbol = strings.Replace(takons[3], " ", "", -1)
			}

			if !location.CheckLocation() {
				return errors.New("find Keyword LOCATION,but error grammar")
			}
			T.State.Stack = append(T.State.Stack, location.EndSymbol)
			T.Locations = append(T.Locations, location)
			continue
		}
		if len(T.State.Stack) != 0 {
			location := &(T.Locations[len(T.Locations)-1])
			if strings.Contains(c, "try_files") {

				takons := strings.Split(c, " ")
				if len(takons) != 3 {
					return errors.New("try_files: error grammar")
				}
				(*location).TryFiles = strings.Trim(takons[2], "@ ;")
				continue
			}
			if strings.Contains(c, "}") {
				T.State.Stack = T.State.Stack[:len(T.State.Stack)-1]
				(*location).Lastline = c
				continue
			}
			(*location).Content = append((*location).Content, c)
		}
	}
	return nil
}

func (L *Location) CheckLocation() bool {
	if L.EqualSymbol == "~*" || L.EqualSymbol == "~" || L.EqualSymbol == "=" || L.EqualSymbol == "^~" || L.EqualSymbol == "" {
		if L.EndSymbol == "{" {
			return true
		} else {
			fmt.Println("find {} err:", L.EndSymbol)
		}
	} else {
		fmt.Println("find err:", L)
	}
	return false
}
