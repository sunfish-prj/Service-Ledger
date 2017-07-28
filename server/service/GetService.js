'use strict';


/**
 * Retrieving a value by its key 
 *
 * getId Get-request-body Body in JSON
 * returns response
 **/
exports.getPOST = function(getId) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = {
  "message" : "aeiou"
};
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}

