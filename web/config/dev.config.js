const constant = require('./constant')
console.log(constant);

module.exports = {
  '/atomci': {
    target: constant.atomci,
    changeOrigin: true,
    secure: false,
    ws: true,
  }
}
