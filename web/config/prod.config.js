const constant = require('./constant')
// let constant = JSON.parse(fs.readFileSync(path.resolve("/etc/config/config.json"), 'utf8'));
console.log(constant);

module.exports = {
  '/atomci': {
    target: constant.atomci,
    changeOrigin: true,
    secure: false,
  },
}
