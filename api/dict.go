package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Link struct {
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	Tags       []string `json:"tags"`
	Contest    int      `json:"contest"`
	Difficulty int      `json:"difficulty"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 開発段階として'*'、本番は特定URLに変更推奨
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	words := r.URL.Query().Get("words")
	var contest string = r.URL.Query().Get("contest")
	var difficulty string = r.URL.Query().Get("difficulty")

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

	minDifficulty := 0
	maxDifficulty := 3600
	if len(difficultyArr) >= 2 {
		minDifficulty = difficultyArr[0]
		maxDifficulty = difficultyArr[1]
	}

	file, err := os.Open("data.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to open data.json: " + err.Error()})
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data []Link
	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to decode JSON: " + err.Error()})
		return
	}

	var filteredByDifficultyContest []Link
	for _, link := range data {
		if !(minDifficulty <= link.Difficulty && link.Difficulty <= maxDifficulty) {
			continue
		}
		if contest == "" {
			filteredByDifficultyContest = append(filteredByDifficultyContest, link)
			continue
		}
		for _, num := range contestArr {
			if link.Contest == num {
				filteredByDifficultyContest = append(filteredByDifficultyContest, link)
				break
			}
		}
	}

	var filtered []Link
	for _, link := range filteredByDifficultyContest {
		matched := false
		for _, tag := range link.Tags {
			for _, word := range wordsArr {
				if tag == word {
					filtered = append(filtered, link)
					matched = true
					break
				}
			}
			if matched {
				break
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	if words == "" {
		json.NewEncoder(w).Encode(filteredByDifficultyContest)
		return
	}
	json.NewEncoder(w).Encode(filtered)
}
