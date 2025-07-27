package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Link struct {
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Tags       []string `json:"tags"`
	Contest    int      `json:"contest"`
	Difficulty int      `json:"difficulty"`
}

func getProblems(c echo.Context) error {
	words := c.QueryParam("words")
	var contest string = c.QueryParam("contest")
	var difficulty string = c.QueryParam("difficulty")
	wordsArr := []string{}
	if words != "" {
		wordsArr = append(wordsArr, strings.Split(words, ",")...)
	}

	contestArr := []int{}
	if contest != "" {
		contestStrArr := strings.Split(contest, ",")
		for _, c := range contestStrArr {
			if c == "" {
				continue
			}
			num, err := strconv.Atoi(c)
			if err != nil {
				continue
			}
			contestArr = append(contestArr, num)
		}
	}

	difficultyArr := []int{}
	difficultyStrArr := strings.Split(difficulty, ",")
	for _, c := range difficultyStrArr {
		if c == "" {
			continue
		}
		num, err := strconv.Atoi(c)
		if err != nil {
			continue
		}
		difficultyArr = append(difficultyArr, num)
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

	var filteredByDifficultyContest []Link
	for _, link := range data {
		if !(difficultyArr[0] <= link.Difficulty && link.Difficulty <= difficultyArr[1]) {
			continue
		}
		if contest == "" {
			filteredByDifficultyContest = append(filteredByDifficultyContest, link)
			continue
		}
		for _, num := range contestArr {
			if link.Contest == num {
				filteredByDifficultyContest = append(filteredByDifficultyContest, link)
			}
		}
	}

	var filtered []Link
	for _, link := range filteredByDifficultyContest {
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

	if words == "" {
		u := &filteredByDifficultyContest
		return c.JSON(http.StatusCreated, u)
	}
	u := &filtered
	return c.JSON(http.StatusCreated, u)
}

func main() {
	e := echo.New()

	// CORSミドルウェアを使う（許可するオリジンを指定）
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"https://atcoder-dictionary.vercel.app", // フロントエンドのURLを指定
			"http://localhost:5173",                 // 開発用にlocalhostも許可（任意）
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/dict", getProblems)

	e.Logger.Fatal(e.Start(":1323"))
}
