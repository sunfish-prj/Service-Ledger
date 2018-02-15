'use strict';

var url = require('url');

var Store = require('./StoreService');

module.exports.deletePOST = function deletePOST (req, res, next) {
  Store.deletePOST(req.swagger.params, res, next);
};

module.exports.getKeysPOST = function getKeysPOST (req, res, next) {
  Store.getKeysPOST(req.swagger.params, res, next);
};

module.exports.getPOST = function getPOST (req, res, next) {
  Store.getPOST(req.swagger.params, res, next);
};

module.exports.putPOST = function putPOST (req, res, next) {
  Store.putPOST(req.swagger.params, res, next);
};
