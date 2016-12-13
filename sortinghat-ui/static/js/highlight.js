angular.module('highlight', []).filter('highlight_vars', function (){
  return function (text){
    if (text !== undefined){
        text.replace(new RegExp("%s",g), "<code>{variable}</code>");
    }
  };
});
