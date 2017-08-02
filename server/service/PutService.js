'use strict';


/**
 * Storing a key-value pair 
 *
 * putSpec Put-request-body Body in JSON
 * returns response
 **/
exports.putPOST = function(putSpec) {
  return new Promise(function(resolve, reject) {
    var examples = {};
    examples['application/json'] = {
  "message" : "aeiou-PUT"
};
    if (Object.keys(examples).length > 0) {
      resolve(examples[Object.keys(examples)[0]]);
    } else {
      resolve();
    }
  });
}

