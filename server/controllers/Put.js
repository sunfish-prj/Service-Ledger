'use strict';

var utils = require('../../hyperledger/fabric/utils/writer.js');
var Put = require('../service/PutService');

//var res = require('response');

module.exports.putPOST = function putPOST (req, res, next) {
  var putSpec = req.swagger.params['putSpec'].value;
  
  Put.putPOST(putSpec)
    .then(function (response) {
      utils.writeJson(res, response, 200);
    })
    .catch(function (response) {
      utils.writeJson(res, response, 400);
    });
};
