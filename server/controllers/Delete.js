'use strict';

var utils = require('../utils/writer.js');
var Delete = require('../service/DeleteService');

module.exports.deletePOST = function deletePOST (req, res, next) {
  var body = req.swagger.params['body'].value;
  Delete.deletePOST(body)
    .then(function (response) {
      utils.writeJson(res, response);
    })
    .catch(function (response) {
      utils.writeJson(res, response);
    });
};
