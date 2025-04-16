declare module '@rsbuild/core' {
  export function defineConfig(config: any): any;
}

declare module '@rsbuild/plugin-react' {
  export function pluginReact(): any;
}

declare module '@rsbuild/plugin-typescript' {
  export function pluginTypeScript(): any;
}

declare module '@rsbuild/plugin-electron' {
  export function pluginElectron(options: any): any;
}
