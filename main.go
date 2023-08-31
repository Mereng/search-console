package main

import (
	"bufio"
	"context"
	"flag"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	filePath := flag.String("f", "urls.txt", "file of list of pages")
	credsPath := flag.String("c", "creds.json", "file with creds")
	flag.Parse()

	creds, err := os.ReadFile(*credsPath)
	if err != nil {
		log.Fatalf("failed to read creds file: %v", err)
	}

	// Load the Google API credentials
	credentials, err := google.CredentialsFromJSON(context.Background(), creds, indexing.IndexingScope)
	if err != nil {
		log.Fatalf("Failed to get Google API credentials: %v", err)
	}

	client := oauth2.NewClient(context.Background(), credentials.TokenSource)
	indexingService, err := indexing.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Failed to create Indexing API service client: %v", err)
	}

	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalln(err)
	}
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		urlNotification := &indexing.UrlNotification{
			Url:  fileScanner.Text(),
			Type: "URL_UPDATED",
		}

		// Send the URL notification request
		resp, err := indexingService.UrlNotifications.Publish(urlNotification).Do()
		if err != nil {
			log.Fatalf("Failed to send URL notification %s: %v", urlNotification.Url, err)
		}

		log.Printf("%s indexed, status %d", urlNotification.Url, resp.HTTPStatusCode)
	}

}
