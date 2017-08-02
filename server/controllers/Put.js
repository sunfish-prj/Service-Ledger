'use strict';

var utils = require('../utils/writer.js');
var Put = require('../service/PutService');

module.exports.putPOST = function putPOST (req, res, next) {
  var putSpec = req.swagger.params['putSpec'].value;
  
  Put.putPOST(putSpec)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
