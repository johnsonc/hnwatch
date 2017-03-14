// TODO: split code into multiple files
// 	 is premature optimization the root of all evil ??
package main

import (
	"bufio"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	version     = "0.1"
	baseURL     = "https://news.ycombinator.com/"
	commentsURL = "https://news.ycombinator.com/item?id="
)

var (
	emailTemplate = `<!doctype html>
	<html lang="en">
	<head>
	<meta charset="utf-8">
	<style>
	.story {margin-bottom:10px;}
	.title {font-weight:bold;}
	</style>
	</head>
	<body>
	%s
	</body>
	</html>`

	textTemplate = "* %s\n\t%s\n\n"
)


type Item struct {
	id    string
	URL   string
	title string
	time  int64
}

func b64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func readConfig(fileName, keyPrefix string) (map[string]string, error) {

	f, err := os.Open(fileName)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	config := make(map[string]string, 20)
	kvPatt := fmt.Sprintf(`(%s[^\s]+)[\s\t]+([^\n]+)`, keyPrefix)
	re := regexp.MustCompile(kvPatt)
	lines := bufio.NewScanner(f)
	ln := 0

	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		kv := re.FindStringSubmatch(line)
		ln++

		if kv != nil && strings.HasPrefix(kv[0], "#") == false {
			k := strings.TrimSpace(kv[1])
			v := strings.TrimSpace(kv[2])
			if _, ok := config[k]; ok {
				err := fmt.Errorf("Dup key: %s, line %d", k, ln)
				return nil, err
			}
			config[k] = v
		}
	}

	return config, nil
}

func email(config map[string]string, subject, body string) error {
	// TODO: TLS???
	// Set up plain authentication information.
	auth := smtp.PlainAuth("",
		config["smtp_user"],
		config["smtp_password"],
		config["smtp_server"],
	)

	to := []string{config["smtp_to_addr"]}
	date := time.Now().Format(time.RFC1123)
	msg := fmt.Sprintf(
		"MIME-Version: 1.0\r\n"+
			"To: <%s>\r\n"+
			"Subject: "+"=?UTF-8?B?"+"%s"+"?="+"\r\n"+
			"From: %s<%s>\r\n"+
			"Date: %s\r\n"+
			"Content-Type: text/html; charset=utf-8\r\n"+
			"Content-Transfer-Encoding: base64\r\n"+
			"\r\n"+
			"%s\r\n",
			to, b64(subject), config["smtp_from_name"],
			config["smtp_from_addr"], date, body)

	err := smtp.SendMail(config["smtp_server"]+":"+config["smtp_port"],
		auth, config["smtp_from_addr"], to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func fetchPage(url string, timeout time.Duration) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: timeout * time.Second}
	//TODO: rethink
	client.CheckRedirect =
		func(req *http.Request, via []*http.Request) error {
			url = req.URL.String()
			req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HNWatch; +https://github.com/vetelko/hnwatch)")
			req.Header.Set("Accept", "*/*")
			req.Header.Set("Accept-Language", "en-US,en;q=0.8,*;q=0.4")
			req.Header.Set("Accept-Charset", "UTF-8;q=0.8,*;q=0.7")
			return nil
		}

	reqt, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	reqt.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HNWatch; +https://github.com/vetelko/hnwatch)")
	reqt.Header.Set("Accept", "*/*")
	reqt.Header.Set("Accept-Language", "en-US,en;q=0.8,*;q=0.4")
	reqt.Header.Set("Accept-Charset", "UTF-8;q=0.8,*;q=0.7")

	resp, err := client.Do(reqt)
	if resp != nil {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {

		}
		defer resp.Body.Close()
		return string(data)
	}
	return ""
}

func internalURL(URL string) string {
	if URL[0:4] == "http" {
		return URL
	}
	return baseURL + URL
}

func (item *Item) parseItems(url, term string) []Item {
	var items []Item
	data := fetchPage(url, 10)
	re := regexp.MustCompile(`(?sU)<tr class='athing' id='(\d+)'>.+<a href="([^"]+)" class="storylink"[^>]*>(.+)<\/a>`)
	matches := re.FindAllStringSubmatch(data, -1)
	for _, match := range matches {
		term = strings.ToLower(term)
		re2 := regexp.MustCompile(`.*` + term + `.*`)
		if re2.MatchString(strings.ToLower(match[3])) {
			item.id = match[1]
			item.URL = internalURL(match[2])
			item.title = match[3]
			item.time = time.Now().Unix()
			items = append(items, *item)
		}
	}
	return items
}

func dbItemFilter(items []Item, term string) (outText, outHTML string) {
	db, _ := sql.Open("sqlite3", "hnwatch.db")
	defer db.Close()
	_, _ = db.Exec("CREATE TABLE IF NOT EXISTS item(term text, item_id text, title text, stamp text, url text)")
	for _, item := range items {
		var rowid int
		row := db.QueryRow("SELECT rowid FROM item WHERE url=? AND term=?", item.URL, term).Scan(&rowid)
		if row == sql.ErrNoRows {
			_, _ = db.Exec("INSERT INTO item VALUES(?, ?, ?, ?, ?)", term, item.id, item.title, item.time, item.URL)
			outText = outText + fmt.Sprintf(textTemplate, item.title, item.URL)
			// sometimes there is a space inserted into commentsURL variable value, I don't know why!?!
			// the comments link is broken in email. Something is wrong! Test more :)
			outHTML = outHTML + fmt.Sprintf(`<div class="story"><a class="title" href="%s">%s</a><br/>[<a href="%s%s">%s</a>]</div>`,
				item.URL, item.title, commentsURL, item.id, "comments")
		}
	}
	return outText, outHTML
}

func main() {
	var url = flag.String("u", "https://news.ycombinator.com/", "url containing items")
	var term = flag.String("t", "", "term(s) to find in item title, can be regexp")
	var c = flag.String("c", "hnwatch.cfg", "path to configuration file")
	var r = flag.Int("r", 10, "repeat checking after N minutes")
	var e = flag.Bool("e", true, "send output as email message")
	flag.Parse()

	var round = 1
	var cfg, _ = readConfig(*c, "")
	var item Item

	for {
		if *term != "" {
			fmt.Printf("Round: %d for term: %s\n", round, *term)
		} else {
			fmt.Printf("Round: %d\n", round)
		}

		items := item.parseItems(*url, *term)
		outText, outHTML := dbItemFilter(items, *term)

		fmt.Print(outText)
		if *e {
			email(cfg, "Test", fmt.Sprintf(emailTemplate, outHTML))
		}
		fmt.Printf("Next round in %d minute(s)\n\n", *r)
		round++
		time.Sleep(time.Duration(*r) * time.Minute)
	}
}
