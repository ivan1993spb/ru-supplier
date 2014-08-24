
function urldecode(str) {
	return decodeURIComponent((str + '')
		.replace(/%(?![\da-f]{2})/gi, function() {
			return '%25';
		})
		.replace(/\+/g, '%20'));
}

function urlencode(str) {
	str = (str + '').toString();
	return encodeURIComponent(str)
		.replace(/!/g, '%21')
		.replace(/'/g, '%27')
		.replace(/\(/g, '%28')
		.replace(/\)/g, '%29')
		.replace(/\*/g, '%2A')
		.replace(/%20/g, '+');
}

function Query(query) {
	this.items = {};
	if (typeof query == 'string') {
		if (query.length > 0) {
			if (query[0] == '?') {
				query = query.substring(1);
			}
			var vars = query.split('&');
			for (var i = 0; i < vars.length; i++) {
				var pair = vars[i].split('=');
				pair[0] = urldecode(pair[0]);
				if (!(pair[0] in this.items)) {
					this.items[pair[0]] = [];
				}
				if (pair.length > 1) {
					this.items[pair[0]].push(urldecode(pair[1]));
				} else {
					this.items[pair[0]].push(null);
				}
			}
		}
	} else if (typeof query == 'object') {
		for (var key in query) {
			if (Object.prototype.toString.call(query[key]) === '[object Array]') {
				this.items[key] = query[key];
			} else {
				this.items[key] = [query[key]];
			}
		}
	} else {
		throw new Error('Required type string or object');
	}
}

Query.prototype.toString = function() {
	var pairs = [];
	for (var key in this.items) {
		for (var i = 0; i < this.items[key].length; i++) {
			pairs.push(urlencode(key) + '=' + urlencode(this.items[key][i]));
		}
	}
	if (pairs.length > 0) {
		return '?' + pairs.join('&');
	}
	return '';
};

Query.prototype.get = function(key) {
	if (key in this.items && this.items[key].length > 0) {
		return this.items[key][0];
	}
	return '';
};

Query.prototype.set = function(key, val) {
	if (Object.prototype.toString.call(val) !== '[object Array]') {
		val = [val];
	}
	this.items[key] = val;
};

function URL(urlStr) {
	if (typeof urlStr == 'string') {
		var a = document.createElement('a');
		a.href = urlStr;
		for (var key in a) {
			this[key] = a[key];
		}
	} else {
		throw new Error("Required type string");
	}
}

URL.ErrRequiredURL = new Error("Required type URL");

URL.prototype.query = function() {
	if (this.search.length > 0) {
		return new Query(this.search);
	} else {
		return new Query({});
	}
};

URL.prototype.toString = function() {
	var url = "";
	if (this.protocol.length > 0) {
		url += this.protocol;
		if (this.protocol[this.protocol.length-1] != ':') {
			url += ':';
		}
		url += '//';
	}
	if (this.username.length > 0) {
		url += this.username;
		if (this.password) {
			url += ':' + this.password;
		}
		url += '@';
	}
	url += this.hostname;
	if (this.port != '80' && this.port != '') {
		url += ':' + this.port;
	}
	if (this.pathname.length > 0) {
		if (this.pathname[0] != '/') {
			url += '/';
		}
		url += this.pathname;
	}
	if (this.search.length > 0) {
		if (this.search[0] != '?') {
			url += '?';
		}
		url += this.search;
	}
	if (this.hash.length > 0) {
		if (this.hash[0] != '#') {
			url += '#';
		}
		url += this.hash;
	}
	return url;
};

var URLZakupkiGovRu = {

	'ErrInvalidPath': new Error("Invalid path"),

	'makeProxyRSSURL': function(hostname, port, query) {
		var url = "http://" + hostname;
		if (port != '80' && port != '') {
			url += ':' + port;
		}
		url += "/rss";
		if (query instanceof Query) {
			url += query.toString();
		}
		return url;
	},

	'isQuickSearchRequest': function(path) {
		switch (path) {
		case "/epz/order/extendedsearch/search.html":
		case "/epz/order/orderCsvSettings/extendedSearch/download.html":
			return false;
		case "/epz/order/quicksearch/search.html":
		case "/epz/order/quicksearch/update.html":
		case "/epz/order/orderCsvSettings/quickSearch/download.html":
			return true;
		}
		throw this.ErrInvalidPath;
	},

	'isValidPath': function(path) {
		switch (path) {
		case "/epz/order/extendedsearch/search.html":
		case "/epz/order/quicksearch/search.html":
		case "/epz/order/quicksearch/update.html":
			return true;
		}
		return false;
	},

	'isValidSearchURL': function(url) {
		if (!(url instanceof URL)) {
			throw URL.ErrRequiredURL;
		}
		if (url.protocol != 'http:') {
			return false;
		}
		if (url.hostname != 'zakupki.gov.ru' && url.hostname != 'www.zakupki.gov.ru') {
			return false;
		}
		if (!this.isValidPath(url.pathname)) {
			return false;
		}
		return true;
	},

	'rewriteURL': function(url) {
		if (url instanceof URL) {
			switch (url.pathname) {
			case "/epz/order/extendedsearch/search.html":
				url.pathname = "/epz/order/orderCsvSettings/extendedSearch/download.html";
				break;
			case "/epz/order/quicksearch/search.html":
			case "/epz/order/quicksearch/update.html":
				url.pathname = "/epz/order/orderCsvSettings/quickSearch/download.html";
				break;
			default:
				throw this.ErrInvalidPath;
			}

			var query = url.query();

			if (this.isQuickSearchRequest(url.pathname)) {
				query.set("quickSearch", "true");
			} else {
				query.set("quickSearch", "false");
			}
			query.set("sortBy", "PUBLISH_DATE");
			query.set("sortDirection", "false");
			query.set("userId", "null");
			query.set("conf", "true;true;true;true;true;true;true;" +
				"true;true;true;true;true;true;true;true;true;true;");
			
			url.search = query.toString();
			
			return url;
		} else {
			throw URL.ErrRequiredURL;
		}
	}

};
