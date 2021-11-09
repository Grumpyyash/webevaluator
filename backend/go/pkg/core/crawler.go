package core

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/Aman-Codes/backend/go/pkg/stringset"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var DefaultHTTPTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:    100,
	MaxConnsPerHost: 1000,
	IdleConnTimeout: 30 * time.Second,
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateOnceAsClient},
}

type Crawler struct {
	C                   *colly.Collector
	LinkFinderCollector *colly.Collector
	Output              *Output
	// Data                *[]SpiderOutput

	subSet  *stringset.StringFilter
	awsSet  *stringset.StringFilter
	jsSet   *stringset.StringFilter
	urlSet  *stringset.StringFilter
	formSet *stringset.StringFilter

	site       *url.URL
	domain     string
	Input      string
	Quiet      bool
	JsonOutput bool

	filterLength_slice []int
}

type SpiderOutput struct {
	Input      string `json:"input"`
	Source     string `json:"source"`
	OutputType string `json:"type"`
	Output     string `json:"output"`
	StatusCode int    `json:"status"`
}

func NewCrawler(
	site *url.URL, jsonOutput bool, maxDepth int, concurrent int, delay int, randomDelay int,
	subs bool, proxy string, timeout int, noRedirect bool, randomUA string, outputFolder string) *Crawler {
	domain := GetDomain(site)
	if domain == "" {
		Logger.Error("Failed to parse domain")
		os.Exit(1)
	}
	Logger.Infof("Start crawling: %s", site)

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(maxDepth),
		colly.IgnoreRobotsTxt(),
	)

	// Setup http client
	client := &http.Client{}

	// Set proxy
	if proxy != "" {
		Logger.Infof("Proxy: %s", proxy)
		pU, err := url.Parse(proxy)
		if err != nil {
			Logger.Error("Failed to set proxy")
		} else {
			DefaultHTTPTransport.Proxy = http.ProxyURL(pU)
		}
	}

	// Set request timeout
	if timeout == 0 {
		Logger.Info("Your input timeout is 0. We will set it to 10 seconds")
		client.Timeout = 10 * time.Second
	} else {
		client.Timeout = time.Duration(timeout) * time.Second
	}

	// Disable redirect
	if noRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			nextLocation := req.Response.Header.Get("Location")
			Logger.Debugf("Found Redirect: %s", nextLocation)
			// Allow in redirect from http to https or in same hostname
			// We just check contain hostname or not because we set URLFilter in main collector so if
			// the URL is https://otherdomain.com/?url=maindomain.com, it will reject it
			if strings.Contains(nextLocation, site.Hostname()) {
				Logger.Infof("Redirecting to: %s", nextLocation)
				return nil
			}
			return http.ErrUseLastResponse
		}
	}

	// Set client transport
	client.Transport = DefaultHTTPTransport
	c.SetClient(client)

	// Set User-Agent
	switch ua := strings.ToLower(randomUA); {
	case ua == "mobi":
		extensions.RandomMobileUserAgent(c)
	case ua == "web":
		extensions.RandomUserAgent(c)
	default:
		c.UserAgent = ua
	}

	// Set referer
	extensions.Referer(c)

	// Init Output
	var output *Output
	if outputFolder != "" {
		filename := strings.ReplaceAll(site.Hostname(), ".", "_")
		output = NewOutput(outputFolder, filename)
	}

	// Set url whitelist regex
	reg := ""
	if subs {
		reg = site.Hostname()
	} else {
		reg = "(?:https|http)://" + site.Hostname()
	}

	sRegex := regexp.MustCompile(reg)
	c.URLFilters = append(c.URLFilters, sRegex)

	// Set Limit Rule
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: concurrent,
		Delay:       time.Duration(delay) * time.Second,
		RandomDelay: time.Duration(randomDelay) * time.Second,
	})
	if err != nil {
		Logger.Errorf("Failed to set Limit Rule: %s", err)
		os.Exit(1)
	}

	// default disallowed regex
	disallowedRegex := `(?i)\.(png|apng|bmp|gif|ico|cur|jpg|jpeg|jfif|pjp|pjpeg|svg|tif|tiff|webp|xbm|3gp|aac|flac|mpg|mpeg|mp3|mp4|m4a|m4v|m4p|oga|ogg|ogv|mov|wav|webm|eot|woff|woff2|ttf|otf|css|pdf)(?:\?|#|$)`
	c.DisallowedURLFilters = append(c.DisallowedURLFilters, regexp.MustCompile(disallowedRegex))

	linkFinderCollector := c.Clone()
	// Try to request as much as Javascript source and don't care about domain.
	// The result of link finder will be send to Link Finder Collector to check is it working or not.
	linkFinderCollector.URLFilters = nil

	// var data *[]SpiderOutput

	return &Crawler{
		C:                   c,
		LinkFinderCollector: linkFinderCollector,
		site:                site,
		Input:               site.String(),
		JsonOutput:          jsonOutput,
		domain:              domain,
		Output:              output,
		// Data:                data,
		urlSet:  stringset.NewStringFilter(),
		subSet:  stringset.NewStringFilter(),
		jsSet:   stringset.NewStringFilter(),
		formSet: stringset.NewStringFilter(),
		awsSet:  stringset.NewStringFilter(),
	}
}

