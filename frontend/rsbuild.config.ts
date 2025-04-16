import { defineConfig } from '@rsbuild/core';

export default defineConfig({
  tools: {
    webpack: (config: any) => {
      config.module.rules.push({
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: 'ts-loader',
      });
      return config;
    },
    tailwindcss: true,
  },

  source: {
    entry: {
      index: './src/entry.client.tsx',
    },
  },
  dev: {
    port: 3000,
    writeToDisk: false,
  },
  html: {
    template: './public/index.html',
    title: 'Agent Runtime',
  },
  output: {
    distPath: {
      root: 'dist',
      js: 'static/js',
      css: 'static/css',
      html: '',
      image: 'static/images',
      media: 'static/media',
      font: 'static/fonts',
    },
    filename: {
      js: '[name].[contenthash:8].js',
      css: '[name].[contenthash:8].css',
    },
    publicPath: '/',
  },
  server: {
    publicDir: 'public',
  },
});
