
angular.module('timeFilters', []).filter('timeFilter', function (){

  return function (items, last_period) {
    var timestamp = Date.now();
    var oldest_timestamp = 0;
    if (last_period) {
       oldest_timestamp = timestamp - (last_period * 60 * 1000);
    }

    var result = [];
    if (typeof items != 'undefined') {
        for (var i=0; i<items.length; i++) {
            if ((items[i].timestamp * 1000) > oldest_timestamp)  {
                result.push(items[i]);
            }
        }
    }
    return result;

  };
});