func (crawler *Crawler) feedLinkfinder(jsFileUrl string, OutputType string, source string) {

	if !crawler.jsSet.Duplicate(jsFileUrl) {
		// outputFormat := fmt.Sprintf("[%s] - %s", OutputType, jsFileUrl)
		var outputFormat string

		if crawler.JsonOutput {
			sout := SpiderOutput{
				Input:      crawler.Input,
				Source:     source,
				OutputType: OutputType,
				Output:     jsFileUrl,
			}
			crawler.Output.WriteToList(sout)
			if data, err := jsoniter.MarshalToString(sout); err == nil {
				outputFormat = data
				fmt.Println(outputFormat)
			}

		} else if !crawler.Quiet {
			fmt.Println(outputFormat)
		}

		if crawler.Output != nil {
			crawler.Output.WriteToFile(outputFormat)
		}

		// If JS file is minimal format. Try to find original format
		if strings.Contains(jsFileUrl, ".min.js") {
			originalJS := strings.ReplaceAll(jsFileUrl, ".min.js", ".js")
			_ = crawler.LinkFinderCollector.Visit(originalJS)
		}

		// Send Javascript to Link Finder Collector
		_ = crawler.LinkFinderCollector.Visit(jsFileUrl)

	}
}

func (crawler *Crawler) Start(linkfinder bool) {
	// Setup Link Finder
	if linkfinder {
		crawler.setupLinkFinder()
	}

	// Handle url
	crawler.C.OnHTML("[href]", func(e *colly.HTMLElement) {
		urlString := e.Request.AbsoluteURL(e.Attr("href"))
		urlString = FixUrl(crawler.site, urlString)
		if urlString == "" {
			return
		}
		if !crawler.urlSet.Duplicate(urlString) {
			outputFormat := fmt.Sprintf("[href] - %s", urlString)
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "form",
					Output:     urlString,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
					fmt.Println(outputFormat)
				}
			} else if !crawler.Quiet {
				fmt.Println(outputFormat)
			}
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}
			_ = e.Request.Visit(urlString)
		}
	})

	// Handle form
	crawler.C.OnHTML("form[action]", func(e *colly.HTMLElement) {
		formUrl := e.Request.URL.String()
		if !crawler.formSet.Duplicate(formUrl) {
			outputFormat := fmt.Sprintf("[form] - %s", formUrl)
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "form",
					Output:     formUrl,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
					fmt.Println(outputFormat)
				}
			} else if !crawler.Quiet {
				fmt.Println(outputFormat)
			}
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}

		}
	})

	// Find Upload Form
	uploadFormSet := stringset.NewStringFilter()
	crawler.C.OnHTML(`input[type="file"]`, func(e *colly.HTMLElement) {
		uploadUrl := e.Request.URL.String()
		if !uploadFormSet.Duplicate(uploadUrl) {
			outputFormat := fmt.Sprintf("[upload-form] - %s", uploadUrl)
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "upload-form",
					Output:     uploadUrl,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
					fmt.Println(outputFormat)
				}
			} else if !crawler.Quiet {
				fmt.Println(outputFormat)
			}
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}
		}

	})

	// Handle js files
	crawler.C.OnHTML("[src]", func(e *colly.HTMLElement) {
		jsFileUrl := e.Request.AbsoluteURL(e.Attr("src"))
		jsFileUrl = FixUrl(crawler.site, jsFileUrl)
		if jsFileUrl == "" {
			return
		}

		fileExt := GetExtType(jsFileUrl)
		if fileExt == ".js" || fileExt == ".xml" || fileExt == ".json" {
			crawler.feedLinkfinder(jsFileUrl, "javascript", "body")
		}
	})

	crawler.C.OnResponse(func(response *colly.Response) {
		respStr := DecodeChars(string(response.Body))

		if len(crawler.filterLength_slice) == 0 || !contains(crawler.filterLength_slice, len(respStr)) {

			// Verify which link is working
			u := response.Request.URL.String()
			// outputFormat := fmt.Sprintf("[url] - [code-%d] - %s", response.StatusCode, u)
			var outputFormat string
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "url",
					StatusCode: response.StatusCode,
					Output:     u,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
				}
			}
			fmt.Println(outputFormat)
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}
			if InScope(response.Request.URL, crawler.C.URLFilters) {
				crawler.findSubdomains(respStr)
				crawler.findAWSS3(respStr)
			}
		}
	})

	crawler.C.OnError(func(response *colly.Response, err error) {
		Logger.Debugf("Error request: %s - Status code: %v - Error: %s", response.Request.URL.String(), response.StatusCode, err)
		/*
			1xx Informational
			2xx Success
			3xx Redirection
			4xx Client Error
			5xx Server Error
		*/
		if response.StatusCode == 404 || response.StatusCode == 429 || response.StatusCode < 100 || response.StatusCode >= 500 {
			return
		}

		u := response.Request.URL.String()
		// outputFormat := fmt.Sprintf("[url] - [code-%d] - %s", response.StatusCode, u)
		var outputFormat string
		if crawler.JsonOutput {
			sout := SpiderOutput{
				Input:      crawler.Input,
				Source:     "body",
				OutputType: "url",
				StatusCode: response.StatusCode,
				Output:     u,
			}
			crawler.Output.WriteToList(sout)
			if data, err := jsoniter.MarshalToString(sout); err == nil {
				outputFormat = data
				fmt.Println(outputFormat)
			}
		}
		// else if crawler.Quiet {
		// 	fmt.Println(u)
		// } else {
		// 	fmt.Println(outputFormat)
		// }

		if crawler.Output != nil {
			crawler.Output.WriteToFile(outputFormat)
		}
	})

	err := crawler.C.Visit(crawler.site.String())
	if err != nil {
		Logger.Errorf("Failed to start %s: %s", crawler.site.String(), err)
	}
}

