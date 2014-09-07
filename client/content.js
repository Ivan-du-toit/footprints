chrome.extension.onRequest.addListener(
    function(request, sender, callback) {
        if(request.method == "getText"){
            callback({data: document.all[0].innerText, method: "getText"}); //same as innerText
        }
    }
);

// Warn if the browser doesn't support addEventListener or the Page Visibility API
if (typeof document.addEventListener === "undefined" || typeof hidden === "undefined") {
	console.log('Replace  this to return the content if it does not work!!!')
//	alert("This demo requires a browser, such as Google Chrome or Firefox, that supports the Page Visibility API.");
} else {
	// Handle page visibility change   
	console.log('The page is being viewed. Send the content');
	//document.addEventListener(visibilityChange, handleVisibilityChange, false);
}