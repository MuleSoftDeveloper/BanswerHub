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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Credentials struct used for AnswerHub
// These are the parameters of the credentials.json file
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

// Attachment is a attachment object that might exist in a node.
type Attachment struct {
	ID            int    `json:"id"`
	FileName      string `json:"fileName"`
	Size          int    `json:"size"`
	SizeFormatted string `json:"sizeFormatted"`
	URL           string `json:"url"`
	Image         bool   `json:"image"`
}

// Topics are the question topics
type Topics struct {
	ID                    int    `json:"id"`
	CreationDate          int    `json:"creationDate"`
	CreationDateFormatted string `json:"creationDateFormatted"`
	Name                  string `json:"name"`
	Author                Author `json:"author"`
	UsedCount             int    `json:"usedCount"`
}

// Node is the general item object for a node
type Node struct {
	ID                    int          `json:"id"`
	Type                  string       `json:"type"`
	CreationDate          int          `json:"creationDate"`
	CreationDateFormatted string       `json:"creationDateFormatted"`
	Title                 string       `json:"title"`
	Body                  string       `json:"body"`
	BodyAsHTML            string       `json:"bodyAsHTML"`
	Author                Author       `json:"author"`
	LastEditedAction      int          `json:"lastEditedAction"`
	ActiveRevisionID      int          `json:"activeRevisionId"`
	RevisionIDs           []int        `json:"revisionIDs"`
	LastActiveUserID      int          `json:"lastActiveUserId"`
	LastActiveDate        int          `json:"lastActiveDate"`
	Attachments           []Attachment `json:"attachments"`
	ChildrenIDs           []int        `json:"childrenIds"`
	CommentIDs            []int        `json:"commentIds"`
	Marked                bool         `json:"marked"`
	Topics                []Topics     `json:"topics"`
	PrimaryContainerID    int          `json:"primaryContainerId"`
	ContainerIDs          []int        `json:"containerIds"`
	Slug                  string       `json:"slug"`
	Wiki                  bool         `json:"wiki"`
	Score                 int          `json:"score"`
	Depth                 int          `json:"depth"`
	ViewCount             int          `json:"viewCount"`
	UpVoteCount           int          `json:"upVoteCount"`
	DownVoteCount         int          `json:"downVoteCount"`
	NodeStates            []string     `json:"nodeStates"`
	Answers               []int        `json:"answers"`
	AnswerCount           int          `json:"answerCount"`
}

// Questions user questions data type
// This is the main data type for questions retrieved
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

// Action is an individual action in a list of actions
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

