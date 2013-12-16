chrome.tabs.onUpdated.addListener( function (tabId, changeInfo, tab) {
	
	if (changeInfo.status == 'complete') {
		if (tab.url.indexOf('localhost') > -1 || tab.url.indexOf('.google.') > -1) {
			console.log('Opened: ' + tab.url);
			console.log('not logging');
			return;
		}
		var url = tab.url;
		chrome.tabs.sendRequest(tab.id, {method: 'getText'}, function(response) {
		
			if(response.method=="getText"){
				var data = {
					"add": {
						"doc":{
							"id": url,
							"url": url,
							"content": response.data
						},
						"boost":1.0,
						"overwrite":true,
						"commitWithin":1000
					}
				};
				
				
				
				var xhr = new XMLHttpRequest();
				xhr.onreadystatechange = function(){
					if(xhr.readyState == 4 && xhr.status == 200) {
						console.log('Data sent: ' + xhr.responseText);
					} else if (xhr.readyState == 4) {
						console.log('Error: ' + xhr.responseText);
					}
				};
				
				xhr.open('POST', "http://localhost:8983/solr/collection1/update?wt=json", true);
				xhr.setRequestHeader("Content-type", "application/json");
				xhr.send(JSON.stringify(data));
			}
		});
	}
});