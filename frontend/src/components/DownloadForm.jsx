import { useState } from 'react'

const DownloadForm = ({ onDownloadStart, onDownloadUpdate, onDownloadComplete }) => {
  const [formData, setFormData] = useState({
    url: '',
    format: 'video',
    resolution: '720',
    videoFormat: 'mp4'
  })
  const [isDownloading, setIsDownloading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [progressText, setProgressText] = useState('')
  const [error, setError] = useState('')

  const handleInputChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({ ...prev, [name]: value }))
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!formData.url.trim()) {
      setError('Please enter a valid URL')
      return
    }

    setIsDownloading(true)
    setError('')
    setProgress(0)
    setProgressText('')

    const downloadId = Date.now().toString()
    
    onDownloadStart({
      id: downloadId,
      url: formData.url,
      format: formData.format,
      resolution: formData.resolution,
      videoFormat: formData.videoFormat,
      status: 'downloading',
      progress: 0,
      startTime: new Date()
    })

    try {
      // Use streaming endpoint for real-time progress
      const response = await fetch(`/api/download/stream?${new URLSearchParams({
        url: formData.url,
        format: formData.format,
        resolution: formData.resolution,
        videoFormat: formData.videoFormat
      })}`)

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === 'completed') {
              setProgress(100)
              setProgressText('Download completed!')
              onDownloadUpdate(downloadId, { 
                status: 'completed', 
                progress: 100,
                endTime: new Date()
              })
              onDownloadComplete()
              break
            } else if (data.startsWith('{')) {
              try {
                const fileInfo = JSON.parse(data)
                if (fileInfo.filename) {
                  onDownloadUpdate(downloadId, { 
                    filename: fileInfo.filename,
                    downloadUrl: fileInfo.downloadUrl
                  })
                }
              } catch (e) {
                // Ignore JSON parse errors for progress data
              }
            } else {
              // Progress update
              setProgressText(data)
              const match = data.match(/(\d+(?:\.\d+)?)%/)
              if (match) {
                setProgress(parseFloat(match[1]))
                onDownloadUpdate(downloadId, { progress: parseFloat(match[1]) })
              }
            }
          }
        }
      }
    } catch (error) {
      console.error('Download failed:', error)
      setError(`Download failed: ${error.message}`)
      onDownloadUpdate(downloadId, { 
        status: 'error', 
        error: error.message,
        endTime: new Date()
      })
    } finally {
      setIsDownloading(false)
      setProgress(0)
      setProgressText('')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="form-group">
        <label htmlFor="url">Video URL</label>
        <input
          type="url"
          id="url"
          name="url"
          className="form-control"
          placeholder="https://www.youtube.com/watch?v=..."
          value={formData.url}
          onChange={handleInputChange}
          disabled={isDownloading}
          required
        />
      </div>

      <div className="form-row">
        <div className="form-group">
          <label htmlFor="format">Format</label>
          <select
            id="format"
            name="format"
            className="form-control"
            value={formData.format}
            onChange={handleInputChange}
            disabled={isDownloading}
          >
            <option value="video">Video</option>
            <option value="audio">Audio (MP3)</option>
          </select>
        </div>

        {formData.format === 'video' && (
          <div className="form-group">
            <label htmlFor="resolution">Quality</label>
            <select
              id="resolution"
              name="resolution"
              className="form-control"
              value={formData.resolution}
              onChange={handleInputChange}
              disabled={isDownloading}
            >
              <option value="360">360p</option>
              <option value="480">480p</option>
              <option value="720">720p (HD)</option>
              <option value="1080">1080p (Full HD)</option>
              <option value="">Best Available</option>
            </select>
          </div>
        )}
      </div>

      {formData.format === 'video' && (
        <div className="form-group">
          <label htmlFor="videoFormat">Video Format</label>
          <select
            id="videoFormat"
            name="videoFormat"
            className="form-control"
            value={formData.videoFormat}
            onChange={handleInputChange}
            disabled={isDownloading}
          >
            <option value="mp4">MP4 (Most Compatible)</option>
            <option value="webm">WebM (Good Compression)</option>
            <option value="mkv">MKV (High Quality)</option>
            <option value="avi">AVI (Legacy)</option>
            <option value="mov">MOV (QuickTime)</option>
            <option value="flv">FLV (Flash)</option>
            <option value="3gp">3GP (Mobile)</option>
            <option value="best">Best Available</option>
          </select>
        </div>
      )}

      {error && (
        <div className="status-badge status-error" style={{ marginBottom: '1rem', display: 'block' }}>
          {error}
        </div>
      )}

      {isDownloading && (
        <div className="progress-container">
          <div className="progress-bar">
            <div 
              className="progress-fill" 
              style={{ width: `${progress}%` }}
            ></div>
          </div>
          <p style={{ marginTop: '0.5rem', fontSize: '0.9rem', color: '#6c757d' }}>
            {progressText || 'Starting download...'}
          </p>
        </div>
      )}

      <button 
        type="submit" 
        className="btn btn-primary" 
        disabled={isDownloading}
        style={{ width: '100%', marginTop: '1rem' }}
      >
        {isDownloading ? 'Downloading...' : 'Download'}
      </button>
    </form>
  )
}

export default DownloadForm
