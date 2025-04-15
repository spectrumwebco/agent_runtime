import { defineConfig } from '@rsbuild/core';
import { pluginReact } from '@rsbuild/plugin-react';
import { pluginTypeScript } from '@rsbuild/plugin-typescript';
import { pluginElectron } from '@rsbuild/plugin-electron';

export default defineConfig({
  plugins: [
    pluginReact(),
    pluginTypeScript(),
    pluginElectron({
      main: {
        entry: {
          main: './src/electron/main.ts',
        },
      },
      preload: {
        entry: {
          preload: './src/electron/preload.ts',
        },
      },
    }),
  ],
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
  tools: {
    tailwindcss: true,
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
