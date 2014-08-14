# Ru-supplier #

* [Russian documentation](docs/index.html)
* [Bitbucket repo](https://bitbucket.org/pushkin_ivan/ru-supplier)
* [Pushkin Ivan](mailto://pushkin13@bk.ru)

### Ru-supplier ###

* version 2.0
* is written on [Golang](http://golang.org)
* helps to monitor orders on zakupki.gov.ru;
* makes work with orders really faster;
* is simplest local proxy server with filter;
* works together with any rss client;
* works only on windows 7+ platforms;

### Ru-supplier can: ###

* to read csv stream from zakupki.gov.ru and parse orders;
* to filter orders by PCRE regular expressions;
* to form human friendly designed and fast readable rss feed with orders;
* to cache last order;

### Repo directories ###
See [ru-supplier source on bitbucket](https://bitbucket.org/pushkin_ivan/ru-supplier) if you are interested in [Golang](http://golang.org)

```
#!
./     - contains ru-supplier source files
./docs - contains html documentation files
./urls - contains url generator source files
```

### Compilation ###

You have to download Golang compiler [here](http://golang.org)

```
#!
C:\> cd path-to-ru-supplier
C:\path-to-ru-supplier> go build -ldflags="-H windowsgui"
C:\path-to-ru-supplier> cd urls
C:\path-to-ru-supplier\urls> go build -ldflags="-H windowsgui"
```
or

```
#!
C:\> cd path-to-ru-supplier
C:\path-to-ru-supplier> builder.bat
C:\path-to-ru-supplier> start build
```
then will be created build directory with program files

### How it works ###
1. go to the [http://zakupki.gov.ru](http://zakupki.gov.ru)
2. do any request
3. copy request url
4. convert request url into rss feed link by url generator
5. run ru-supplier
6. open your rss client
7. subscribe to rss feed
