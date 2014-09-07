var DEBUG = true;
var SERVER_URL = 'http://192.168.1.102:8888/';
var BLACK_LIST = ['localhost', 'google', 'amazon'];

//Check to see if URL is blacklisted.
//This function should be improved to have a dynamically saved blacklist. Possibly pulled from the server.
function shouldCaptureURL(url) {
	for (var index in BLACK_LIST) {
		if (url.match(BLACK_LIST[index]) !== null) {
			return false;
		}
	}
	return true;
}

function debugLog(message) {
	if (DEBUG === true) {
		console.log(message)
	}
}

chrome.tabs.onUpdated.addListener( function (tabId, changeInfo, tab) {

	if (changeInfo.status == 'complete') {
		if (!shouldCaptureURL(tab.url)) {
		//if (tab.url.indexOf('localhost') > -1 || tab.url.indexOf('.google.') > -1) {
			debugLog("Black listed URL encountered. Not capturing: " + tab.url)
			return;
		}
		chrome.tabs.sendRequest(tab.id, {method: 'getText'}, function(response) {
			if(response.method=="getText"){
				var data = {
					"Content": response.data,
					"Title":tab.title,
					"Url":tab.url
				};

				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function(){
					if (xhr.readyState == 4 && ( xhr.status == 200 || xhr.status == 201)) {
						console.log('Data sent: ' + xhr.responseText);
					} else if (xhr.readyState == 4) {
						console.log('Error: ' + xhr.responseText);
					}
				};

				xhr.open('POST', SERVER_URL+ "pages", true);
				xhr.setRequestHeader("Content-type", "application/json");
				xhr.send(JSON.stringify(data));
			}
		});
	}
});
