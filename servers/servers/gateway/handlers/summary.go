package handlers

import (
	"io"
	"golang.org/x/net/html"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	"errors"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()
	if len(query) == 0 {
		http.Error(w, "Missing Query Paramter(s)", http.StatusBadRequest)
		return
	}
	fetchedHTML, err := fetchHTML(query.Get("url"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	summary, err := extractSummary(query.Get("url"), fetchedHTML)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fetchedHTML.Close()
	finalErr := json.NewEncoder(w).Encode(summary)
	if finalErr != nil {
		http.Error(w, finalErr.Error(), 500)
		return
	}

}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	response, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		err = errors.New("Status code above 400")
		return nil, err
	} else if !strings.HasPrefix(response.Header.Get("Content-type"), "text/html") {
		err = errors.New("Content type is incorrect")
		return nil, err
	}
	return response.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
	tokenizer := html.NewTokenizer(htmlStream)
	summaryData := &PageSummary{}
	var imgArray []*PreviewImage
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if token.Data == "meta" {
				implicitURL, success := getMetaContent(token, "og:image")
				if success {
					imgArray = append(imgArray, &PreviewImage{})
					if implicitURL != "" {
						if strings.Contains(implicitURL, "://") {
							imgArray[len(imgArray) - 1].URL = implicitURL
						} else {
							imgArray[len(imgArray) - 1].URL = updateURL(pageURL, implicitURL)
						}
					}
				}
				explicitURL, success := getMetaContent(token, "og:image:url")
				if success {
					if strings.Contains(explicitURL, "://") {
						imgArray[len(imgArray) - 1].URL = explicitURL
					} else {
						imgArray[len(imgArray) - 1].URL = updateURL(pageURL, explicitURL)
					}
				}
				safeURL, success := getMetaContent(token, "og:image:secure_url")
				if success {
					imgArray[len(imgArray) - 1].SecureURL = safeURL
				}
				imgType, success := getMetaContent(token, "og:image:type")
				if success {
					imgArray[len(imgArray) - 1].Type = imgType
				}
				width, success := getMetaContent(token, "og:image:width")
				if success && len(width) != 0 {
					widthNum, err := strconv.Atoi(width)
					if err != nil {
						break
						//log.Fatal(err)
					}
					imgArray[len(imgArray) - 1].Width = widthNum
				}
				height, success := getMetaContent(token, "og:image:height")
				if success && len(height) != 0 {
					heightNum, err := strconv.Atoi(height)
					if err != nil {
						break
						//log.Fatal(err)
					}
					imgArray[len(imgArray) - 1].Height = heightNum
				}
				altTag, success := getMetaContent(token, "og:image:alt")
				if success {
					imgArray[len(imgArray) - 1].Alt = altTag
				}
				Type, success := getMetaContent(token, "og:type")
				if success {
					summaryData.Type = Type
				}
				URL, success := getMetaContent(token, "og:url")
				if success {
					summaryData.URL = URL
				}
				Title, success := getMetaContent(token, "og:title")
				if success {
					summaryData.Title = Title
				}
				siteName, success := getMetaContent(token, "og:site_name")
				if success {
					summaryData.SiteName = siteName
				}
				description, success := getMetaContent(token, "og:description")
				if success {
					summaryData.Description = description
				} else {
					tempDescription, success := getMetaContent(token, "description")
					if success && summaryData.Description == "" {
						summaryData.Description = tempDescription
					}
				}
				author, success := getMetaContent(token, "author")
				if success {
					summaryData.Author = author
				}
				keywords, success := getMetaContent(token, "keywords")
				if success {
					keywordsArray := strings.Split(keywords, ",")
					var trimmedKeywords []string
					for _, word := range keywordsArray {
						trimmedKeywords = append(trimmedKeywords, strings.TrimSpace(word))
					}
					summaryData.Keywords = trimmedKeywords
				}
			}
			if token.Data == "link" {
				for _, attribute := range token.Attr {
					if attribute.Key == "rel" && attribute.Val == "icon" {
						iconImagePreview := getIcon(token, pageURL)
						summaryData.Icon = iconImagePreview
					}
				}
			}
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken && summaryData.Title == "" {
					summaryData.Title = tokenizer.Token().Data
				}
			}
		}
	}
	if len(imgArray) != 0 {
		summaryData.Images = imgArray
	}
	return summaryData, nil
}

func getMetaContent(token html.Token, propertyName string) (string, bool) {
	givenAttribute := false
	givenContent := ""
	for _, attribute := range token.Attr {
		if (attribute.Key == "property" || attribute.Key == "name") && attribute.Val == propertyName {
			givenAttribute = true
		}
		if attribute.Key == "content" {
			givenContent = attribute.Val
		}
	}
	return givenContent, givenAttribute
}

func updateURL(pageURL string, implicitURL string) string {
	tempURL := strings.Split(pageURL, "/")
	tempURL[len(tempURL) - 1] = strings.Trim(implicitURL, "/")
	resultURL := strings.Join(tempURL, "/")
	return resultURL
}

func getIcon(token html.Token, pageURL string) *PreviewImage {
	iconImagePreview := new(PreviewImage)
	for _, attribute := range token.Attr {
		if attribute.Key == "href" {
			if strings.Contains(attribute.Val, "://") {
				iconImagePreview.URL = attribute.Val
			} else {
				iconImagePreview.URL = updateURL(pageURL, attribute.Val)
			}
		}
		if attribute.Key == "type" {
			iconImagePreview.Type = attribute.Val
		}
		if attribute.Key == "sizes" {
			if attribute.Val != "any" && len(attribute.Val) != 0 {
				heightWidth := strings.Split(attribute.Val, "x")
				heightNum, err := strconv.Atoi(heightWidth[0])
				if err != nil {
					break
					//log.Fatal(err)
				}
				widthNum, err := strconv.Atoi(heightWidth[1])
				if err != nil {
					 break
					 //log.Fatal(err)
				}
				iconImagePreview.Height = heightNum
				iconImagePreview.Width = widthNum
			}
		}
	}
	return iconImagePreview
}
