package main

import (
	"fmt"
	"flag"
	"net/http"
	"regexp"
	"encoding/json"
	"golang.org/x/net/html"
	"github.com/PuerkitoBio/goquery"
	"strconv"

	. "github.com/sbaildon/hn-json/types"
	"github.com/sbaildon/hn-json/validate"
)

const (
	HACKER_NEWS = "https://news.ycombinator.com"
	POSTS_PER_PAGE = 30
	FOCUS_ELEMENT = "table"
	FOCUS_CLASS = "itemlist"
)

var (
	postsToPrint *int = flag.Int("posts", POSTS_PER_PAGE, "How many posts to print")
)


func main() {
	flag.Parse()

	if (*postsToPrint > 100) || (*postsToPrint <= 0) {
		panic("Unsupported number of posts")
	}

	results := &Result{}

	/*
	 * Keep fetching pages until we've got enough posts
	 * to print
	 */
	page := 1
	for len(results.Posts) <= *postsToPrint {
		html := fetchHN(page)
		parsePosts(html, results)
		page = page+1
	}

	js, _ := json.Marshal(results.Posts[:*postsToPrint])
	fmt.Println(string(js))
}

/* Fetches the HN desired HN page, and returns
 * the entire DOM tree
 */
func fetchHN(pageNumber int) *html.Node {
	req, _ := http.NewRequest("GET", HACKER_NEWS, nil)
	query := req.URL.Query()
	query.Add("p", strconv.Itoa(pageNumber))
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Not 200")
	}

	root, err := html.Parse(res.Body)

	if err != nil {
		panic(err)
	}

	return root
}

/* Posts are rows in a table, taking the form
 * row1 -> title
 * row2 -> meta
 * row3 -> empty space
 * To find all posts, we have to look for all ".athing"s
 * which is the title, then look for the following sibling
 * to read the metadata
 */
func parsePosts(html *html.Node, results *Result) {
	doc := goquery.NewDocumentFromNode(html)

	element := fmt.Sprintf("%s.%s", FOCUS_ELEMENT, FOCUS_CLASS)

	posts := make(chan Post)
	done := make(chan bool)

	/* Read Posts from channel until it's closed
	 * by the publisher. Write posts to results
	 */
	go func() {
		for {
			post, more := <- posts
			if more {
				results.Posts = append(results.Posts, post)
			} else {
				done <- true
				return
			}
		}
	}()

	doc.Find(element).Find("tbody").Find(".athing").Each(func(i int, s *goquery.Selection) {
		heading := s.Nodes[0]
		meta := s.Next().Nodes[0]

		parsePost(heading, meta, posts)
	})

	close(posts)

	<-done
}

func parsePost(headingNode *html.Node, metaNode *html.Node, posts chan Post) {
	heading := &Heading{}
	parseHeading(headingNode, heading)

	meta := &Meta{}
	parseMeta(metaNode, meta)

	posts <- Post{
		Heading: heading,
		Meta: meta,
	}
}

func parseHeading(headingNode *html.Node, heading *Heading) {
	element := goquery.NewDocumentFromNode(headingNode)

	reg, _ := regexp.Compile("[^0-9]+")

	_rank := reg.ReplaceAllString(element.Find("span.rank").Text(), "")
	rank, _ := strconv.Atoi(_rank)
	if valid := validate.Rank(rank); !valid {
		panic(fmt.Sprintf("invalid rank %d", rank))
	}

	title := element.Find("a.storylink").Text()
	if valid := validate.Title(title); !valid {
		panic(fmt.Sprintf("invalid title %s", title))
	}

	uri, _ := element.Find("a.storylink").Attr("href")
	if valid := validate.URI(uri); !valid {
		panic(fmt.Sprintf("invalid uri %s", uri))
	}

	heading.Rank = rank
	heading.Title = title
	heading.Uri = uri
}

func parseMeta(metaNode *html.Node, meta *Meta) {
	element := goquery.NewDocumentFromNode(metaNode)

	reg, _ := regexp.Compile("[^0-9]+")

	_points := reg.ReplaceAllString(element.Find("span.score").Text(), "")
	points, _ := strconv.Atoi(_points)
	if valid := validate.Points(points); !valid {
		panic(fmt.Sprintf("invalid points %d", points))
	}

	author := element.Find("a.hnuser").Text()
	if valid := validate.Author(author); !valid {
		// job postings don't have an author, not
		// sure what the best way to fail is
		// so set it to an obvious value
		author = "internal_job_posting"
	}

	_comments := reg.ReplaceAllString(element.Find("td.subtext").Find("a").Last().Text(), "")
	comments, _ :=  strconv.Atoi(_comments)
	if valid := validate.Comments(comments); !valid {
		panic(fmt.Sprintf("invalid comments %d", comments))
	}

	meta.Points = points
	meta.Author = author
	meta.Comments = comments
}
