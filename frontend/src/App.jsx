import { useState, useEffect } from 'react'
import DownloadForm from './components/DownloadForm'
import DownloadHistory from './components/DownloadHistory'
import FilesList from './components/FilesList'

function App() {
  const [downloads, setDownloads] = useState([])
  const [files, setFiles] = useState([])
  const [activeTab, setActiveTab] = useState('download')

  const addDownload = (download) => {
    setDownloads(prev => [download, ...prev])
  }

  const updateDownload = (id, updates) => {
    setDownloads(prev => prev.map(d => d.id === id ? { ...d, ...updates } : d))
  }

  const fetchFiles = async () => {
    try {
      const response = await fetch('/api/files')
      const data = await response.json()
      setFiles(data.files || [])
    } catch (error) {
      console.error('Failed to fetch files:', error)
    }
  }

  useEffect(() => {
    fetchFiles()
  }, [])

  return (
    <div className="container">
      <div className="header">
        <h1>🎥 YouTube Downloader</h1>
        <p>Download videos and audio in multiple formats and qualities</p>
      </div>

      <div className="card">
        <div style={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
          <button
            className={`btn ${activeTab === 'download' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setActiveTab('download')}
          >
            Download
          </button>
          <button
            className={`btn ${activeTab === 'history' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setActiveTab('history')}
          >
            History ({downloads.length})
          </button>
          <button
            className={`btn ${activeTab === 'files' ? 'btn-primary' : 'btn-secondary'}`}
            onClick={() => setActiveTab('files')}
          >
            Files ({files.length})
          </button>
        </div>

        {activeTab === 'download' && (
          <DownloadForm
            onDownloadStart={addDownload}
            onDownloadUpdate={updateDownload}
            onDownloadComplete={fetchFiles}
          />
        )}

        {activeTab === 'history' && (
          <DownloadHistory downloads={downloads} />
        )}

        {activeTab === 'files' && (
          <FilesList files={files} onRefresh={fetchFiles} />
        )}
      </div>
    </div>
  )
}

export default App