// loadCredentials reads the credentials file
// returns a Credentials object
func loadCredentials(file string) Credentials {
	var config Credentials
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Printf("Please enter all credentials flags or include credentials.json file.\nUse -h flag to see list of available arguments and options\n\n")
		log.Fatal(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

// getAPIURL appends the API services to for making a AnswerHub API call
// returns a concatinated string
func getAPIURL(base string) string {
	return base + "/services/v2/"
}

// makeRequest is a flexible HTTP request function call that takes
// the method, resource path, body as a byte array, and Basic authentication credentials
// returns a io byte array of the response body
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

// customBody is a default body for updating a spam post with this predefined body
// returns a stringified JSON object
func customBody() string {
	return `{"body":"Nothing here"}`
}

// deleteQuestion deletes a question by taking the question ID
// calls makeRequest() with a PUT to delete
// a 204 is a successful response
func deleteQuestion(ID int, auth Credentials) {
	path := "node/" + strconv.Itoa(ID) + "/delete.json"
	println("\nDeleting question ID:", ID)
	makeRequest("PUT", path, nil, auth)
}

// updateQuestion updates a question by taking the question ID
// calls makeRequest() with a PUT to change
// a 204 is a successful response
func updateQuestion(ID int, auth Credentials) {
	path := "question/" + strconv.Itoa(ID) + ".json"
	body := customBody()
	println("\nUpdating question ID:", ID)
	makeRequest("PUT", path, []byte(body), auth)
}

// deactivateUser will deactivate a user's account
// takes the User ID parameter from the original CLI command
func deactivateUser(userID string, auth Credentials) {
	path := "user/" + userID + "/deactivateUser.json"
	println("\nDeactivating user:", userID)
	makeRequest("PUT", path, nil, auth)
}

// parseQuestionList takes the array of questions that a user has posted
// and iterates through it. It then deletes/updates each of these questions
// by calling updateQuestion() or deleteQuestion()
func parseQuestionList(list []Node, auth Credentials) {
	for _, v := range list {
		// updateQuestion(v.ID, auth)
		deleteQuestion(v.ID, auth)
	}
}

// getUserQuestionsByID retrieves the questions posted by the user
// returns the http body as io byte array
func getUserQuestionsByID(userID string, auth Credentials) []byte {
	path := "user/" + userID + "/question.json"
	println("\nFetching user's questions:", userID)
	return makeRequest("GET", path, nil, auth)
}

// processQuestionsBody takes the io byte array of the body and
// parses it as a Questions object
func processQuestionsBody(body []byte) *Questions {
	q := new(Questions)
	err := json.Unmarshal(body, &q)
	if err != nil {
		panic(err)
	}
	return q
}

// processUserQuestions is the logical processor of questions
// it takes the list of questions and sends to be processed/parsed
// for updating/deletion.
// If the number of questions in the current page is than the total,
// it will run this processor again to fetch the next list after processing
func processUserQuestions(userID string, auth Credentials) {
	qs := processQuestionsBody(getUserQuestionsByID(userID, auth))
	parseQuestionList(qs.List, auth)
	if qs.TotalCount > qs.ListCount {
		processUserQuestions(userID, auth)
	}
}

// getUserActionsByID fetches the actions a user performed
// returns a io byte array of the http response.
func getUserActionsByID(userID string, pageNumber int, auth Credentials) []byte {
	path := "user/" + userID + "/action.json?page=" + strconv.Itoa(pageNumber)
	println("\nFetching user's actions:", userID)
	return makeRequest("GET", path, nil, auth)
}

// deleteNode will make a HTTP request to delete a node by it's ID
func deleteNode(ID int, auth Credentials) {
	path := "node/" + strconv.Itoa(ID) + "/delete.json"
	println("\nDeleting node ID:", ID)
	makeRequest("PUT", path, nil, auth)
}

// parseActionList takes each of the action array items and sends for deletion
func parseActionList(list []Action, auth Credentials) {
	for _, v := range list {
		deleteNode(v.Node.ID, auth)
	}
}

// processActionsBody takes the response from HTTP to the /action.json endpoint
// return a Actions object
func processActionsBody(body []byte) *Actions {
	a := new(Actions)
	err := json.Unmarshal(body, &a)
	if err != nil {
		panic(err)
	}
	return a
}

// processUserActions is the logical processor for Actions
// It fetches the list of user actions,
// calls the parser to sift the list and perform deletion of actions.
// If the current page is not the last page, it will call the next page
// Actions that are deleted will still exist, they're just not published.
// The list size never reduces.
func processUserActions(userID string, lastPage int, auth Credentials) {
	lastPage++
	as := processActionsBody(getUserActionsByID(userID, lastPage, auth))
	parseActionList(as.List, auth)
	if as.PageCount > as.Page {
		processUserActions(userID, lastPage, auth)
	}
}

// main starts the app
// Gets the credentials from credentials.json located at the root
// Get the CLI options (user ID)
// begins the HTTP requests for getting user action
// finishes by deactivating user
// NOTE:	It currently does not get and process the user's questions!
// 			this is because questions count as actions.
func main() {
	unamePtr := flag.String("user", "", "Your Answerhub username")
	passPtr := flag.String("pass", "", "Your Answerhub password")
	forumURLPtr := flag.String("url", "", "The root URL for your Answerhub instance")
	banPtr := flag.String("ban", "", "ID of person to deactivate")
	flag.Parse()

	credentials := Credentials{}
	loadedCredentials := Credentials{}
	if *unamePtr == "" || *passPtr == "" || *forumURLPtr == "" {
		credPath := "credentials.json"
		loadedCredentials = loadCredentials(credPath)
	}

	if *unamePtr != "" {
		credentials.Username = *unamePtr
	} else {
		credentials.Username = loadedCredentials.Username
	}

	if *passPtr != "" {
		credentials.Password = *passPtr
	} else {
		credentials.Password = loadedCredentials.Password
	}

	if *forumURLPtr != "" {
		credentials.AnswerHubBaseURL = *forumURLPtr
	} else {
		credentials.AnswerHubBaseURL = loadedCredentials.AnswerHubBaseURL
	}

	if *banPtr != "" {
		// processUserQuestions(userID, credentials)
		processUserActions(*banPtr, 0, credentials)
		deactivateUser(*banPtr, credentials)
	} else {
		println("Provide user ID to ban using the -ban flag (e.g. -ban=1234)")
	}
}
