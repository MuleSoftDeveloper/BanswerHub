# AnswerHub Banishment Tool

This utility fetches a spammer's actions (questions and answers), deletes the content, and deactivates the user's account.

## Optional Configuration
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
`go install`

## Usage After Building
If using without creating the optional credentials.json file:

`banswerhub -user={{YOUR_USERNAME}} -pass={{YOUR_PASSWORD}} -url={{ANSWERHUB_ROOT_URL}} -ban={{USER ID OF USER TO BAN}}`

If using with the credentials.json:

`banswerhub -ban={{USER ID OF USER TO BAN}}`
