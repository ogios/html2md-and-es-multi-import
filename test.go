package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/mattn/godown"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Blog struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

const (
	INDEX_NAME = "pblog"
	MAPPING    = `{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "term_vector": "with_positions_offsets"
      },
      "content": {
        "type": "text",
        "term_vector": "with_positions_offsets"
      }
    }
  }
}`
)

var CLI *elasticsearch.Client

func InitCLI() {
	fmt.Println("Connecting")
	cli, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{GLOBAL_CONFIG.ESAddr},
	})
	if err != nil {
		panic(err)
	}
	CLI = cli
}

func CreateIndex() {
	fmt.Println("Check if exist")
	res, err := CLI.Indices.Exists([]string{INDEX_NAME})
	fmt.Println(*res)
	if err != nil {
		panic(err)
	}
	if res.StatusCode == 200 {
		return
	} else {
		res, err = CLI.Indices.Create(INDEX_NAME, CLI.Indices.Create.WithBody(strings.NewReader(MAPPING)))
		if err != nil {
			panic(err)
		}
		fmt.Println(*res)
	}
}

func DeleteIndex() {
	fmt.Println("Check if exist")
	res, err := CLI.Indices.Exists([]string{INDEX_NAME})
	fmt.Println(*res)
	if err != nil {
		panic(err)
	}
	if res.StatusCode == 200 {
		res, err = CLI.Indices.Delete([]string{INDEX_NAME})
		if err != nil {
			panic(err)
		}
		fmt.Println(*res)
	}
}

func NewTable(id int, b *Blog) {
	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	res, err := CLI.Index(INDEX_NAME, bytes.NewReader(data), CLI.Index.WithDocumentID(strconv.Itoa(id)))
	if err != nil {
		panic(err)
	}
	if res.StatusCode == 200 || res.StatusCode == 201 {
		fmt.Println("[ ok  ]:", id)
	} else {
		fmt.Println("[fatal]:", id, "-", res.StatusCode)
	}
}

func NewTableTest() {
	data, err := json.Marshal(Blog{
		Title:   "TEST2 TITLE",
		Content: "TEST2 CONTENT",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	res, err := CLI.Index(INDEX_NAME, bytes.NewReader(data), CLI.Index.WithDocumentID("1"))
	if err != nil {
		panic(err)
	}
	fmt.Println(*res)
}

var DB *gorm.DB

func DBInit() {
	d := GLOBAL_CONFIG.MysqlDSM
	db, err := gorm.Open(mysql.Open(d), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	DB = db
}

func addOneBlog(b *ITC) {
	res, err := http.Get("http://" + GLOBAL_CONFIG.Base + "/raw/text/" + b.Content)
	fmt.Println(res.Request.URL)
	if err != nil {
		panic(err)
	}
	defer res.Request.Body.Close()
	if res.StatusCode != 200 {
		panic(fmt.Sprintln("wrong status: ", res.StatusCode))
	}
	c := &CustomWriter{}
	err = godown.Convert(c, res.Body, nil)
	if err != nil {
		panic(err)
	}
	blog := Blog{
		Title:   b.Title,
		Content: string(c.Content),
	}
	NewTable(b.Id, &blog)
}

type ITC struct {
	Title   string
	Content string
	Id      int
}

func Test() {
	s := []*ITC{}
	DB.Table("t_blog").Select("id", "title", "content").Find(&s)
	fmt.Println(s)
	for _, b := range s {
		addOneBlog(b)
	}
}
