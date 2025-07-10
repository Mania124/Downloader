
# YouTube Downloader Backend (Go + Gin)

This is a production-hardened backend API for downloading YouTube videos or audio, extracting thumbnails, and streaming download progress. The service wraps around `yt-dlp` to provide these features via HTTP endpoints.

---

## Features

- ✅ **Download Video or Audio** with multiple format and quality options
  - **Video formats**: MP4, WebM, MKV, AVI, MOV, FLV, 3GP
  - **Quality options**: 360p, 480p, 720p, 1080p
  - **Audio formats**: MP3 with 192K quality
- ✅ **Fetch Thumbnail** of any valid YouTube video
- ✅ **Stream Download Progress** to clients via Server-Sent Events (SSE)
- ✅ **List Downloaded Files** with metadata and download URLs
- ✅ **Serve Downloaded Files** with proper streaming and content headers
- ✅ **Health Check Endpoint**
- ✅ Production-ready with input validation, context timeouts, and error handling
- ✅ CORS-configurable via environment variable

---

## Endpoints

### Health Check
```http
GET /
```
**Response:**
```json
{
  "success": true,
  "status": 200,
  "message": "The server was successfully connected"
}
```

---

### Download Video/Audio
```http
POST /download
```
**Request Body:**
```json
{
  "url": "https://youtube.com/...",
  "format": "video", // or "audio"
  "resolution": "720", // optional: "360", "480", "720", "1080"
  "videoFormat": "mp4" // optional: "mp4", "webm", "mkv", "avi", "best"
}
```
**Response:**
```json
{
  "message": "Download completed",
  "filename": "video.mp4",
  "size": 12345678,
  "downloadUrl": "/files/video.mp4"
}
```

---

### Get Thumbnail
```http
POST /thumbnail
```
**Request Body:**
```json
{
  "url": "https://youtube.com/..."
}
```
**Response:**
```json
{
  "thumbnail": "https://example.com/thumbnail.jpg"
}
```

---

### Stream Download Progress
```http
GET /download/stream?url=<VIDEO_URL>&format=video|audio&resolution=720&videoFormat=mp4
```
**Parameters:**
- `url`: Video URL (required)
- `format`: "video" or "audio" (required)
- `resolution`: "360", "480", "720", "1080" (optional)
- `videoFormat`: "mp4", "webm", "mkv", "avi", "best" (optional)

**Response:** Stream of download progress via SSE. When complete, includes file information.

---

### List Downloaded Files
```http
GET /files
```
**Response:**
```json
{
  "files": [
    {
      "name": "video.mp4",
      "size": 12345678,
      "modTime": "2025-07-10T16:30:00Z",
      "downloadUrl": "/files/video.mp4",
      "type": "video"
    }
  ],
  "count": 1
}
```

---

### Download File
```http
GET /files/{filename}
```
**Response:** Streams the requested file with appropriate headers for download.

---

## Format Selection

The API supports intelligent format selection based on your preferences:

### Video Formats
- **mp4**: Most compatible, works on all devices
- **webm**: Good compression, modern browsers
- **mkv**: High quality, supports multiple audio tracks
- **avi**: Legacy format, wide compatibility
- **mov**: Apple QuickTime format
- **flv**: Flash video format
- **3gp**: Mobile-optimized format
- **best**: Let yt-dlp choose the best available format

### Quality Options
- **360p**: Low quality, small file size
- **480p**: Standard definition
- **720p**: HD quality (recommended)
- **1080p**: Full HD quality
- **No specification**: Best available quality

### Format Priority
When both resolution and format are specified, the system will:
1. Try to find the exact format and resolution combination
2. Fall back to the specified resolution in any format
3. Fall back to the specified format in any resolution
4. Use the best available quality and format

---

## Environment Variables

| Variable         | Description                          | Default Value            |
|-----------------|--------------------------------------|------------------------|
| `FRONTEND_ORIGIN` | Allowed CORS Origin for Frontend     | `http://localhost:5173` |

---

## Running the Server

```bash
go run main.go
```

Or build:

```bash
go build -o downloader
./downloader
```

---

## Requirements

- Go 1.21+
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed and in `$PATH` (Check in [installation guide](how_to_download_yt-dlp.md))
- [ffmep](https://www.gyan.dev/ffmpeg/builds/) installed (Check in [installation guide](ffmpeg_installation.md))

---

## Testing

```bash
go test ./...
```

Unit tests exist for:
- URL validation
- Helper functions
- Thumbnail handler

---

## Project Structure

```
/downloader
│
├── cmd/               # Entrypoint
│   └── main.go
│
├── handlers/          # Route Handlers
│   ├── download.go
│   ├── thumbnail.go
│   ├── health.go
│   ├── download_progress.go
│
├── router/            # Routes Setup
│   └── routes.go
│
├── utils/             # Utility functions
│   ├── url.go
│   ├── helper.go
│
├── go.mod
└── go.sum
```

---

## Notes

- All `yt-dlp` commands are wrapped with Go contexts for timeout control.
- Environment-specific CORS origin setup via `FRONTEND_ORIGIN`.
- SSE used for download progress streaming.
- Production-ready error handling.

---

## License

MIT
