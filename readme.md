# About

[![Build Status](https://travis-ci.org/vetelko/hnwatch.svg?branch=master)](https://travis-ci.org/vetelko/hnwatch)

I actually started programming last year in my 43, so I appologize for any
inconvenience :) My first language is Golang. I also learn TCL even it's not
a mainstream language.

The idea behind this little project is don't procrastinate on HN and never
miss a story, job, or project I'm interested in.
You can fetch all new stories, or stories containing particular terms
in the title. To define term, you can use plain text or regular expressions.
SQLite 3 is used to store fetched stories.
You will never get the same story for the same term twice.

Program uses never-ending for{} cycle with pause between cycles,
default pause is 30 minutes according to HN robots.txt. Please don't abuse HN :)

### Build
```
$ git clone https://github.com/vetelko/hnwatch
$ cd hnwatch
$ go build -ldflags="-s -w"
```

### Config
Configuration is stored in hnwatch.cfg file. Currently there is only configuration
for sending e-mail alerts.

### Storage
SQLite3 database file hnwatch.db is created automatically on the first run of
the application

### Run
```
Usage of ./hnwatch:
  -c string
        path to configuration file (default "hnwatch.cfg")
  -e bool
        send output as email message (default true)
  -r int
        repeat checking after N minutes (default 30)
  -t string
        term(s) to find in item title, can be regexp
  -u string
        url containing items (default "https://news.ycombinator.com/")
```

### Examples
Fetch all HN news from homepage. Next check after 30 minutes.
```
./hnwatch
```

Fetch HN news containing term(s) **Ask HN** in title
```
./hnwatch -t "Ask HN"
```

Fetch HN news containing term(s) **Ask HN** in title, repeat checking after 60 minutes (default is 30 minutes)
```
./hnwatch -t "Ask HN" -r 60
```

Fetch HN news containing term(s) **google or microsoft** in title
```
./hnwatch -t "google|microsoft"
```

Fetch HN **job** news containing term(s) **google or devops** in title
```
./hnwatch -u https://news.ycombinator.com/jobs -t "google|devops"
```

Fetch HN **show** news containing term(s) **golang or project** in title
```
./hnwatch -u https://news.ycombinator.com/show -t "golang|project"
```

Fetch HN **show** news containing term(s) **golang or project** in title using **other.cfg** as config file
```
./hnwatch -u https://news.ycombinator.com/show -t "golang|project" -c other.cfg
```

Fetch HN **show** news containing term(s) **golang or project** in title, send results to email
```
./hnwatch -u https://news.ycombinator.com/show -t "golang|project" -e=true
```
