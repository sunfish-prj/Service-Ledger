'use strict';

var utils = require('../utils/writer.js');
var Invoke = require('../service/InvokeService');

//var res = require('response');

module.exports.invokePOST = function invokePOST (req, res, next) {
  var invokeSpec = req.swagger.params['invokeSpec'].value;
  
  Invoke.InvokePOST(invokeSpec)
    .then(function (response) {
      utils.writeJson(res, response, 200);
    })
    .catch(function (response) {
      utils.writeJson(res, response, 400);
    });
};
