/// <reference types="vite/client" />

// File System Access API — Chrome/Edge でのみサポート
interface Window {
  showDirectoryPicker(options?: { mode?: 'read' | 'readwrite' }): Promise<FileSystemDirectoryHandle>
}
