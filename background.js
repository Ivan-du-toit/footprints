//Background
chrome.tabs.onUpdated.addListener( function (tabId, changeInfo, tab) {
	
	if (changeInfo.status == 'complete') {
		console.log('Opened: ' + tab.url);
		var url = tab.url;
		chrome.tabs.sendRequest(tab.id, {method: 'getText'}, function(response) {
			if(response.method=="getText"){
				console.log('Text: ' + response.data);
				
				var data = 'URL=' +  encodeURIComponent(url) + '&data=' + encodeURIComponent(response.data);
				
				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function(){
					if(xhr.readyState == 4 && xhr.status == 200) {
						console.log('Data sent: ' + xhr.responseText);
					} else if (xhr.readyState == 4) {
						console.log('Error: ' + xhr.responseText);
					}
				};
				
				xhr.open('POST', "http://localhost/test/dump.php", true);
				xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
				xhr.send(data);
			}
		});
	}
});
