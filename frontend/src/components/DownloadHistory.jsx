const DownloadHistory = ({ downloads }) => {
  const formatFileSize = (bytes) => {
    if (!bytes) return 'Unknown'
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  const formatDuration = (startTime, endTime) => {
    if (!startTime || !endTime) return 'Unknown'
    const duration = Math.round((new Date(endTime) - new Date(startTime)) / 1000)
    return `${duration}s`
  }

  const getStatusBadge = (status) => {
    const statusClasses = {
      downloading: 'status-downloading',
      completed: 'status-success',
      error: 'status-error'
    }
    
    const statusText = {
      downloading: 'Downloading',
      completed: 'Completed',
      error: 'Failed'
    }

    return (
      <span className={`status-badge ${statusClasses[status] || 'status-badge'}`}>
        {statusText[status] || status}
      </span>
    )
  }

  const handleDownload = (downloadUrl, filename) => {
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  if (downloads.length === 0) {
    return (
      <div style={{ textAlign: 'center', padding: '2rem', color: '#6c757d' }}>
        <h3>No downloads yet</h3>
        <p>Start downloading videos to see them here!</p>
      </div>
    )
  }

  return (
    <div>
      <h3 style={{ marginBottom: '1.5rem' }}>Download History</h3>
      {downloads.map((download) => (
        <div key={download.id} className="download-item">
          <div className="download-info">
            <h4>{download.filename || 'Processing...'}</h4>
            <p>
              <strong>URL:</strong> {download.url.length > 50 ? download.url.substring(0, 50) + '...' : download.url}
            </p>
            <p>
              <strong>Format:</strong> {download.format} 
              {download.format === 'video' && download.resolution && ` • ${download.resolution}p`}
              {download.format === 'video' && download.videoFormat && ` • ${download.videoFormat.toUpperCase()}`}
            </p>
            {download.startTime && (
              <p>
                <strong>Started:</strong> {new Date(download.startTime).toLocaleString()}
                {download.endTime && (
                  <span> • <strong>Duration:</strong> {formatDuration(download.startTime, download.endTime)}</span>
                )}
              </p>
            )}
            {download.error && (
              <p style={{ color: 'var(--error-color)' }}>
                <strong>Error:</strong> {download.error}
              </p>
            )}
          </div>
          
          <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: '0.5rem' }}>
            {getStatusBadge(download.status)}
            
            {download.status === 'downloading' && download.progress > 0 && (
              <div style={{ fontSize: '0.9rem', color: '#6c757d' }}>
                {download.progress}%
              </div>
            )}
            
            {download.status === 'completed' && download.downloadUrl && (
              <button
                className="btn btn-download"
                onClick={() => handleDownload(download.downloadUrl, download.filename)}
                style={{ fontSize: '0.8rem', padding: '0.5rem 1rem' }}
              >
                Download File
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}

export default DownloadHistory
