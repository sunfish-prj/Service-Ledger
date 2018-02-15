'use strict';

var url = require('url');

var Exec = require('./ExecService');

module.exports.invokePOST = function invokePOST (req, res, next) {
  Exec.invokePOST(req.swagger.params, res, next);
};
