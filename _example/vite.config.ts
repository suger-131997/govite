import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import entriesConfig from './entries.gen.json'
import { spawn ,ChildProcess} from 'child_process'

export default defineConfig({
  plugins: [
    react(),
    {
      name: 'run-go-server',
      configureServer(server) {
        const startGoServer = () => spawn('go', ['run', './entrypoint/dev/main.go'], {
          stdio: 'inherit',
          detached: true,
        })

        let goProcess = startGoServer()
        let isRestarting = false

        const stopGoServer = (proc: ChildProcess) => {
          if (proc.pid) {
            try {
              process.kill(-proc.pid, 'SIGTERM')
            } catch {
              proc.kill()
            }
          } else {
            proc.kill()
          }
        }

        // server.watcher.add(projectRoot)
        server.watcher.on('change', (file) => {
          if (file.endsWith('.go') && !isRestarting) {
            isRestarting = true

            const startNew = () => {
              goProcess = startGoServer()
              isRestarting = false
              server.ws.send({ type: 'full-reload' })
            }

            if (goProcess.exitCode !== null) {
              startNew()
            } else {
              goProcess.once('exit', startNew)
              stopGoServer(goProcess)
            }
          }
        })

      const serverClose = () => {
        server.close()
      }
        process.on('SIGINT', serverClose)
        process.on('SIGTERM', serverClose)

        server.httpServer?.on('close', () => {
          stopGoServer(goProcess)
          process.off("SIGINT", serverClose);
          process.off("SIGTERM", serverClose);
        })

      },
    },
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
