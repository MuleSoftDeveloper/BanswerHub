// UPDATE: This app now gets all actions by a user instead of Questions
// TODO: Break up into packages and add CLI options

// This app takes in a user ID from AnswerHub forums
// Using the ID it gets all questions written by user
// It updates all the body content of those questions
// Finally it deactivates the user

package main

import (
	"bytes"
	"encoding/json"
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
	Password         string `json:"password"`
}

// Author user / author object
type Author struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Realname   string `json:"realname"`
	Reputation int    `json:"reputation"`
}

// Topics question topics
type Topics struct {
	ID                    int    `json:"id"`
	CreationDate          int    `json:"creationDate"`
	CreationDateFormatted string `json:"creationDateFormatted"`
	Name                  string `json:"name"`
	Author                Author `json:"author"`
	UsedCount             int    `json:"usedCount"`
}

// Node general item object
type Node struct {
	ID                    int      `json:"id"`
	Type                  string   `json:"type"`
	CreationDate          int      `json:"creationDate"`
	CreationDateFormatted string   `json:"creationDateFormatted"`
	Title                 string   `json:"title"`
	Body                  string   `json:"body"`
	BodyAsHTML            string   `json:"bodyAsHTML"`
	Author                Author   `json:"author"`
	LastEditedAction      int      `json:"lastEditedAction"`
	ActiveRevisionID      int      `json:"activeRevisionId"`
	RevisionIDs           []int    `json:"revisionIDs"`
	LastActiveUserID      int      `json:"lastActiveUserId"`
	LastActiveDate        int      `json:"lastActiveDate"`
	Attachments           []string `json:"attachments"`
	ChildrenIDs           []int    `json:"childrenIds"`
	CommentIDs            []int    `json:"commentIds"`
	Marked                bool     `json:"marked"`
	Topics                []Topics `json:"topics"`
	PrimaryContainerID    int      `json:"primaryContainerId"`
	ContainerIDs          []int    `json:"containerIds"`
	Slug                  string   `json:"slug"`
	Wiki                  bool     `json:"wiki"`
	Score                 int      `json:"score"`
	Depth                 int      `json:"depth"`
	ViewCount             int      `json:"viewCount"`
	UpVoteCount           int      `json:"upVoteCount"`
	DownVoteCount         int      `json:"downVoteCount"`
	NodeStates            []string `json:"nodeStates"`
	Answers               []int    `json:"answers"`
	AnswerCount           int      `json:"answerCount"`
}

// Questions user questions data type
//
type Questions struct {
	Name       string `json:"name"`
	Sort       string `json:"sort"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	PageCount  int    `json:"pageCount"`
	ListCount  int    `json:"listCount"`
	TotalCount int    `json:"totalCount"`
	List       []Node `json:"list"`
}

//Action is an individual action in a list of actions
type Action struct {
	ID            int    `json:"id"`
	IP            string `json:"ip"`
	User          Author `json:"user"`
	ActionDate    int    `json:"actionDate"`
	Canceled      bool   `json:"canceled"`
	PrivateAction bool   `json:"privateAction"`
	Verb          string `json:"verb"`
	Node          Node   `json:"node"`
	RootNode      Node   `json:"rootNode"`
}

// Actions is the user's actions data type
type Actions struct {
	Name       string   `json:"name"`
	Sort       string   `json:"sort"`
	Page       int      `json:"page"`
	PageSize   int      `json:"pageSize"`
	PageCount  int      `json:"pageCount"`
	ListCount  int      `json:"listCount"`
	TotalCount int      `json:"totalCount"`
	List       []Action `json:"list"`
}

func loadCredentials(file string) Credentials {
	var config Credentials
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
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

func deleteQuestion(ID int, auth Credentials) {
	path := "node/" + strconv.Itoa(ID) + "/delete.json"
	println("\nDeleting question ID:", ID)
	makeRequest("PUT", path, nil, auth)

}

func updateQuestion(ID int, auth Credentials) {
	path := "question/" + strconv.Itoa(ID) + ".json"
	body := customBody()
	println("\nUpdating question ID:", ID)
	makeRequest("PUT", path, []byte(body), auth)
}

func deactivateUser(userID string, auth Credentials) {
	path := "user/" + userID + "/deactivateUser.json"
	println("\nDeactivating user:", userID)
	makeRequest("PUT", path, nil, auth)
}

func parseQuestionList(list []Node, auth Credentials) {
	for _, v := range list {
		// updateQuestion(v.ID, auth)
		deleteQuestion(v.ID, auth)
	}
}

func getUserQuestionsByID(userID string, auth Credentials) []byte {
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

func processUserQuestions(userID string, auth Credentials) {
	qs := processQuestionsBody(getUserQuestionsByID(userID, auth))
	parseQuestionList(qs.List, auth)
	if qs.TotalCount > qs.ListCount {
		processUserQuestions(userID, auth)
	}
}

func getUserActionsByID(userID string, pageNumber int, auth Credentials) []byte {
	path := "user/" + userID + "/action.json?page=" + strconv.Itoa(pageNumber)
	println("\nFetching user's actions:", userID)
	return makeRequest("GET", path, nil, auth)
}

func deleteNode(ID int, auth Credentials) {
	path := "node/" + strconv.Itoa(ID) + "/delete.json"
	println("\nDeleting node ID:", ID)
	makeRequest("PUT", path, nil, auth)

}

func parseActionList(list []Action, auth Credentials) {
	for _, v := range list {
		deleteNode(v.Node.ID, auth)
	}
}

func processActionsBody(body []byte) *Actions {
	a := new(Actions)
	err := json.Unmarshal(body, &a)
	if err != nil {
		panic(err)
	}
	return a
}

func processUserActions(userID string, lastPage int, auth Credentials) {
	lastPage++
	as := processActionsBody(getUserActionsByID(userID, lastPage, auth))
	parseActionList(as.List, auth)
	if as.PageCount > as.Page {
		processUserActions(userID, lastPage, auth)
	}
}

func main() {
	credPath := "credentials.json"
	credentials := loadCredentials(credPath)

	if len(os.Args) > 1 {
		userID := os.Args[1]
		// processUserQuestions(userID, credentials)
		processUserActions(userID, 0, credentials)
		deactivateUser(userID, credentials)
	} else {
		println("ProvIDe userID")
	}
}
