// This app takes in a user ID from AnswerHub forums
// Using the ID it gets all questions written by user
// It updates all the body content of those questions
// Finally it deactivates the user

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Credentials used for AnswerHub
type Credentials struct {
	AnswerHubBaseURL string `json:"answerHubBaseURL"`
	Username         string `json:"username"`
	Password         string `json:"password`
}

// Author question author
type Author struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Realname   string `json:"realname"`
	Reputation int    `json:"reputation"`
}

// Topics question topics
type Topics struct {
	Id                    int    `json:"id"`
	CreationDate          int    `json:"creationDate"`
	CreationDateFormatted string `json:"creationDateFormatted"`
	Name                  string `json:"name"`
	Author                Author `json:"author"`
	UsedCount             int    `json:"usedCount"`
}

// Question object
type Question struct {
	Id                 int      `json:"id"`
	Type               string   `json:"type"`
	CreationDate       int      `json:"creationDate"`
	Title              string   `json:"creationDateFormatted"`
	Body               string   `json:"body"`
	BodyAsHTML         string   `json:"bodyAsHTML"`
	Author             Author   `json:"Author"`
	LastEditedAction   int      `json:"lastEditedAction"`
	ActiveRevisionId   int      `json:"activeRevisionId`
	RevisionIds        []int    `json:"revisionIds"`
	LastActiveUserId   int      `json:"lastActiveUserId"`
	LastActiveDate     int      `json:"lastActiveDate"`
	Attachments        []string `json:"attachments"`
	ChildrenIds        []int    `json:"childrenIds"`
	CommentIds         []int    `json:"commentIds"`
	Marked             bool     `json:"marked"`
	Topics             []Topics `json:"topics"`
	PrimaryContainerId int      `json:"primaryContainerId"`
	ContainerIds       []int    `json:"containerIds"`
	Slug               string   `json:"slug"`
	Wiki               bool     `json:"wiki"`
	Score              int      `json:"score"`
	Depth              int      `json:"depth"`
	ViewCount          int      `json:"viewCount"`
	UpVoteCount        int      `json:"upVoteCount"`
	DownVoteCount      int      `json:"downVoteCount"`
	NodeStates         []string `json:"nodeStates"`
	Answers            []int    `json:"answers"`
	AnswerCount        int      `json:"answerCount"`
}

// Questions user questions data type
//
type Questions struct {
	Name       string     `json:"name"`
	Sort       string     `json:"sort"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
	PageCount  int        `json:"pageCount"`
	ListCount  int        `json:"listCount"`
	TotalCount int        `json:"totalCount"`
	List       []Question `json:"list"`
}

func loadCredentials(file string) Credentials {
	var config Credentials
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func getAPIURL(base string) string {
	return base + "/services/v2/"
}

func makeRequest(method string, path string, body []byte, auth Credentials) []byte {
	url := getAPIURL(auth.AnswerHubBaseURL) + path
	println("Making request to:", url)
	b := bytes.NewReader(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(auth.Username, auth.Password)
	req.Header.Add("content-type", "application/json")
	newClient := &http.Client{Timeout: time.Second * 10}
	resp, err := newClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	println("Request status", resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return respBody
}

// func base64Credentials(auth string) string {
// 	return base64.StdEncoding.EncodeToString([]byte(auth))
// }

func customBody() string {
	return `{"body":"Nothing here"}`
}

func deleteQuestion(id int, auth Credentials) {
	path := "node/" + strconv.Itoa(id) + "/delete.json"
	println("\nDeleting question Id:", id)
	makeRequest("PUT", path, nil, auth)

}

func updateQuestion(id int, auth Credentials) {
	path := "question/" + strconv.Itoa(id) + ".json"
	body := customBody()
	println("\nUpdating question Id:", id)
	makeRequest("PUT", path, []byte(body), auth)
}

func deactivateUser(userID string, auth Credentials) {
	path := "user/" + userID + "/deactivateUser.json"
	println("\nDeactivating user:", userID)
	makeRequest("PUT", path, nil, auth)
}

func parseQuestionList(list []Question, auth Credentials) {
	for _, v := range list {
		updateQuestion(v.Id, auth)
		deleteQuestion(v.Id, auth)
	}
}

func getUserByID(userID string, auth Credentials) []byte {
	path := "user/" + userID + "/question.json"
	println("\nFetching user's questions:", userID)
	return makeRequest("GET", path, nil, auth)
}

func processQuestionsBody(body []byte) *Questions {
	q := new(Questions)
	err := json.Unmarshal(body, &q)
	if err != nil {
		panic(err)
	}
	return q
}

func startBanishment(userID string, auth Credentials) {
	qs := processQuestionsBody(getUserByID(userID, auth))
	parseQuestionList(qs.List, auth)
	if qs.TotalCount > qs.ListCount {
		startBanishment(userID, auth)
	}
	deactivateUser(userID, auth)
}

func main() {
	credPath := "credentials.json"
	credentials := loadCredentials(credPath)

	if len(os.Args) > 1 {
		userID := os.Args[1]
		startBanishment(userID, credentials)
	} else {
		println("Provide userID")
	}
}
