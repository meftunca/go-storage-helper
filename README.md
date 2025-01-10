# go-ffmpeg-helper
# Media Converter API Documentation

This API allows you to convert, resize, crop, and manipulate video and image files. It is built using the FFmpeg library.

---

## Basic Information

- **API URL**: `http://localhost:3000/convert`
- **HTTP Method**: `GET`
- **Media Types**: `Video` and `Image`
- **Supported Formats**: `mp4`, `webm`, `webp`, `jpg`, `png`, `gif`, `avif`

---

## Parameters

The following parameters can be passed as query parameters to the API. Only the required parameters need to be used.

### General Parameters

| Parameter | Description                                                                 | Example Value     |
|-----------|-----------------------------------------------------------------------------|-------------------|
| `input`   | Path to the input file (required).                                          | `input.mp4`       |
| `format`  | Output format (`mp4`, `webp`, `gif`, `jpg`, etc.).                          | `webp`            |
| `width`   | Width of the output (in pixels).                                            | `1280`            |
| `height`  | Height of the output (in pixels).                                           | `720`             |
| `quality` | Quality of the output (0-100, default: 90).                                 | `80`              |
| `crop`    | Cropping dimensions (`w:h:x:y` format).                                     | `1920:1080:0:0`   |

### Video-Specific Parameters

| Parameter | Description                                                                 | Example Value     |
|-----------|-----------------------------------------------------------------------------|-------------------|
| `fps`     | Frame rate (FPS).                                                           | `30`              |
| `start`   | Start time for clipping (in seconds).                                       | `10`              |
| `end`     | End time for clipping (in seconds).                                         | `20`              |
| `bitrate` | Bitrate (e.g., `2M`).                                                       | `2M`              |
| `crf`     | Constant Rate Factor (0-51, lower values mean higher quality).              | `23`              |

---

## Warnings

- **Video-Specific Parameters**: Parameters like `fps`, `start`, `end`, `bitrate`, and `crf` are only valid for `VideoKind`. If used with `ImageKind`, a warning message will be returned, and the parameters will be ignored.
- **Required Parameter**: The `input` parameter is always required.

---

## Example Requests

### 1. **Convert Video to MP4**
```bash
GET http://localhost:3000/convert?input=input.mp4&format=mp4&width=1280&height=720&fps=30&bitrate=2M&crf=23
```
**Response:**


```bash   {    "message": "Conversion successful",    "outputFile": "./output/output.mp4",    "warnings": []  }   ```

### 2. **Convert Image to WebP**




```json    GET http://localhost:3000/convert?input=input.jpg&format=webp&width=800&height=600&quality=90   ```

**Response:**



```json   {    "message": "Conversion successful",    "outputFile": "./output/output.webp",    "warnings": []  } ```

### 3. **Create GIF from Video**

 

```bash   GET http://localhost:3000/convert?input=input.mp4&format=gif&width=320&height=240&fps=15&start=10&end=20 ```

**Response:**

```json    {    "message": "Conversion successful",    "outputFile": "./output/output.gif",    "warnings": []  } ```

### 4. **Convert Image with Video-Specific Parameters**

bashCopy

```bash   GET http://localhost:3000/convert?input=input.jpg&format=webp&width=800&height=600&fps=30&bitrate=2M&crf=23   ```

**Response:**

  ```json {    "message": "Conversion successful",    "outputFile": "./output/output.webp",    "warnings": [      "fps parameter is ignored for images",      "bitrate parameter is ignored for images",      "crf parameter is ignored for images"    ]  }   `

Error Messages
--------------

*   jsonCopy{ "error": "Input file does not exist"}
    
*   jsonCopy{ "error": "ffmpeg error: "}
    

Installation and Usage
----------------------

1.  **Install FFmpeg**: The API uses the FFmpeg library. Ensure FFmpeg is installed on your system.
    
    *   bashCopysudo apt install ffmpeg
        
2.  bashCopygo run main.go
    
3.  **Use the API**:
    
    *   Send requests to the API using a browser or tools like curl.
        
    *   bashCopycurl "http://localhost:3000/convert?input=input.mp4&format=mp4&width=1280&height=720"
