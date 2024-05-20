const path = require('path');

module.exports = {
  entry: './src/index.ts',
  output: {
    filename: 'scripts.js',
    path: path.resolve(__dirname, 'static')
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js']
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      },
      {
        enforce: 'pre',
        test: /\.js$/,
        loader: 'source-map-loader'
      }
    ]
  },
  devtool: 'source-map', // This helps with debugging by providing source maps
};