// Find subdomains from response
func (crawler *Crawler) findSubdomains(resp string) {
	subs := GetSubdomains(resp, crawler.domain)
	for _, sub := range subs {
		if !crawler.subSet.Duplicate(sub) {
			outputFormat := fmt.Sprintf("[subdomains] - %s", sub)

			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "subdomain",
					Output:     sub,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
				}
				fmt.Println(outputFormat)
			} else if !crawler.Quiet {
				outputFormat = fmt.Sprintf("[subdomains] - http://%s", sub)
				fmt.Println(outputFormat)
				outputFormat = fmt.Sprintf("[subdomains] - https://%s", sub)
				fmt.Println(outputFormat)
			}
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}
		}
	}
}

// Find AWS S3 from response
func (crawler *Crawler) findAWSS3(resp string) {
	aws := GetAWSS3(resp)
	for _, e := range aws {
		if !crawler.awsSet.Duplicate(e) {
			outputFormat := fmt.Sprintf("[aws-s3] - %s", e)
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     "body",
					OutputType: "aws",
					Output:     e,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
				}
			}
			fmt.Println(outputFormat)
			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}
		}
	}
}

// Setup link finder
func (crawler *Crawler) setupLinkFinder() {
	crawler.LinkFinderCollector.OnResponse(func(response *colly.Response) {
		if response.StatusCode == 404 || response.StatusCode == 429 || response.StatusCode < 100 {
			return
		}

		respStr := string(response.Body)

		if len(crawler.filterLength_slice) == 0 || !contains(crawler.filterLength_slice, len(respStr)) {

			// Verify which link is working
			u := response.Request.URL.String()
			// outputFormat := fmt.Sprintf("[url] - [code-%d] - %s", response.StatusCode, u)
			var outputFormat string
			if crawler.JsonOutput {
				sout := SpiderOutput{
					Input:      crawler.Input,
					Source:     response.Request.URL.String(),
					OutputType: "linkfinder",
					Output:     u,
					StatusCode: response.StatusCode,
				}
				crawler.Output.WriteToList(sout)
				if data, err := jsoniter.MarshalToString(sout); err == nil {
					outputFormat = data
				}
			}
			// fmt.Println(outputFormat)

			if crawler.Output != nil {
				crawler.Output.WriteToFile(outputFormat)
			}

			if InScope(response.Request.URL, crawler.C.URLFilters) {

				crawler.findSubdomains(respStr)
				crawler.findAWSS3(respStr)

				paths, err := LinkFinder(respStr)
				if err != nil {
					Logger.Error(err)
					return
				}

				currentPathURL, err := url.Parse(u)
				currentPathURLerr := false
				if err != nil {
					currentPathURLerr = true
				}

				for _, relPath := range paths {
					var outputFormat string
					// JS Regex Result
					if crawler.JsonOutput {
						sout := SpiderOutput{
							Input:      crawler.Input,
							Source:     response.Request.URL.String(),
							OutputType: "linkfinder",
							Output:     relPath,
						}
						crawler.Output.WriteToList(sout)
						if data, err := jsoniter.MarshalToString(sout); err == nil {
							outputFormat = data
						}
					} else if !crawler.Quiet {
						outputFormat = fmt.Sprintf("[linkfinder] - [from: %s] - %s", response.Request.URL.String(), relPath)
					}
					fmt.Println(outputFormat)

					if crawler.Output != nil {
						crawler.Output.WriteToFile(outputFormat)
					}
					rebuildURL := ""
					if !currentPathURLerr {
						rebuildURL = FixUrl(currentPathURL, relPath)
					} else {
						rebuildURL = FixUrl(crawler.site, relPath)
					}

					if rebuildURL == "" {
						continue
					}

					// Try to request JS path
					// Try to generate URLs with main site
					fileExt := GetExtType(rebuildURL)
					if fileExt == ".js" || fileExt == ".xml" || fileExt == ".json" || fileExt == ".map" {
						crawler.feedLinkfinder(rebuildURL, "linkfinder", "javascript")
					} else if !crawler.urlSet.Duplicate(rebuildURL) {

						if crawler.JsonOutput {
							sout := SpiderOutput{
								Input:      crawler.Input,
								Source:     response.Request.URL.String(),
								OutputType: "linkfinder",
								Output:     rebuildURL,
							}
							crawler.Output.WriteToList(sout)
							if data, err := jsoniter.MarshalToString(sout); err == nil {
								outputFormat = data
							}
						} else if !crawler.Quiet {
							outputFormat = fmt.Sprintf("[linkfinder] - %s", rebuildURL)
						}

						fmt.Println(outputFormat)

						if crawler.Output != nil {
							crawler.Output.WriteToFile(outputFormat)
						}
						_ = crawler.C.Visit(rebuildURL)
					}

					// Try to generate URLs with the site where Javascript file host in (must be in main or sub domain)

					urlWithJSHostIn := FixUrl(crawler.site, relPath)
					if urlWithJSHostIn != "" {
						fileExt := GetExtType(urlWithJSHostIn)
						if fileExt == ".js" || fileExt == ".xml" || fileExt == ".json" || fileExt == ".map" {
							crawler.feedLinkfinder(urlWithJSHostIn, "linkfinder", "javascript")
						} else {
							if crawler.urlSet.Duplicate(urlWithJSHostIn) {
								continue
							} else {

								if crawler.JsonOutput {
									sout := SpiderOutput{
										Input:      crawler.Input,
										Source:     response.Request.URL.String(),
										OutputType: "linkfinder",
										Output:     urlWithJSHostIn,
									}
									crawler.Output.WriteToList(sout)
									if data, err := jsoniter.MarshalToString(sout); err == nil {
										outputFormat = data
									}
								} else if !crawler.Quiet {
									outputFormat = fmt.Sprintf("[linkfinder] - %s", urlWithJSHostIn)
								}
								fmt.Println(outputFormat)

								if crawler.Output != nil {
									crawler.Output.WriteToFile(outputFormat)
								}
								_ = crawler.C.Visit(urlWithJSHostIn) //not print care for lost link
							}
						}

					}

				}
			}
		}
	})
}
