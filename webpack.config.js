const path = require('path');

module.exports = {
    target: 'es5',
    mode: 'development',
    entry: {
      index: './client/src/index.ts',
      register: './client/src/register.ts',
      login: './client/src/login.ts',
  },
    devtool: 'inline-source-map',
  module: {
    rules: [
      {
        test: /\.ts?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    fallback: {
        "crypto": require.resolve("crypto-browserify"),
        "stream": require.resolve("stream-browserify"),
        "util": require.resolve("util-promisify/"),
        "timers": require.resolve("timers-browserify")
    },
    extensions: ['.tsx', '.ts', '.js'],
  },
  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, './client/build'),
    chunkFormat: 'module',
    clean: true,
  },
  devServer: {
    static: path.join(__dirname, "dist"),
    compress: true,
    port: 4000,
  },
};
