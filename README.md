
# ru-supplier #

This tool helps to monitor orders on zakupki.gov.ru. ru-supplier is simplest local proxy server with filter.

### ru-supplier can: ###

* to read csv stream from zakupki.gov.ru and parse orders;
* to filter orders by PCRE regular expressions;
* to form human friendly designed and fast readable rss feed with orders;
* to cache last order;

ru-supplier works together with any rss client.

// To get rid of the cmd window, instead run
// go build -ldflags="-H windowsgui"
