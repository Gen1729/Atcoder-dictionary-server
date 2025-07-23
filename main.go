package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type Link struct {
	Name string   `json:"name"`
	Url  string   `json:"url"`
	Tags []string `json:"tags"`
}

func getUser(c echo.Context) error {
	words := c.QueryParam("words")
	wordsArr := []string{}
	if words != "" {
		wordsArr = append(wordsArr, strings.Split(words, ",")...)
	}

	file, err := os.Open("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data []Link
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	var filtered []Link
	for _, link := range data {
		D := false
		for _, tag := range link.Tags {
			for _, word := range wordsArr {
				if tag == word {
					filtered = append(filtered, link)
					D = true
					break
				}
			}
			if D {
				break
			}
		}
	}
	u := &filtered

	return c.JSON(http.StatusCreated, u)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/dict", getUser)

	e.Logger.Fatal(e.Start(":1323"))
}
