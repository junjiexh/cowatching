# Video Player Guide

This guide explains how to use the Cowatching video player.

## Features

- **Video Upload**: Upload videos up to 500MB
- **Video Playback**: Play videos using HTML5 video player
- **Video List**: Browse all uploaded videos with metadata
- **Video Management**: Delete videos you no longer need
- **Seek Support**: Fast-forward and rewind using the video player controls

## Getting Started

### 1. Start the Backend

```bash
cd backend
go run cmd/api/main.go
```

The backend will start on `http://localhost:8080`.

### 2. Start the Frontend

```bash
cd frontend
npm install  # First time only
npm run dev
```

The frontend will start on `http://localhost:5173` (or another port if 5173 is busy).

### 3. Access the Player

1. Open your browser to `http://localhost:5173`
2. Click "Go to Video Player" button
3. You'll see the player page with:
   - Main video player in the center
   - Upload section on the right
   - Video list below the upload section

## Using the Video Player

### Uploading Videos

1. Click "Click to select video" in the upload section
2. Choose a video file from your computer (any video format supported by your browser)
3. Enter a title for the video (optional - defaults to filename)
4. Click "Upload"
5. Wait for the upload to complete (progress bar shows status)
6. The video will appear in the video list automatically

### Playing Videos

1. Find the video you want to play in the video list on the right
2. Click on the video card
3. The video will load in the main player and start buffering
4. Use the HTML5 video controls to:
   - Play/Pause
   - Adjust volume
   - Seek to different timestamps
   - Toggle fullscreen

### Deleting Videos

1. Hover over a video in the list
2. Click the red trash icon that appears
3. Confirm the deletion
4. The video will be removed from the list and storage

## API Endpoints

The backend provides these endpoints:

- `GET /api/v1/videos` - List all videos
- `POST /api/v1/videos/upload` - Upload a new video
- `GET /api/v1/videos/stream/{filename}` - Stream a video
- `DELETE /api/v1/videos/{filename}` - Delete a video

## Storage

Videos are stored in `backend/uploads/videos/` directory. The directory is created automatically on first upload.

## Supported Video Formats

The player supports any video format that:
1. Your browser can play natively (typically MP4, WebM, Ogg)
2. Has a `video/*` MIME type

Common supported formats:
- MP4 (H.264)
- WebM
- Ogg

## Limitations

- Maximum upload size: 500MB per video
- Videos are stored on the local filesystem (not database)
- No video transcoding (videos must be in browser-compatible formats)
- No authentication (any user can upload/delete videos)

## Troubleshooting

### Video won't play

- Make sure the video format is supported by your browser
- Check browser console for errors
- Verify the backend is running on port 8080

### Upload fails

- Check video file is under 500MB
- Ensure backend has write permissions to `uploads/videos` directory
- Check backend logs for errors

### CORS errors

- Make sure backend is running on port 8080
- Check that frontend is accessing from `http://localhost:5173` or `http://localhost:3000`
- Backend CORS is configured for these origins

## Development Notes

### Backend Structure

```
backend/internal/handlers/video.go
├── Upload()  - Handles multipart file uploads
├── List()    - Returns all videos as JSON
├── Stream()  - Serves video files with range support
└── Delete()  - Removes video files
```

### Frontend Components

```
frontend/app/components/
├── VideoPlayer.tsx  - Main video player with HTML5 <video>
├── VideoList.tsx    - Scrollable list of videos
└── VideoUpload.tsx  - Upload form with progress tracking
```

### Video Metadata

Each video includes:
```typescript
{
  id: string;          // Timestamp-based ID
  title: string;       // User-provided or filename
  filename: string;    // Stored filename
  url: string;         // Streaming endpoint URL
  size: number;        // File size in bytes
  contentType: string; // MIME type
  uploadedAt: string;  // ISO timestamp
}
```

## Next Steps

Future enhancements could include:
- Video transcoding to ensure browser compatibility
- Thumbnail generation
- Video duration extraction
- User authentication and permissions
- Video categories/tags
- Search and filtering
- Playlist support
- Video quality selection
