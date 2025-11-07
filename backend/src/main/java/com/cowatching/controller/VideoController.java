package com.cowatching.controller;

import com.cowatching.dto.VideoDto;
import com.cowatching.dto.VideoUploadResponse;
import com.cowatching.entity.User;
import com.cowatching.entity.Video;
import com.cowatching.service.UserService;
import com.cowatching.service.VideoService;
import com.cowatching.service.VideoStorageService;
import lombok.RequiredArgsConstructor;
import org.springframework.core.io.Resource;
import org.springframework.core.io.UrlResource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.net.MalformedURLException;
import java.nio.file.Path;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/videos")
@RequiredArgsConstructor
public class VideoController {

    private final VideoService videoService;
    private final UserService userService;
    private final VideoStorageService videoStorageService;

    @PostMapping("/upload")
    public ResponseEntity<?> uploadVideo(
            @RequestParam("file") MultipartFile file,
            @RequestParam("title") String title,
            @RequestParam(value = "description", required = false) String description,
            Authentication authentication) {

        if (authentication == null || !authentication.isAuthenticated()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not authenticated");
        }

        try {
            User user = userService.findByUsername(authentication.getName())
                    .orElseThrow(() -> new RuntimeException("User not found"));

            Video video = videoService.uploadVideo(file, title, description, user);

            VideoUploadResponse response = new VideoUploadResponse();
            response.setVideoId(video.getId());
            response.setMessage("Video uploaded successfully");
            response.setFilename(video.getFilename());
            response.setFileSize(video.getFileSize());

            return ResponseEntity.status(HttpStatus.CREATED).body(response);
        } catch (Exception e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Failed to upload video: " + e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(error);
        }
    }

    @GetMapping("/my-videos")
    public ResponseEntity<?> getMyVideos(Authentication authentication) {
        if (authentication == null || !authentication.isAuthenticated()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not authenticated");
        }

        try {
            User user = userService.findByUsername(authentication.getName())
                    .orElseThrow(() -> new RuntimeException("User not found"));

            List<VideoDto> videos = videoService.getUserVideos(user);
            return ResponseEntity.ok(videos);
        } catch (Exception e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Failed to retrieve videos: " + e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(error);
        }
    }

    @GetMapping("/{videoId}")
    public ResponseEntity<?> getVideoInfo(@PathVariable Long videoId, Authentication authentication) {
        if (authentication == null || !authentication.isAuthenticated()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not authenticated");
        }

        try {
            Video video = videoService.getVideoById(videoId);

            // Check if the video belongs to the authenticated user
            if (!video.getUser().getUsername().equals(authentication.getName())) {
                return ResponseEntity.status(HttpStatus.FORBIDDEN).body("Access denied");
            }

            VideoDto videoDto = new VideoDto();
            videoDto.setId(video.getId());
            videoDto.setTitle(video.getTitle());
            videoDto.setDescription(video.getDescription());
            videoDto.setFilename(video.getFilename());
            videoDto.setContentType(video.getContentType());
            videoDto.setFileSize(video.getFileSize());
            videoDto.setUploadedAt(video.getUploadedAt());
            videoDto.setUploaderUsername(video.getUser().getUsername());

            return ResponseEntity.ok(videoDto);
        } catch (Exception e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(error);
        }
    }

    @GetMapping("/{videoId}/stream")
    public ResponseEntity<?> streamVideo(@PathVariable Long videoId, Authentication authentication) {
        if (authentication == null || !authentication.isAuthenticated()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not authenticated");
        }

        try {
            Video video = videoService.getVideoById(videoId);

            // Check if the video belongs to the authenticated user
            if (!video.getUser().getUsername().equals(authentication.getName())) {
                return ResponseEntity.status(HttpStatus.FORBIDDEN).body("Access denied");
            }

            Path filePath = videoStorageService.loadFile(video.getFilename());
            Resource resource = new UrlResource(filePath.toUri());

            if (!resource.exists() || !resource.isReadable()) {
                return ResponseEntity.status(HttpStatus.NOT_FOUND).body("Video file not found");
            }

            String contentType = video.getContentType();
            if (contentType == null) {
                contentType = "application/octet-stream";
            }

            return ResponseEntity.ok()
                    .contentType(MediaType.parseMediaType(contentType))
                    .header(HttpHeaders.CONTENT_DISPOSITION, "inline; filename=\"" + video.getFilename() + "\"")
                    .body(resource);

        } catch (MalformedURLException e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", "Invalid file path: " + e.getMessage());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(error);
        } catch (Exception e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.status(HttpStatus.NOT_FOUND).body(error);
        }
    }

    @DeleteMapping("/{videoId}")
    public ResponseEntity<?> deleteVideo(@PathVariable Long videoId, Authentication authentication) {
        if (authentication == null || !authentication.isAuthenticated()) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Not authenticated");
        }

        try {
            User user = userService.findByUsername(authentication.getName())
                    .orElseThrow(() -> new RuntimeException("User not found"));

            videoService.deleteVideo(videoId, user);

            Map<String, String> response = new HashMap<>();
            response.put("message", "Video deleted successfully");
            return ResponseEntity.ok(response);
        } catch (Exception e) {
            Map<String, String> error = new HashMap<>();
            error.put("error", e.getMessage());
            return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(error);
        }
    }
}
