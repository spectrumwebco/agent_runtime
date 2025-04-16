#!/usr/bin/env node

import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';

const colors = {
  reset: '\x1b[0m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
};

console.log(`${colors.cyan}Setting up Shadcn UI components...${colors.reset}`);

const components = [
  'accordion',
  'alert',
  'alert-dialog',
  'aspect-ratio',
  'avatar',
  'badge',
  'button',
  'calendar',
  'card',
  'checkbox',
  'collapsible',
  'command',
  'context-menu',
  'dialog',
  'dropdown-menu',
  'form',
  'hover-card',
  'input',
  'label',
  'menubar',
  'navigation-menu',
  'popover',
  'progress',
  'radio-group',
  'scroll-area',
  'select',
  'separator',
  'sheet',
  'skeleton',
  'slider',
  'switch',
  'table',
  'tabs',
  'textarea',
  'toast',
  'toggle',
  'tooltip',
];

const componentsJsonPath = path.join(process.cwd(), 'components.json');
if (!fs.existsSync(componentsJsonPath)) {
  console.log(`${colors.yellow}Creating components.json configuration...${colors.reset}`);
  
  const componentsJson = {
    "$schema": "https://ui.shadcn.com/schema.json",
    "style": "default",
    "rsc": false,
    "tsx": true,
    "tailwind": {
      "config": "tailwind.config.ts", // Updated to .ts extension
      "css": "src/styles/globals.css",
      "baseColor": "slate",
      "cssVariables": true
    },
    "aliases": {
      "components": "@/components",
      "utils": "@/lib/utils"
    }
  };
  
  fs.writeFileSync(componentsJsonPath, JSON.stringify(componentsJson, null, 2));
  console.log(`${colors.green}Created components.json${colors.reset}`);
}

const utilsDir = path.join(process.cwd(), 'src', 'lib');
const utilsPath = path.join(utilsDir, 'utils.ts');

if (!fs.existsSync(utilsDir)) {
  fs.mkdirSync(utilsDir, { recursive: true });
}

if (!fs.existsSync(utilsPath)) {
  console.log(`${colors.yellow}Creating utils.ts...${colors.reset}`);
  
  const utilsContent = `import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`;
  
  fs.writeFileSync(utilsPath, utilsContent);
  console.log(`${colors.green}Created utils.ts${colors.reset}`);
}

console.log(`${colors.blue}Installing Shadcn UI components...${colors.reset}`);

const componentsDir = path.join(process.cwd(), 'src', 'components');
const uiComponentsDir = path.join(componentsDir, 'ui');

if (!fs.existsSync(componentsDir)) {
  fs.mkdirSync(componentsDir, { recursive: true });
}

if (!fs.existsSync(uiComponentsDir)) {
  fs.mkdirSync(uiComponentsDir, { recursive: true });
}

components.forEach((component) => {
  try {
    console.log(`${colors.magenta}Installing ${component}...${colors.reset}`);
    execSync(`npx shadcn-ui@latest add ${component} --yes`, { stdio: 'inherit' });
    console.log(`${colors.green}Successfully installed ${component}${colors.reset}`);
  } catch (error) {
    console.error(`Error installing ${component}:`, (error as Error).message);
  }
});

console.log(`${colors.green}Shadcn UI setup complete!${colors.reset}`);
