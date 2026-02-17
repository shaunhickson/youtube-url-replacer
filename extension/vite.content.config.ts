import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
  build: {
    emptyOutDir: false,
    rollupOptions: {
      input: {
        content: resolve(__dirname, 'src/content.ts'),
      },
      output: {
        entryFileNames: 'assets/[name].js',
        format: 'iife',
        inlineDynamicImports: true
      }
    }
  }
})
