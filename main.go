package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"

	"github.com/spf13/cobra"
)

var (
	credentialsFile string
	tokenFile       string

	fromEmail    string
	toEmail      string
	subjectEmail string
	bodyEmail    string
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := tokenFile
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

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
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

func sendEmail(cmd *cobra.Command, args []string) error {
	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return err
	}
	client := getClient(config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	var message gmail.Message

	messageStr := []byte(
		fmt.Sprintf("From: %s\r\n", fromEmail) +
			fmt.Sprintf("To: %s\r\n", toEmail) +
			fmt.Sprintf("Subject: %s\r\n\r\n", subjectEmail) +
			bodyEmail)

	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	fmt.Println("Message sent!")
	return nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "email-notification",
		Short: "Send email using Google's Gmail API",
		Long: `Use Gmail's API with OAuth2 to send notification emails.

If you have a server you want to manage and send yourself notification emails
from time to time there's multiple ways you can go about it. One way is to
set up your own mail server. But that has quite a few heavy requirements (on
top of having to configure them), and your IP might just be blacklisted
anyway.

When researching I found out Gmail has an API you can use with OAuth2.
Meaning you can set up a new Gmail account and use their servers and email
service, which will pretty much always work and never be blacklisted, to send
any given notification email to whatever other account you want.`,
	}

	usr, _ := user.Current()
	dir := usr.HomeDir

	rootCmd.PersistentFlags().StringVar(
		&credentialsFile,
		"credentials",
		filepath.Join(dir, ".email-notification/credentials.json"),
		"OAuth2 Credentials JSON file",
	)

	rootCmd.PersistentFlags().StringVar(
		&tokenFile,
		"token",
		filepath.Join(dir, ".email-notification/token.json"),
		"OAuth2 Token JSON file",
	)

	cmdSend := &cobra.Command{
		Use:   "send",
		Short: "Send email command",
		RunE:  sendEmail,
	}

	cmdSend.Flags().StringVar(&fromEmail, "from", "", "From email")
	cmdSend.Flags().StringVar(&toEmail, "to", "", "To email")
	cmdSend.Flags().StringVar(&subjectEmail, "subject", "", "Email subject")
	cmdSend.Flags().StringVar(&bodyEmail, "body", "", "Email body")

	rootCmd.AddCommand(cmdSend)
	rootCmd.Execute()
}
