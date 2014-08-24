
document.addEventListener('DOMContentLoaded', function(){
	chrome.tabs.query({
		active: true,
		lastFocusedWindow: true
	}, function(tabs) {
		var tab = tabs[0];
		RuSupplierPreferences.get(function(items){

			// init storage api
			window.RuSupplierPreferences.init(new Preferences(items));

			// generate and insert url
			var tabUrl = new URL(tab.url);
			if (URLZakupkiGovRu.isValidSearchURL(tabUrl)) {
				tabUrl = URLZakupkiGovRu.rewriteURL(tabUrl);
				document.getElementById('inputURL').value = URLZakupkiGovRu.makeProxyRSSURL(
					window.RuSupplierPreferences.last.items.Host,
					window.RuSupplierPreferences.last.items.Port,
					new Query({'url': tabUrl.toString()})
				);
			}
			document.getElementById('inputURL').addEventListener('click', function(){
				this.select();
			});

			// bind click listener
			document.getElementById('buttonCopy').addEventListener('click', function(){
				document.getElementById("inputURL").select();
				document.execCommand("Copy", false, null);
				window.close();
			});

		});
	});
});
