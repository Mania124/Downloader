## FFmpeg Installation Guide

### Why FFmpeg?

This project requires **FFmpeg** for audio and video processing (used by `yt-dlp` for format conversion and post-processing tasks).

---

### Installation Instructions

#### **Ubuntu / Debian**

```bash
sudo apt update
sudo apt install ffmpeg
```

#### **macOS (using Homebrew)**

```bash
brew install ffmpeg
```

#### **Windows**

1. Download the latest static build from:
   [https://www.gyan.dev/ffmpeg/builds/](https://www.gyan.dev/ffmpeg/builds/)

2. Extract the ZIP file.

3. Add the `bin` folder (inside the extracted directory) to your system **PATH**:

   * Control Panel → System → Advanced System Settings → Environment Variables → Edit `Path`.
   * Add the full path to the `bin` folder.

4. Verify installation:

```bash
ffmpeg -version
```

---

### More details:

🔗 [Official FFmpeg Download Page](https://ffmpeg.org/download.html)
