import { useState } from 'react'

const FilesList = ({ files, onRefresh }) => {
  const [filter, setFilter] = useState('all')
  const [sortBy, setSortBy] = useState('modTime')

  const formatFileSize = (bytes) => {
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString()
  }

  const getFileIcon = (type) => {
    switch (type) {
      case 'video':
        return '🎥'
      case 'audio':
        return '🎵'
      default:
        return '📄'
    }
  }

  const handleDownload = (downloadUrl, filename) => {
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const filteredFiles = files.filter(file => {
    if (filter === 'all') return true
    return file.type === filter
  })

  const sortedFiles = [...filteredFiles].sort((a, b) => {
    switch (sortBy) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'size':
        return b.size - a.size
      case 'modTime':
      default:
        return new Date(b.modTime) - new Date(a.modTime)
    }
  })

  if (files.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem', color: '#6c757d' }}>
        <h3>No files available</h3>
        <p>Download some videos to see them here!</p>
        <button className="btn btn-primary" onClick={onRefresh} style={{ marginTop: '1rem' }}>
          Refresh
        </button>
      </div>
    )
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h3>Downloaded Files ({files.length})</h3>
        <button className="btn btn-primary" onClick={onRefresh}>
          Refresh
        </button>
      </div>

      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1.5rem', flexWrap: 'wrap' }}>
        <div>
          <label htmlFor="filter" style={{ marginRight: '0.5rem', fontSize: '0.9rem' }}>Filter:</label>
          <select
            id="filter"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            className="form-control"
            style={{ width: 'auto', display: 'inline-block' }}
          >
            <option value="all">All Files</option>
            <option value="video">Videos</option>
            <option value="audio">Audio</option>
          </select>
        </div>

        <div>
          <label htmlFor="sort" style={{ marginRight: '0.5rem', fontSize: '0.9rem' }}>Sort by:</label>
          <select
            id="sort"
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value)}
            className="form-control"
            style={{ width: 'auto', display: 'inline-block' }}
          >
            <option value="modTime">Date Modified</option>
            <option value="name">Name</option>
            <option value="size">Size</option>
          </select>
        </div>
      </div>

      <div className="files-grid">
        {sortedFiles.map((file) => (
          <div key={file.name} className="file-card">
            <div className="file-info">
              <h4>
                {getFileIcon(file.type)} {file.name}
              </h4>
              <div className="file-meta">
                <p><strong>Size:</strong> {formatFileSize(file.size)}</p>
                <p><strong>Type:</strong> {file.type}</p>
                <p><strong>Modified:</strong> {formatDate(file.modTime)}</p>
              </div>
            </div>
            <button
              className="btn btn-download"
              onClick={() => handleDownload(file.downloadUrl, file.name)}
            >
              Download
            </button>
          </div>
        ))}
      </div>

      {filteredFiles.length === 0 && filter !== 'all' && (
        <div style={{ textAlign: 'center', padding: '2rem', color: '#6c757d' }}>
          <p>No {filter} files found.</p>
        </div>
      )}
    </div>
  )
}

export default FilesList
