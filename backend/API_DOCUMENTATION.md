# Video Upload and Streaming API Documentation

## Overview
This is a Spring Boot backend API for video uploading and streaming. Users can register, login, upload videos, view their video list, and stream videos.

## Architecture

### Entities
- **User**: Stores user account information (username, email, password)
- **Video**: Stores video metadata (title, description, filename, file path, content type, file size, upload date)

### Key Components
- **Controllers**: Handle HTTP requests (AuthController, VideoController)
- **Services**: Business logic layer (UserService, VideoService, VideoStorageService)
- **Repositories**: Data access layer (UserRepository, VideoRepository)
- **DTOs**: Data transfer objects for API requests/responses
- **Security**: Spring Security with BCrypt password encoding

### Storage
- Videos are stored in the file system at `./uploads/videos/`
- Video metadata is stored in H2 database (in-memory for development)
- Maximum file size: 500MB

## API Endpoints

### Authentication Endpoints

#### 1. Register User
```
POST /api/auth/register
Content-Type: application/json

Request Body:
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}

Response (201 Created):
{
  "message": "User registered successfully",
  "username": "testuser",
  "email": "test@example.com"
}

Error Response (400 Bad Request):
{
  "error": "Username already exists"
}
```

#### 2. Login User
```
POST /api/auth/login
Content-Type: application/json

Request Body:
{
  "username": "testuser",
  "password": "password123"
}

Response (200 OK):
{
  "message": "Login successful",
  "username": "testuser"
}

Error Response (401 Unauthorized):
{
  "error": "Invalid username or password"
}
```

#### 3. Get Current User
```
GET /api/auth/me
Authorization: Basic <base64(username:password)>

Response (200 OK):
{
  "username": "testuser"
}

Error Response (401 Unauthorized):
"Not authenticated"
```

#### 4. Logout
```
POST /api/auth/logout
Authorization: Basic <base64(username:password)>

Response (200 OK):
{
  "message": "Logout successful"
}
```

### Video Endpoints

#### 5. Upload Video
```
POST /api/videos/upload
Authorization: Basic <base64(username:password)>
Content-Type: multipart/form-data

Request Form Data:
- file: (video file)
- title: "My Video Title"
- description: "Video description" (optional)

Response (201 Created):
{
  "videoId": 1,
  "message": "Video uploaded successfully",
  "filename": "testuser_123e4567-e89b-12d3-a456-426614174000.mp4",
  "fileSize": 15728640
}

Error Response (401 Unauthorized):
"Not authenticated"

Error Response (500 Internal Server Error):
{
  "error": "Failed to upload video: <error message>"
}
```

#### 6. Get My Videos
```
GET /api/videos/my-videos
Authorization: Basic <base64(username:password)>

Response (200 OK):
[
  {
    "id": 1,
    "title": "My Video Title",
    "description": "Video description",
    "filename": "testuser_123e4567-e89b-12d3-a456-426614174000.mp4",
    "contentType": "video/mp4",
    "fileSize": 15728640,
    "uploadedAt": "2025-11-07T10:30:00",
    "uploaderUsername": "testuser"
  }
]

Error Response (401 Unauthorized):
"Not authenticated"
```

#### 7. Get Video Info
```
GET /api/videos/{videoId}
Authorization: Basic <base64(username:password)>

Response (200 OK):
{
  "id": 1,
  "title": "My Video Title",
  "description": "Video description",
  "filename": "testuser_123e4567-e89b-12d3-a456-426614174000.mp4",
  "contentType": "video/mp4",
  "fileSize": 15728640,
  "uploadedAt": "2025-11-07T10:30:00",
  "uploaderUsername": "testuser"
}

Error Response (401 Unauthorized):
"Not authenticated"

Error Response (403 Forbidden):
"Access denied"

Error Response (404 Not Found):
{
  "error": "Video not found with id: 1"
}
```

#### 8. Stream Video
```
GET /api/videos/{videoId}/stream
Authorization: Basic <base64(username:password)>

Response (200 OK):
Content-Type: video/mp4
Content-Disposition: inline; filename="testuser_123e4567-e89b-12d3-a456-426614174000.mp4"
[Binary video data]

Error Response (401 Unauthorized):
"Not authenticated"

Error Response (403 Forbidden):
"Access denied"

Error Response (404 Not Found):
"Video file not found"
```

#### 9. Delete Video
```
DELETE /api/videos/{videoId}
Authorization: Basic <base64(username:password)>

Response (200 OK):
{
  "message": "Video deleted successfully"
}

Error Response (401 Unauthorized):
"Not authenticated"

Error Response (400 Bad Request):
{
  "error": "You don't have permission to delete this video"
}
```

### Health Check

#### 10. Health Check
```
GET /api/health

Response (200 OK):
{
  "status": "UP",
  "timestamp": "2025-11-07T10:30:00"
}
```

## Authentication

This API uses HTTP Basic Authentication. For all protected endpoints, you need to include the Authorization header:

```
Authorization: Basic <base64_encoded_credentials>
```

Where `<base64_encoded_credentials>` is the Base64 encoding of `username:password`.

Example:
```
Username: testuser
Password: password123
Base64: dGVzdHVzZXI6cGFzc3dvcmQxMjM=
Header: Authorization: Basic dGVzdHVzZXI6cGFzc3dvcmQxMjM=
```

## Testing with curl

### 1. Register a user
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 3. Upload a video
```bash
curl -X POST http://localhost:8080/api/videos/upload \
  -u testuser:password123 \
  -F "file=@/path/to/video.mp4" \
  -F "title=My Test Video" \
  -F "description=This is a test video"
```

### 4. Get your videos
```bash
curl -X GET http://localhost:8080/api/videos/my-videos \
  -u testuser:password123
```

### 5. Stream a video
```bash
curl -X GET http://localhost:8080/api/videos/1/stream \
  -u testuser:password123 \
  --output downloaded_video.mp4
```

### 6. Delete a video
```bash
curl -X DELETE http://localhost:8080/api/videos/1 \
  -u testuser:password123
```

## Running the Application

### Prerequisites
- Java 17 or higher
- Maven 3.6 or higher

### Build and Run
```bash
cd backend
mvn clean install
mvn spring-boot:run
```

The application will start on `http://localhost:8080`

### H2 Database Console
Access the H2 console at: `http://localhost:8080/h2-console`
- JDBC URL: `jdbc:h2:mem:cowatchingdb`
- Username: `sa`
- Password: (leave empty)

## Configuration

Edit `application.properties` to configure:
- Server port
- Database settings
- File upload limits
- Video storage location
- CORS settings

## Security Notes

1. **Password Storage**: Passwords are hashed using BCrypt before storing in the database
2. **Access Control**: Users can only view, stream, and delete their own videos
3. **CORS**: Configured to allow requests from `http://localhost:5173` and `http://localhost:3000`
4. **File Upload**: Maximum file size is 500MB (configurable)

## Future Enhancements

Consider adding:
- JWT tokens for better authentication
- Video transcoding for multiple quality levels
- Thumbnail generation
- Video sharing between users
- Comments and likes
- PostgreSQL/MySQL for production database
- Cloud storage (AWS S3, Google Cloud Storage)
- Rate limiting
- Video compression
