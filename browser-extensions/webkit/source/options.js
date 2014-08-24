
(function() {

	function UpdateForm() {
		// update port input and use-custom-port checkbox
		with (document) {
			if (getElementById('textPort').value == RuSupplierPreferences.default.items.Port) {
				getElementById('textPort').disabled = true;
				getElementById('checkboxUseCustomPort').checked = false;
			} else {
				getElementById('textPort').disabled = false;
				getElementById('checkboxUseCustomPort').checked = true;
			}
		}

		var currentPreferences = new Preferences({
			Host: document.getElementById('textHost').value,
			Port: document.getElementById('textPort').value
		});

		// update restore-default-preferences button
		if (currentPreferences.equals(RuSupplierPreferences.default)) {
			document.getElementById('restoreDefaults').disabled = true;
		} else {
			document.getElementById('restoreDefaults').disabled = false;
		}

		// update restore-last-preferences button and save button
		if (currentPreferences.equals(RuSupplierPreferences.last)) {
			document.getElementById('restoreLast').disabled = true;
			document.getElementById('savePreferences').disabled = true;
		} else {
			document.getElementById('restoreLast').disabled = false;
			document.getElementById('savePreferences').disabled = false;
		}
	}

	function RestoreDefaults() {
		with (RuSupplierPreferences.default.items) {
			document.getElementById('textHost').value = Host;
			document.getElementById('textPort').value = Port;
		}
		UpdateForm();
	}

	function RestoreLast() {
		with (RuSupplierPreferences.last.items) {
			document.getElementById('textHost').value = Host;
			document.getElementById('textPort').value = Port;
		}
		UpdateForm();
	}

	function Save() {
		var host = document.getElementById('textHost').value,
			port = parseInt(document.getElementById('textPort').value);
		if (!RuSupplierPreferences.validator.Host(host)) {
			alert("Введен не правильный хост");
		} else if (!RuSupplierPreferences.validator.Port(port)) {
			alert("Введен не правильный порт");
		} else {
			RuSupplierPreferences.set(new Preferences({
				'Host': host,
				'Port': port
			}), function(){
				UpdateForm();
				alert("Настройки сохранены");
			});
		}
	}

	function InitForm() {
		RestoreLast();
		UpdateForm();
	}

	function BindListeners() {
		document.getElementById('restoreDefaults').addEventListener('click', RestoreDefaults);
		document.getElementById('restoreLast').addEventListener('click', RestoreLast);
		document.getElementById('savePreferences').addEventListener('click', Save);
		document.getElementById('checkboxUseCustomPort').addEventListener('click', function(){
			if (this.checked) {
				document.getElementById('textPort').value = '';
			} else {
				document.getElementById('textPort').value = RuSupplierPreferences.default.items.Port;
			}
			UpdateForm();
		});
		document.getElementById('textHost').addEventListener('change', UpdateForm, false);
		document.getElementById('textPort').addEventListener('change', UpdateForm, false);
	}
	
	// initialization
	// wait for content loading
	document.addEventListener('DOMContentLoaded', function(){
		// load last preferences
		RuSupplierPreferences.get(function(items){
			// init storage api
			window.RuSupplierPreferences.init(new Preferences(items));
			// load last preferences and update form buttons
			InitForm();
			// bind event listeners
			BindListeners();
		});
	});

})();
