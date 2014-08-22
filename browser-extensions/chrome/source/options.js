if (preferences != undefined) {

	function UpdatePortInputs() {
		if (document.getElementById('textPort').value == preferences.default.Port) {
			document.getElementById('textPort').disabled = true;
			document.getElementById('checkboxUseCustomPort').checked = false;
		} else {
			document.getElementById('textPort').disabled = false;
			document.getElementById('checkboxUseCustomPort').checked = true;
		}
	}

	function EnableSaveButton() {
		var button = document.getElementById('savePreferences');
		if (button.disabled) {
			button.disabled = false;
		}
	}

	function DisableSaveButton() {
		var button = document.getElementById('savePreferences');
		if (!button.disabled) {
			button.disabled = true;
		}
	}

	function RestoreDefaults() {
		document.getElementById('textHost').value = preferences.default.Host;
		document.getElementById('textPort').value = preferences.default.Port;

		UpdatePortInputs();
		EnableSaveButton();
	}

	function RestoreLast() {
		preferences.get(function(items){
			document.getElementById('textHost').value = items.Host;
			document.getElementById('textPort').value = items.Port;

			UpdatePortInputs();
		});
	}

	function Save() {
		var host = document.getElementById('textHost').value,
			port = parseInt(document.getElementById('textPort').value);
		if (!preferences.validator.Host(host)) {
			alert("Введен не правильный хост");
		} else if (!preferences.validator.Port(port)) {
			alert("Введен не правильный порт");
		} else {
			preferences.set({
				'Host': host,
				'Port': port
			}, function(){
				alert("Настройки сохранены");
			});
		}
	}

	document.addEventListener('DOMContentLoaded', function() {

		RestoreLast();

		document.getElementById('restoreDefaults').addEventListener('click', function(){
			RestoreDefaults();
		});

		document.getElementById('restoreLast').addEventListener('click', function(){
			RestoreLast();
		});

		document.getElementById('checkboxUseCustomPort').addEventListener('click', function(){
			if (this.checked) {
				document.getElementById('textPort').value = '';
				document.getElementById('textPort').focus();
			} else {
				document.getElementById('textPort').value = preferences.default.Port;
			}
			UpdatePortInputs();
		});

		document.getElementById('savePreferences').addEventListener('click', function(){
			Save();
			DisableSaveButton();
		});

		document.getElementById('textHost').addEventListener('change', EnableSaveButton, false);

		document.getElementById('textPort').addEventListener('change', function() {
			UpdatePortInputs();
			EnableSaveButton();
		}, false);

	});

} else {
	console.log('undefined preferences');
}