
/**
 * 
 * Show pageAction icon when was loaded search page on zakupki.gov.ru
 * 
 */
chrome.webNavigation.onCommitted.addListener(function(e) {
	chrome.pageAction.show(e.tabId);
}, {
	url: [
		{
			hostSuffix: 'zakupki.gov.ru',
			pathEquals: '/epz/order/extendedsearch/search.html'
		},
		{
			hostSuffix: 'zakupki.gov.ru',
			pathEquals: '/epz/order/quicksearch/search.html'
		},
		{
			hostSuffix: 'zakupki.gov.ru',
			pathEquals: '/epz/order/quicksearch/update.html'
		}
	]
});
