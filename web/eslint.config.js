import vue from 'eslint-plugin-vue'
import typescript from '@vue/eslint-config-typescript'
import { includeIgnoreFile } from '@eslint/compat'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const gitignorePath = path.resolve(__dirname, '.gitignore')

export default [
  includeIgnoreFile(gitignorePath),
  {
    ignores: ['**/dist/**', '**/node_modules/**']
  },
  ...vue.configs['flat/recommended'],
  ...typescript(),
  {
    files: ['**/*.{vue,js,jsx,cjs,mjs,ts,tsx,cts,mts}'],
    rules: {
      // Prettier rules are handled by running prettier separately
      // Allow single-word component names for pages
      'vue/multi-word-component-names': 'off',
      // Allow any type (can be tightened later)
      '@typescript-eslint/no-explicit-any': 'warn',
      // Allow unused vars with underscore prefix
      '@typescript-eslint/no-unused-vars': ['error', {
        argsIgnorePattern: '^_',
        varsIgnorePattern: '^_'
      }]
    }
  },
  {
    files: ['**/*.config.{js,cjs,mjs}', '**/vite.config.*', '**/tailwind.config.*'],
    rules: {
      // Allow require in config files
      '@typescript-eslint/no-require-imports': 'off'
    }
  }
]
