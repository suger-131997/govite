import { spawn, ChildProcess } from 'child_process'
import type { Plugin } from 'vite'

export interface GoDevRunnerOptions {
  /** Go entry file to run. Defaults to `'./entrypoint/dev/main.go'` */
  entry?: string
}

export default function goDevRunner(options?: GoDevRunnerOptions): Plugin {
  const entry = options?.entry ?? './entrypoint/dev/main.go'

  return {
    name: 'vite-plugin-go-dev-runner',
    configureServer(server) {
      const startGoServer = () =>
        spawn('go', ['run', entry], {
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
        process.off('SIGINT', serverClose)
        process.off('SIGTERM', serverClose)
      })
    },
  }
}
