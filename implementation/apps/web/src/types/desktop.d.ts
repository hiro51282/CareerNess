// Electron preload（apps/desktop/src/preload.ts）が公開するブリッジの型。
// ブラウザ実行時は undefined（存在チェックで desktop モードを判定する）。
export {}

declare global {
  interface Window {
    careerness?: {
      desktop: boolean
      pickDirectory(): Promise<string | null>
    }
  }
}
