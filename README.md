# hn-json

## Intro

Crawls Hacker News, and converts any posts into machine friendly JSON.

Why? because why not

## Usage

### Go

`go get github.com/sbaildon/hn-json`

`hn-json [--posts int]`

Output is pretty ugly. Consider using a tool such as `jq` for pretty-printing.

`hn-json | jq`

### Docker

`docker build -t hn-json .`

## Improvements

* Fetch HN pages concurrently. Concurrency existed in the beginning, but the JSON results were out of order, so concurrency was removed in `d77a10d`.  Find a way to preserve ordering and reimplement concurrency.
* Job postings don't have an author. Currently an obvious value is set as a replacement, maybe not the best idea
* There are too many magic-words when searching for DOM elements. Make these constants
* Error handling. Every error is currently handled with a `panic`. Not really appropriate.
* Sometimes the last post on a page is the same as the first post on the following page. Probably a HN caching thing?
* All posts on a page are parsed, even if only a subset is required. Limit parsing to total number of posts required
* Respect `robots.txt`. There's a crawl-delay specified that this project ignores
