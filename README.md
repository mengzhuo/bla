# bla
another blog based on golang

![Screenshot](https://raw.githubusercontent.com/mengzhuo/bla_default_template/master/bla_screenshot.png)

## Highlights
* Fully static file serving
* WYSIWYG (No markdown or offline editing)
* Tag/Relation finding
* Golang blog style template
* Easy deploy with no php, no Mysql (only golang)

## Simple Tutorial

``` shell
go get -u github.com/mengzhuo/bla/cmd/bla
bla new
```

edit your configure file "bla.config.json" and just run bla

## How to add post to blog?
http://localhost:7080/.add

## How to edit post of blog?
http://localhost:7080/.edit?doc=<docname>

## How to quit? I don't want CSRF
http://localhost:7080/.quit


## TODO
* manauls?
* testcases
* media uploads
* automatically download template
