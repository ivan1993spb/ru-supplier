
/*
 * This file contains default preferences, and tools for
 * editing chrome extension preferences
 * 
 * @author Pushkin Ivan
 */

function Preferences(prefs) {
	if (typeof prefs != 'object') {
		throw new Error('Required type is object');
	}
	this.items = prefs;
}

Preferences.ErrRequiredPreferences = new Error('Required instance of Preferences');

Preferences.prototype.empty = function() {
	for (var key in this.items) {
		if (Object.prototype.hasOwnProperty.call(this.items, key)) {
			return false;
		}
	}
	return true;
}

Preferences.prototype.equals = function(prefs) {
	if (prefs instanceof Preferences) {
		for (var key in this.items) {
			if (!(key in prefs.items) || this.items[key] != prefs.items[key]) {
				return false;
			}
		}
		for (var key in prefs.items) {
			if (!(key in this.items) || prefs.items[key] != this.items[key]) {
				return false;
			}
		}
		return true;
	} else {
		throw Preferences.ErrRequiredPreferences;
	}
}

Preferences.prototype.intersectFields = function(prefs) {
	if (prefs instanceof Preferences) {
		if (this.empty() || prefs.empty()) {
			return {};
		}
		var fields = [];
		for (var key in this.items) {
			if (key in prefs.items && this.items === prefs.items) {
				fields.push(key);
			}
		}
		return fields;
	} else {
		throw Preferences.ErrRequiredPreferences;
	}
}

Preferences.prototype.import = function(prefs) {
	if (prefs instanceof Preferences) {
		for (var key in prefs.items) {
			this.items[key] = prefs.items[key];
		}
	} else {
		throw Preferences.ErrRequiredPreferences;
	}
}

var RuSupplierPreferences = {

	'default': new Preferences({
		'Host': 'proxy-zakupki-gov-ru.local',
		'Port': 80
	}),

	'last': undefined,

	'init': function(prefs) {
		if (!(prefs instanceof Preferences)) {
			throw Preferences.ErrRequiredPreferences;
		}
		this.last = prefs;
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
		}
	},

	'get': function(callback) {
		chrome.storage.local.get(this['default']['items'], callback);
	},

	'set': function(prefs, callback) {
		if (!(prefs instanceof Preferences)) {
			throw Preferences.ErrRequiredPreferences;
		}
		if (typeof callback != 'function') {
			throw new Errpr('Second parameter must be callable');
		}

		// remove bad preference fields and throw on invalid value
		for (var key in prefs.items) {
			if (key in this['default']['items']) {
				if (key in this['validator']) {
					if (typeof this['validator'][key] == 'function') {
						if (!this['validator'][key](prefs.items[key])) {
							throw new Error('Invalid option value');
						}
					}
				}
			} else {
				delete prefs.items['key'];
			}
		}

		if (!prefs.empty()) {
			try {
				if (!prefs.equals(this.last)) {
					// import new preferences
					this.last.import(prefs);

					// if any preference like default remove from store
					var defFields = prefs.intersectFields(this['default']);
					if (defFields.length > 0) {
						for (i in defFields) {
							delete prefs.items[defFields[i]];
						}
						chrome.storage.local.remove(defFields);
					}

					// save preferences into browser storage
					if (!prefs.empty()) {
						chrome.storage.local.set(prefs.items, callback);
					} else {
						callback();
					}
				}
			} catch (e) {}
		}
	}

};
