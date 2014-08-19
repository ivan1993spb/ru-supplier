# Ru-supplier #

* [Russian documentation](docs/index.html)
* [Bitbucket repo](https://bitbucket.org/pushkin_ivan/ru-supplier)
* [Pushkin Ivan](mailto://pushkin13@bk.ru)

### Ru-supplier ###

* author: Pushkin Ivan
* version 2.0
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
./docs               - contains html documentation
./src                - contains icons, manifests and other
./go-source          - contains .go source files
```

### Compilation ###

You have to download and install:
* [Golang compiler](http://golang.org/doc/install).You have to set up %GOPATH% env variable (something like this C:\gocode)
* [Git](http://git-scm.com/downloads). You have to install console cmd git version
* [Mercurial](http://mercurial.selenic.com/wiki/Download)

Then download this repo.

```
#!
C:\> cd path-to-ru-supplier
C:\path-to-ru-supplier> builder.bat
```
Then will be created and opened directory with program files.
After compilation you have to create local host in hosts file.

### How it works ###
1. go to the [http://zakupki.gov.ru](http://zakupki.gov.ru)
2. do any request
3. copy request url
4. convert request url into rss feed link by url generator
5. run ru-supplier (proxy)
6. open your rss client
7. subscribe to the rss feed
