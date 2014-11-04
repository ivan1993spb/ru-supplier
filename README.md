# Ru-supplier #

* [Russian documentation](docs/index.html)
* [Bitbucket repo](https://bitbucket.org/pushkin_ivan/ru-supplier)
* [Pushkin Ivan](mailto://pushkin13@bk.ru)

### Ru-supplier ###

* author: Pushkin Ivan
* version 2.1
* is written on [Golang](http://golang.org)
* helps to monitor orders on zakupki.gov.ru
* makes work with orders really faster
* is simplest local proxy server with filter
* works together with any rss client
* works only on windows 7+ platforms

### Ru-supplier can: ###

* to read csv stream from zakupki.gov.ru and parse orders
* to filter orders by PCRE regular expressions
* to form human friendly designed and fast readable rss feed with orders
* to cache last order

### Repo directories ###
See [ru-supplier source on bitbucket](https://bitbucket.org/pushkin_ivan/ru-supplier) if you are interested in [Golang](http://golang.org)

```
#!
./browser-extensions - contains extensions for url generation
./docs               - contains russian documentation
./go-source          - contains .go source files
./src                - another files
```

### Compilation ###

You have to download and install:

* Golang [Golang compiler](http://golang.org/doc/install). You have to set up %GOPATH% env variable (something like this C:\gocode)
* [Git](http://git-scm.com/downloads). You have to install console git version
* [Mercurial](http://mercurial.selenic.com/wiki/Download)

Then download this repo.

```
#!
C:\> cd path-to-ru-supplier
C:\path-to-ru-supplier> make.bat
```

After compilation you have to create local host in [hosts file](https://ru.wikipedia.org/wiki/Hosts).

### How it works ###
1. go to the [http://zakupki.gov.ru](http://zakupki.gov.ru)
2. do any request
3. convert request url into rss feed link by browser extension
4. run ru-supplier
5. open your rss client
6. subscribe to the rss feed
