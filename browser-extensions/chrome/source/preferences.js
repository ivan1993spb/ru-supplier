
/**
 * This file contains default preferences
 */

var preferences = {
	'default': {
		"Host": "proxy-zakupki-gov-ru.local",
		"Port": 80
	},
	'get': function(callback) {
		chrome.storage.local.get(this.default, callback);
	},
	'set': function(prefs) {
		//...
		// chrome.storage.local.set({
		// 	favoriteColor: color,
		// 	likesColor: likesColor
		// })
	}
}
