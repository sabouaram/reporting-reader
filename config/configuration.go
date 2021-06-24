package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var authcode string

// HTTP Server and handler for getting the authcode passing as a query to Reporting-reader backend client after Resource owner authorization
func HandleAuthCode(w http.ResponseWriter, req *http.Request) {
	authcode = req.URL.Query().Get("code")
}
func init() {
	go func() {
		http.HandleFunc("/", HandleAuthCode)
		http.ListenAndServe("localhost:8080", nil)
	}()
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	time.Sleep(20 * time.Second)
	log.Println(authcode)
	tok, err := config.Exchange(context.TODO(), authcode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

type Config struct {
	AuthorizedHTTPClient *http.Client
	Username             string //String User@domain
}

//NewConfig exchange OAUTH credentianls for an access token and return the authorized http client based on the Scope defined in the func args
func NewConfig(filename string, username string) (*Config, error) {
	secret, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	conf, err := google.ConfigFromJSON(secret, gmail.GmailModifyScope)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	client := getClient(conf)
	log.Printf("OAuth2.0 Flow Succeeded: Granted Access for REPORTING-READER BACKEND")
	return &Config{
		AuthorizedHTTPClient: client,
		Username:             username,
	}, nil

}
