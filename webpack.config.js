const path = require('path');

module.exports = {
  entry: './src/index.js',
  output: {
    filename: 'scripts.js',
    path: path.resolve(__dirname, 'static'),
  },
};