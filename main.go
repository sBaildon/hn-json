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
)

const (
	HACKER_NEWS = "https://news.ycombinator.com"
	POSTS_PER_PAGE = 30
	FOCUS_ELEMENT = "table"
	FOCUS_CLASS = "itemlist"
)

var (
	postsToPrint *int = flag.Int("posts", 30, "How many posts to print")
)


func main() {
	flag.Parse()

	if (*postsToPrint > 100) || (*postsToPrint <= 0) {
		panic("Unsupported number of posts")
	}

	html := fetchHN(1)
	results := &Result{}
	parsePosts(html, results)

	js, _ := json.Marshal(results.Posts[:*postsToPrint])

	fmt.Println(string(js))
}

func fetchHN(pageNumber int) *html.Node {
	req, _ := http.NewRequest("GET", HACKER_NEWS, nil)
	req.URL.Query().Add("p", strconv.Itoa(pageNumber))

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

func parsePosts(html *html.Node, results *Result) {
	doc := goquery.NewDocumentFromNode(html)

	element := fmt.Sprintf("%s.%s", FOCUS_ELEMENT, FOCUS_CLASS)

	posts := make(chan Post)
	done := make(chan bool)

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

	title := element.Find("a.storylink").Text()
	uri, _ := element.Find("a.storylink").Attr("href")


	heading.Rank = rank
	heading.Title = title
	heading.Uri = uri
}

func parseMeta(metaNode *html.Node, meta *Meta) {
	element := goquery.NewDocumentFromNode(metaNode)

	reg, _ := regexp.Compile("[^0-9]+")

	_points := reg.ReplaceAllString(element.Find("span.score").Text(), "")
	points, _ := strconv.Atoi(_points)

	author := element.Find("a.hnuser").Text()

	_comments := reg.ReplaceAllString(element.Find("td.subtext").Find("a").Last().Text(), "")
	comments, _ :=  strconv.Atoi(_comments)

	meta.Points = points
	meta.Author = author
	meta.Comments = comments
}
