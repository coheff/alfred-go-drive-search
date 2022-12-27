package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	ac "github.com/coheff/al-co"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
	config = &oauth2.Config{
		ClientID:     Wf.Config.Get("client_id"),
		ClientSecret: Wf.Config.Get("client_secret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:1337/callback",
		Scopes:       []string{drive.DriveMetadataReadonlyScope},
	}
	queryMap = map[string]string{
		"name": "name contains '%s'",
		"full": "fullText contains '%s'",
	}
	typeIconMap = map[string]string{
		"application/vnd.google-apps.document":     "icons/doc.png",
		"application/vnd.google-apps.spreadsheet":  "icons/sheet.png",
		"application/vnd.google-apps.presentation": "icons/slide.png",
		"application/vnd.google-apps.form":         "icons/form.png",
		"application/pdf":                          "icons/pdf.png",
	}
)

// SearchDrive uses the Google Drive API to return a slice of *searchResult for a given search query.
func SearchDrive(q string) []*Result {
	// get cached token
	tok, err := ac.CachedToken(Kc)
	if err != nil {
		log.Printf("Error retrieving cached token; it might not exist: %v", err)

		// get new token
		tok, err = ac.NewToken(config)
		if err != nil {
			log.Fatalf("Error aquiring token: %v", err)
		}

		// store token
		err = ac.CacheToken(Kc, tok)
		if err != nil {
			log.Fatalf("Error storing token: %v", err)
		}
	}

	// create Google Drive service
	// https://github.com/googleapis/google-api-go-client/blob/main/drive/v3/drive-gen.go
	var ctx = context.TODO()
	cli := config.Client(ctx, tok)
	service, err := drive.NewService(ctx, option.WithHTTPClient(cli))
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	// make request to service
	resp, err := service.Files.List().
		Q(buildQuery(q)).
		Fields("files(name,mimeType,modifiedTime,owners,webViewLink)").
		PageSize(20).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Do()
	if err != nil {
		log.Fatalf("Error calling FilesService: %v", err)
	}

	var results []*Result
	for _, file := range resp.Files {
		results = append(
			results,
			&Result{
				file.Name,
				fmt.Sprintf("%s - %s", owners(file.Owners), file.ModifiedTime),
				file.WebViewLink,
				typeIconMap[file.MimeType],
			})
	}

	return results
}

// buildQuery sets the query string used for searching.
func buildQuery(q string) string {
	var queries []string

	if q != "" {
		queries = append(queries, fmt.Sprintf(queryMap["name"], q))
	}

	if opts.FullText != "" {
		queries = append(queries, fmt.Sprintf(queryMap["full"], opts.FullText))
	}

	return strings.Join(queries, " and ")
}

// owners parses display names from a slice of Users.
func owners(user []*drive.User) string {
	var owners []string

	for _, u := range user {
		owners = append(owners, u.DisplayName)
	}

	return strings.Join(owners, ", ")
}
