package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// firstJSON()

	secondJSON()
}

var bookTitle = "Building Modern CLI Applications in Go"
var bookAuthor = "Packt Publishing Ltd"

// firstJSON map 编码解密
func firstJSON() {
	var book = make(map[string]interface{})
	book["id"] = 100
	book["title"] = bookTitle
	book["author"] = bookAuthor

	bs, _ := json.Marshal(book)
	fmt.Println(string(bs))

	var book2 = make(map[string]interface{})
	_ = json.Unmarshal(bs, &book2)
	fmt.Println(book2)
}

// secondJSON struct 序列名和反序列化（编码/解码）
func secondJSON() {
	type Book struct {
		ID     int
		Title  string
		Author string
	}

	var book = Book{
		ID:     100,
		Title:  bookTitle,
		Author: bookAuthor,
	}

	// 编码
	bs, _ := json.Marshal(book)
	fmt.Println(string(bs))

	// 格式化编码
	bsFormat, _ := json.MarshalIndent(book, "", " ")
	fmt.Println(string(bsFormat))

	var b Book
	_ = json.Unmarshal(bs, &b)
	fmt.Println(b.ID, b.Title, b.Author)

	var input = `
	{
		"ID": 100,
		"Title": "Building Modern CLI Applications in Go",
		"Author": "Packt Publishing Ltd",
		"Price": 55
	}
	`
	var book2 Book
	// Price会被忽略, 结构体中不存在的字段会被忽略
	_ = json.Unmarshal([]byte(input), &book2)
	fmt.Printf("%v\n", book2)

}
