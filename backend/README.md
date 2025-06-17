
# YouTube Downloader Backend (Go + Gin)

This is a production-hardened backend API for downloading YouTube videos or audio, extracting thumbnails, and streaming download progress. The service wraps around `yt-dlp` to provide these features via HTTP endpoints.

---

## Features

- ✅ **Download Video or Audio** in specified formats
- ✅ **Fetch Thumbnail** of any valid YouTube video
- ✅ **Stream Download Progress** to clients via Server-Sent Events (SSE)
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
  "resolution": "720" // optional, only for video
}
```
**Response:**
```json
{
  "message": "Download completed"
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
GET /download/stream?url=<VIDEO_URL>&format=video|audio
```
**Response:** Stream of download progress via SSE.

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
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed and in `$PATH`

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
