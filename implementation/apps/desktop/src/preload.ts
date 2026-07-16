// preload — renderer へ公開する最小ブリッジ（contextIsolation 前提）。
// desktop フラグと、attach 用のネイティブフォルダ選択のみを公開する。
import { contextBridge, ipcRenderer } from 'electron'

contextBridge.exposeInMainWorld('careerness', {
  desktop: true,
  pickDirectory: (): Promise<string | null> => ipcRenderer.invoke('careerness:pickDirectory'),
})
