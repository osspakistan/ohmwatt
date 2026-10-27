import { globalIgnores } from 'eslint/config'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import pluginVue from 'eslint-plugin-vue'
import pluginOxlint from 'eslint-plugin-oxlint'
import skipFormatting from 'eslint-config-prettier/flat'
import pluginCheckFile from 'eslint-plugin-check-file'
import pluginImportX from 'eslint-plugin-import-x'

export default defineConfigWithVueTs(
  {
    name: 'app/files-to-lint',
    files: ['**/*.{vue,ts,mts,tsx}'],
  },

  globalIgnores([
    '**/dist/**',
    '**/dist-ssr/**',
    '**/coverage/**',
    '**/components/ui/**',
    '**/src-tauri/**',
    '**/vite-env.d.ts',
    'vite.config.ts',
    'tsconfig.json',
  ]),

  ...pluginVue.configs['flat/essential'],
  vueTsConfigs.recommended,

  ...pluginOxlint.buildFromOxlintConfigFile('.oxlintrc.json'),

  // Vue template relaxations
  {
    name: 'app/vue-rules',
    rules: {
      'vue/html-self-closing': 'off',
      'vue/no-multiple-template-root': 'off',
    },
  },

  // Enforce KEBAB_CASE naming convention for files and folders
  {
    name: 'app/naming-conventions',
    files: ['src/**/*.{ts,tsx,vue}'],
    plugins: {
      'check-file': pluginCheckFile,
    },
    rules: {
      'check-file/filename-naming-convention': [
        'error',
        {
          '**/*.{ts,tsx,vue}': 'KEBAB_CASE',
        },
        {
          ignoreMiddleExtensions: true,
        },
      ],
      'check-file/folder-naming-convention': [
        'error',
        {
          '**/': 'KEBAB_CASE',
        },
      ],
    },
  },

  // Enforce import order
  {
    name: 'app/import-order',
    files: ['src/**/*.{ts,tsx,vue}'],
    plugins: {
      'import-x': pluginImportX,
    },
    rules: {
      'import-x/order': [
        'error',
        {
          groups: [
            'builtin',
            'external',
            'internal',
            'parent',
            'sibling',
            'index',
          ],
          'newlines-between': 'always',
          alphabetize: { order: 'asc', caseInsensitive: true },
          warnOnUnassignedImports: false,
        },
      ],
    },
  },

  skipFormatting,
)
