
if (preferences != undefined) {
	document.addEventListener('DOMContentLoaded', function(){
		chrome.tabs.query({
			active: true,
			lastFocusedWindow: true
		}, function(array_of_Tabs) {
			var tab = array_of_Tabs[0];
			var url = tab.url;
		});
	});
} else {
	console.log('undefined preferences');
}
