# AnswerHub Banishment Tool

This utility fetches a spammer's actions (questions and answers), deletes the content, and deactivates the user's account.

## Configuration
Create a *credentials.json* file with the following properties:
```JSON
{
    "answerHubBaseURL" : "Your AnswerHub URL (e.g. https://forums.mulesoft.com)",
    "username" : "Your AnswerHub Username",
    "password" : "Your AnswerHub Password"
}
```
## Building
Use GoLang's package builder to build the binary by running
`go build`

## Usage After Building
./banswerhub {{userid}}