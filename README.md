Keywords: Golang, go, web

### Pythia - A web application that allows specific tag-based searching for answers to questions

Pythia is a web application written in Go that allows you to build up a database of questions and answers.  You can assign multiple tagsto each question and then use the search interface to search on one or more specific tags.  I currently use Pythia to hold questions and answers about the boardgame Advanced Squad Leader (ASL).  ASL is a very complex game.  The rulebook is several hundred pages.  When playing a game of ASL, it is not unusual to have to stop and lookup a specific rule many times.  Often, you must jump from one reference to another in the rulebook to finally get the answer you are looking for.  The goal of Pythia is to make a web application, that looks good on phones and tablets, that a person can use to key in specific targeted tags and get back a short list of answers that give him the answer he is looking for.  Hopefully, this will make rulebook lookups much less frequent.

### How to install

~~~
go get github.com/jameycribbs/pythia
~~~

- go get any dependencies
- go build pythia.go
- in the directory where you are going to run the pythia executable, create a "data" directory and two subdirectories "data/answers" and "data/users"
- copy the "1.json" file to the "data/users" directory
- run the pythia executable that you just built
- point your browser to http://localhost:8080


### How to use

To add questions and answers, you will need to be logged in as an admin. You can initially login as login "login" and password "password".  This test user is an admin user which will allow you to go create a real admin user.  Make sure you put "admin" in the level field when creating your own user.  Once you have created your own admin user, you need to go back and delete the test user.
 
Once you have added some records, anyone can go to the front page and key in one or more tags to search for answers.  Only records that have ALL of the tags that are being searched for will show up in the search results.

### Contributions welcome!

Pull requests/forks/bug reports all welcome, and please share your thoughts, questions and feature requests in the [Issues] section or via [Email].

[Email]: mailto:jamey.cribbs@gmail.com
[Issues]: https://github.com/jameycribbs/pythia/issues

