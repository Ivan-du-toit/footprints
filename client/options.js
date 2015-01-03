function saveOptions() {
	var server = document.getElementById('server').value;
	chrome.storage.sync.set({
		server: server
	}, function() {
		// Update status to let user know options were saved.
		var status = document.getElementById('status');
		status.textContent = 'Options saved.';
		setTimeout(function() {
			status.textContent = '';
		}, 750);
	});
}

function init() {
	document.getElementById('save').addEventListener('click', saveOptions);
	restoreOptions();
}

function restoreOptions() {
	
	chrome.storage.sync.get({
		server: 'http://192.168.1.102:8888/'
	}, function(items) {
		document.getElementById('server').value = items.server;
	});
}

document.addEventListener('DOMContentLoaded', init);
