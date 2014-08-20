
/**
 * This file contains default preferences, and
 * tools for preferences editing
 */

var preferences = {

	'default': {
		'Host': 'proxy-zakupki-gov-ru.local',
		'Port': 80
	},

	'validator': {
		
		'Host': function(host) {
			if (typeof host == 'string') {
				return /^[a-z0-9\-\.]+$/i.test(host);
			}
			return false;
		},
		
		'Port': function(port) {
			if (typeof port == 'number') {
				if (port % 1 === 0 && port > 0 && port < 65536) {
					return true;
				}
			}
			return false;
		},

	},

	'get': function(callback) {
		chrome.storage.local.get(this.default, callback);
	},

	'set': function(prefs, callback) {
		if (typeof prefs == 'object') {
			for (var key in prefs) {
				if (key in this.default) {
					if (key in this.validator) {
						if (!this.validator[key](prefs[key])) {
							delete prefs[key];
						}
					}
				} else {
					delete prefs[key];
				}
			}
			var flag = false;
			for (var key in prefs) {
				if (hasOwnProperty.call(prefs, key)) {
					flag = true;
					break;
				}
			}
			if (flag) {
				chrome.storage.local.set(prefs, callback);
			}
		}
	},

};
