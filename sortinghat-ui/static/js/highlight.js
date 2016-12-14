angular.module('highlight', []).filter('highlight_vars', function (){
  return function (text){
      var data = text
    if (text != undefined){
        data = text.replace(new RegExp("%s",'g'), "<code>%s</code>");
    }
    return data
  };
});
