// This is whistle

var app = angular.module('whistleApp', ['igTruncate','ui.bootstrap', 'ngRoute', 'timeFilters']);
var base_path = "/whistle-ui"
app.config(function($routeProvider){
    $routeProvider
         // for the host page
    .when("/host/:hostId", {
        templateUrl: "host.html",
        controller: "HostController"
    })
    .when("/category/:categoryName", {
        templateUrl: "main.html",
        controller: "MainController"
    })
    .when("/settings/",{
        templateUrl: "settings.html",
        controller: "SettingsController"
    })
    .when("/bucket/:bucketId", {
         templateUrl: "bucket.html",
         controller: "BucketController"
     })
    .when("/alerts",{
        templateUrl: "alerts.html",
        controller: "AlertController"
    })
    .otherwise({
        templateUrl: "main.html",
        controller: "MainController"
    })
});


function add_hostnames(record) {
    record.hostnames = new Array();
    for (var j=0; j < record.categories.length; j++) {
        var category_name = record.categories[j];
        var category_info = record.category_info[category_name];
        record.hostnames = record.hostnames.concat(category_info.hostnames);
    }
    return record;
}

// Now the MainController
app.controller('MainController', function($scope, $http, $interval, $location, $routeParams) {
   $scope.isCollapsed = true;
   $scope.content = {ts:'',message:'click on host'}

    $scope.go_host = function ( path  ) {
           $location.path( "/host/"+path  );
    };

    $scope.category = 'all';
    if ($routeParams.categoryName) {
         $scope.category = $routeParams.categoryName;
    }


    $scope.get_category_data = function(category){
        url = base_path+"/category/"+category;
        if (category == "all") {
             url = base_path+"/alerts";
        }
        // This is ugly and needs to be fixed in the backend
        // why are there two APIs for all vs categories
        $http.get(url).success(function(response) {
            response.forEach(function(element){
                if (element.user_message == 'Undefined' ) {
                    element.user_message = element.canonical_message;
                }
            });
            $scope.cat_records = response;
            window.scrollTo(0, 0);
        });
    };

    $scope.get_category_data($scope.category);

    $scope.get_bucket_details = function(bucket_id) {
        $location.path("/bucket/"+JSON.stringify(bucket_id));
    };
});

app.controller('HostController', function($scope, $http, $interval, $routeParams) {
    $scope.get_host_data = function(host){
                                $http.get(base_path+"/host/"+host).success(function(response) {
                                                $scope.content = response;
                                                window.scrollTo(0, 0);
                                });
                            };
    $scope.get_host_data($routeParams.hostId);
    $scope.host = $routeParams.hostId;
});

app.controller('NavController', function($scope, $http, $interval, $routeParams) {
    $scope.navItems = [ {'name':'all', 'class':'inactive'},
                        {'name':'prod', 'class':'active'},
                        {'name':'vmware', 'class':'inactive'},
                        {'name':'dogfood', 'class':'inactive'},
                        {'name':'test', 'class':'inactive'}];
    var clockTick = function() {
         var d = new Date();
         $scope.timeUTC = d.toUTCString();
    };
    console.log($routeParams.categoryName);
    if ($routeParams.categoryName) {
        for (item in $scope.navItems) {
            if (item['name'] == $routeParams.categoryName) {
                console.log('Category '+ item['name']);
                item['class'] = 'active';
            }else {
                console.log('Inactive');
                item['class'] = 'inactive';
                }
        }
    }
    $interval(clockTick, 1000);
});

