package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var (
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

type result struct {
	name     string
	owners   string
	modified string
	icon     string
	link     string
}

// searchDrive use the Google Drive API to return a slice of *results for a given search query.
func searchDrive(arg string, config *oauth2.Config, token *oauth2.Token) ([]*result, error) {
	service, err := service(config, token)
	if err != nil {
		return nil, err
	}

	q := query(arg)

	r, err := service.Files.
		List().
		Q(q).
		Fields("files(name,mimeType,modifiedTime,owners,webViewLink)").
		PageSize(20).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Do()
	if err != nil {
		return nil, err
	}

	var results []*result
	for _, i := range r.Files {
		results = append(
			results,
			&result{
				i.Name,
				owners(i.Owners),
				i.ModifiedTime,
				typeIconMap[i.MimeType],
				i.WebViewLink,
			})
	}
	return results, nil
}

// service creates a new Google Drive Service.
// https://github.com/googleapis/google-api-go-client/blob/main/drive/v3/drive-gen.go
func service(config *oauth2.Config, token *oauth2.Token) (*drive.Service, error) {
	var context = context.TODO()
	client := config.Client(context, token)
	service, err := drive.NewService(context, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return service, err
}

// query sets the query string used for searching.
func query(arg string) string {
	var queries []string

	if arg != "" {
		queries = append(queries, fmt.Sprintf(queryMap["name"], arg))
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
