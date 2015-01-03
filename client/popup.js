'strict';



var footprints= angular.module('footprints', []);

footprints.controller('SearchCtrl', ['$scope', '$http', function ($scope, $http) {
	$scope.searchItems = [];
	$scope.error = '';
	$scope.query = '';
	
	$scope.settings = {};

	$scope.loadSettings = function() {
		chrome.storage.sync.get({
			server: 'http://192.168.1.102:8888/'
		}, function(items) {
			$scope.settings = items;
		});
	}
	
	$scope.search = function() {
		$http.get($scope.settings.server+ "search/" + $scope.query).success(function(data) {
			$scope.searchItems = data.Pages;
			console.log(data);
			console.log($scope.searchItems);
		}).error(function(error) {
			$scope.error = error
			console.log(error);
		});
		//xhr.setRequestHeader("Content-type", "application/json");
		//xhr.send();
	}

	$scope.loadSettings();
}]);