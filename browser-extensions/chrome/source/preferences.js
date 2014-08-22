
/**
 * This file contains default preferences, and
 * tools for editing preferences
 */

if (chrome.storage != undefined) {
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
				// remove bad preference fields
				for (var key in prefs) {
					if (key in this.default) {
						// remove invalid preferences if there is validator
						if (key in this.validator) {
							if (!this.validator[key](prefs[key])) {
								delete prefs[key];
							}
						}
					} else {
						delete prefs[key];
					}
				}
				if (!this.utils.empty(prefs)) {
					// if any preference like default remove from store
					var likeDefaultFields = this.utils.keys(
						this.utils.intersectObjects(this.default, prefs));
					for (key in prefs) {
						if (prefs[key] == this.default[key]) {
							likeDefaultFields.push(key);
							delete prefs[key];
						}
					}
					if (likeDefaultFields.length > 0) {
						chrome.storage.local.remove(likeDefaultFields);
					}
					// save preferences
					if (!this.utils.empty(prefs)) {
						chrome.storage.local.set(prefs, callback);
					} else {
						callback();
					}
				}
			}
		},

		'utils': {
			'empty': function(obj) {
				if (obj == null) {
					return true;
				}
				if (obj.length > 0) {
					return false;
				}
				if (obj.length === 0) {
					return true;
				}
				for (var key in obj) {
					if (Object.prototype.hasOwnProperty.call(obj, key)) {
						return false;
					}
				}
				return true;
			},

			'intersectObjects': function() {
				if (arguments.length == 0) {
					return {};
				}
				var objects = [];
				for (var i = 0; i < arguments.length; i++) {
					if (typeof arguments[i] == 'object') {
						objects.push(arguments[i]);
					}
				}
				if (objects.length == 0) {
					return {};
				}

				var result = objects[0];
				if (this.empty(result)) {
					return {};
				}
				for (var i = 1; i < objects.length; i++) {
					if (this.empty(objects[i])) {
						return {};
					}
					for (var key in result) {
						if (key in objects[i]) {
							if (result[key] === objects[i][key]) {
								continue;
							}
						}
						delete result[key];
						if (this.empty(result)) {
							return {};
						}
					}
				}
				return result;
			},

			'keys': function(obj) {
				var keys = [];
				if (typeof obj == 'object') {
					for (var k in obj) {
						keys.push(k);
					}
				}
				return keys;
			}
		}

	};
} else {
	console.log('undefined chrome.storage');
}