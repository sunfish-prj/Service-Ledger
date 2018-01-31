'use strict';

var utils = require('../../hyperledger/fabric/utils/writer.js');
var Get = require('../service/GetService');

module.exports.getPOST = function getPOST (req, res, next) {
  var getId = req.swagger.params['getId'].value;
  Get.getPOST(getId)
    .then(function (response) {
      utils.writeJson(res, response, 200);
    })
    .catch(function (response) {
      utils.writeJson(res, response, 404);
    });
};
