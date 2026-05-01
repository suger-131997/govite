import type { Plugin } from 'vite';
export interface GoDevRunnerOptions {
    /** Go entry file to run. Defaults to `'./entrypoint/dev/main.go'` */
    entry?: string;
}
export default function goDevRunner(options?: GoDevRunnerOptions): Plugin;
//# sourceMappingURL=index.d.ts.map