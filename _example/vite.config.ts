import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import entriesConfig from './entries.gen.json'
import goDevRunner from 'vite-plugin-go-dev-runner'

export default defineConfig({
  plugins: [
    react(),
    goDevRunner(),
  ],
  resolve: {
    alias: {
      '~': path.resolve(__dirname, "./"),
    },
  },
  build: {
    outDir: path.resolve(__dirname, "./entrypoint/prod/dist"),
    manifest: true,
    rolldownOptions: {
      input: entriesConfig,
    },
  },
  server: {
    host: true,
    watch: {
      ignored: ['**/entries.gen.json'],
    },
  },
})
