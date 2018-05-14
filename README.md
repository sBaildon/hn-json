# hn-json

## Intro

Crawls Hacker News, and converts any posts into machine friendly JSON.

Why? Why not.

## Usage

### Go

`go get github.com/sbaildon/hn-json`

`hn-json [--posts int]`

Output format is not human friendly. Consider using a tool such as `jq` for pretty-printing.

`hn-json | jq`

Need Go? https://golang.org/doc/install

### Docker

`docker build -t hn-json .`

`docker run -it --rm hn-json [-- --posts int]`

Need docker? https://docs.docker.com/install/

## Libraries

`goquery` is the only non standard library used. Mainly because makes working
with DOM lookups much simpler than iterating with `html.Tokenizer`. The API
is pretty robust and enables useful features like filtering, siblings, map
functions etc.

## Improvements

* Fetch HN pages concurrently. Concurrency existed in the beginning, but the JSON results were out of order, so concurrency was removed in `d77a10d`.  Find a way to preserve ordering and reimplement concurrency.
* Job postings don't have an author. Currently an obvious value is set as a replacement, maybe not the best idea
* There are too many magic-words when searching for DOM elements. Make these constants
* Error handling. Every error is currently handled with a `panic`. Not really appropriate.
* Sometimes the last post on a page is the same as the first post on the following page. Probably a HN caching thing?
* All posts on a page are parsed, even if only a subset is required. Limit parsing to total number of posts required
* Respect `robots.txt`. There's a crawl-delay specified that this project ignores
