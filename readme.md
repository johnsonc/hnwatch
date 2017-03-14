# hnwatch
The idea behind this little project is don't procrastinate on HN
and never miss a story, job, or project I'm interested in.

You can fetch all new stories, or stories containing particular terms
in the title. To define term, you can use plain text or regular expressions.

SQLite 3 is used to store fetched stories.
You will never get the same story for the same term twice.

Program runs in for cycle with pause between cycles,
default pause is 10 minutes. Please don't abuse HN :)

Fetch all HN news from homepage
```
./hnwatch -t "Ask HN"
```
Fetch HN news containing term(s) **Ask HN** in title
```
./hnwatch -t "Ask HN"
```
Fetch HN news containing term(s) **Ask HN** in title, repeat checking after 30 minutes (default is 10 minutes)
```
./hnwatch -t "Ask HN" -r 30
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
