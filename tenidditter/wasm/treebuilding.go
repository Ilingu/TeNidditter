package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Webassembly Connected (treebuilding.go)!")

	c := make(chan struct{}, 0)
	js.Global().Set("TreeBuilding", TreeBuilding())
	<-c
}

type TedditCommmentShape struct {
	Id       int `json:"id"`
	ParentId int `json:"parentId"`

	Created     int64  `json:"created"`
	Ups         int    `json:"ups"`
	Body_html   string `json:"body_html"`
	Link_author string `json:"link_author"`
}

func TreeBuilding() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of argument passed"
		}

		stringComments := args[0].String()

		var comments [][]TedditCommmentShape
		err := json.Unmarshal(stringComments, &comments)
		if err != nil {
			return nil
		}

		return BuildTree(comments)
	})
}

func BuildTree(comments [][]TedditCommmentShape) string {
	return ""
}
