## Installing yt-dlp

You can install **yt-dlp** globally (in your system `PATH`) using one of the following methods:

---

### 1. Using pip (Cross-platform, Recommended)

```bash
pip install -U yt-dlp
```

If `yt-dlp` is not found in your shell after installation, ensure your local bin directory is in the `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Add the above line to your `~/.bashrc`, `~/.zshrc`, or equivalent shell config file.

---

### 2. Using the Standalone Binary (No Python Required)

```bash
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

Verify installation:

```bash
yt-dlp --version
```

---

### 3. On Windows

#### Using Scoop:

```powershell
scoop install yt-dlp
```

#### Using Chocolatey:

```powershell
choco install yt-dlp
```

#### Manual Installation:

* Download `yt-dlp.exe` from [yt-dlp Latest Releases](https://github.com/yt-dlp/yt-dlp/releases/latest).
* Place it in a folder (e.g., `C:\tools\yt-dlp\`).
* Add this folder to your system's **Environment Variables > PATH**.

---

## Useful Links

* [yt-dlp GitHub Repository](https://github.com/yt-dlp/yt-dlp)
* [yt-dlp Latest Releases](https://github.com/yt-dlp/yt-dlp/releases/latest)