app.controller('BucketController', function($scope, $http, $interval, $routeParams, $modal) {

    $scope.showFullMessage = "false";
    $scope.get_bucket_data = function(bucketId) {
        $http.get(base_path+"/bucket/details/"+bucketId).success(function(response){
            $scope.bucket_data = response;
            var record = $scope.bucket_data.bucket;
         });
    };

    $scope.mute_bucket = function(bucket_id, msg) {
        $http.post(base_path+'/mute/'+bucket_id+"/"+msg).then(function(response) {
            $scope.mute_success = true;
        }, function(response) {
            $scope.mute_failure = true;
        });
    };

    $scope.save_user_message = function() {
        $http.post(base_path+'/bucket/'+$scope.bucket_id+"/"+$scope.bucket_data.bucket[0].user_message).then(function(response) {
            $scope.mute_success = true;
        }, function(response) {
            $scope.mute_failure = true;
        });
    }


    $scope.account_mute_bucket = function(bucket_id, accounts, msg) {
        data = {'accounts':accounts, 'msg': msg};
        $http.post(base_path+'/mute-account-bucket/'+bucket_id, data).then(function(response) {
            $scope.mute_success = true;
        }, function(response) {
            $scope.mute_failure = true;
        });
    };



    $scope.mute_dialog_open = function(bucketId, hostnames) {
        var modalInstance = $modal.open({
              animation: $scope.animationsEnabled,
              templateUrl: 'MuteModalContent.html',
              controller: 'MuteController',
              size: 'lg',
              resolve: {
                accounts: function () {
                  return hostnames;
                },
                bucket_id: function () {
                  return bucketId;
                }
              }
        });

        modalInstance.result.then(function (retobj) {
            if (retobj.op == 'mute_accounts') {
                $scope.account_mute_bucket(retobj.bucket, retobj.accounts, retobj.msg);
            } else if(retobj.op == 'mute_bucket'){
                $scope.mute_bucket(retobj.bucket, retobj.msg);
            }

        });
    };

    $scope.bucketId=$routeParams.bucketId;
    $scope.get_bucket_data($routeParams.bucketId);
});

app.controller('SettingsController', function($scope, $http) {
	$scope.reset_alert = {type:'hidden'};
    $scope.mute_buckets  = [];
    $scope.reset_buckets = function() {
        $http.post(base_path+'/alerts').then(function(response) {
            $scope.reset_success = true;
        }, function(response) {
            $scope.reset_failure = true;
        });
    };
    
	$scope.get_mute_buckets = function() {
	    $http.get(base_path+'/settings/mute').success(function(response){
	        mute_settings = response;
	        for (var i = 0; i < mute_settings.length; i++) {
	            var url = base_path+'/bucket/summary/'+ JSON.stringify(mute_settings[i].bucket_id);
	            var msg = mute_settings[i].msg;
	            temp_func = function(msg) {
	                return function(response){
                        response.mute_msg = msg;
                        //ßresponse = add_hostnames(response);
                        $scope.mute_buckets.push(response);
	                };
	            }
                success_cb = temp_func(msg)
	            $http.get(url).success(success_cb);
	        }
	    });
	};

	$scope.unmute_bucket = function(bucket_id) {
        $http.post(base_path+'/unmute/'+JSON.stringify(bucket_id)).then(function(response) {
            $scope.umute_success = true;
        }, function(response) {
            $scope.umute_failure = true;
        });
    };
    $scope.mute_bucket_alert = {type:'hidden'};
	$scope.get_mute_buckets();
});

app.controller('MuteController', function($scope, $http, $modalInstance, accounts, bucket_id) {

    $scope.accounts = accounts;
    $scope.bucket_id = bucket_id;
    $scope.selected_accounts = accounts.slice();
    $scope.msg = "";

    // toggle selection for a given account name
    $scope.toggleSelection = function toggleSelection(account_name) {
        var idx = $scope.selected_accounts.indexOf(account_name);

        // is currently selected
        if (idx > -1) {
          $scope.selected_accounts.splice(idx, 1);
        }

        // is newly selected
        else {
          $scope.selected_accounts.push(account_name);
        }
     };

    $scope.mute_bucket = function() {
        $modalInstance.close({'op':'mute_bucket', 'bucket':$scope.bucket_id, 'msg':$scope.msg});
    };

    $scope.mute = function() {
        $modalInstance.close({'op':'mute_accounts', 'bucket':$scope.bucket_id, 'accounts':$scope.selected_accounts, 'msg':$scope.msg});
    };

    $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
    };

});
