import path from 'path';
import HtmlWebpackPlugin from 'html-webpack-plugin';
import MiniCssExtractPlugin from 'mini-css-extract-plugin';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Define environment variables
const APP_MODE = process.env.APP_MODE || 'development';
const APP_SLUG = process.env.APP_SLUG || 'kled';
const GITHUB_CLIENT_ID = process.env.GITHUB_CLIENT_ID || '';
const FEATURE_FLAGS = process.env.FEATURE_FLAGS ? JSON.parse(process.env.FEATURE_FLAGS) : {};

/** @type {import('webpack').Configuration} */
const config = {
  mode: 'development',
  entry: './src/basic-entry.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'static/js/[name].[contenthash:8].js',
    publicPath: '/',
    clean: true,
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx'],
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '#': path.resolve(__dirname, 'public'),
    },
  },
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: 'ts-loader',
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              postcssOptions: {
                plugins: [
                  'tailwindcss',
                  'autoprefixer',
                ],
              },
            },
          },
        ],
      },
      {
        test: /\.(png|jpg|jpeg|gif)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'static/images/[name].[hash:8][ext]',
        },
      },
      {
        test: /\.svg$/i,
        oneOf: [
          {
            test: /\.svg\?react$/,
            issuer: /\.[jt]sx?$/,
            use: ['@svgr/webpack'],
          },
          {
            type: 'asset/resource',
            generator: {
              filename: 'static/images/[name].[hash:8][ext]',
            },
          },
        ],
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
        generator: {
          filename: 'static/fonts/[name].[hash:8][ext]',
        },
      },
    ],
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: './public/index.html',
      title: 'Kled.io',
      templateParameters: {
        APP_MODE,
        APP_SLUG,
        GITHUB_CLIENT_ID,
        FEATURE_FLAGS: JSON.stringify(FEATURE_FLAGS),
      },
    }),
    new MiniCssExtractPlugin({
      filename: 'static/css/[name].[contenthash:8].css',
    }),
  ],
  devServer: {
    static: {
      directory: path.resolve(__dirname, 'public'),
    },
    historyApiFallback: true,
    port: 3000,
    hot: true,
  },
};

export default config;
