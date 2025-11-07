package com.cowatching.service;

import com.cowatching.dto.VideoDto;
import com.cowatching.entity.User;
import com.cowatching.entity.Video;
import com.cowatching.repository.VideoRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class VideoService {

    private final VideoRepository videoRepository;
    private final VideoStorageService videoStorageService;

    @Transactional
    public Video uploadVideo(MultipartFile file, String title, String description, User user) {
        String filename = videoStorageService.storeFile(file, user.getUsername());
        String filePath = videoStorageService.getStorageLocation().resolve(filename).toString();

        Video video = new Video();
        video.setTitle(title);
        video.setDescription(description);
        video.setFilename(filename);
        video.setFilePath(filePath);
        video.setContentType(file.getContentType());
        video.setFileSize(file.getSize());
        video.setUser(user);

        return videoRepository.save(video);
    }

    public List<VideoDto> getUserVideos(User user) {
        List<Video> videos = videoRepository.findByUserOrderByUploadedAtDesc(user);
        return videos.stream()
                .map(this::convertToDto)
                .collect(Collectors.toList());
    }

    public Video getVideoById(Long videoId) {
        return videoRepository.findById(videoId)
                .orElseThrow(() -> new RuntimeException("Video not found with id: " + videoId));
    }

    @Transactional
    public void deleteVideo(Long videoId, User user) {
        Video video = getVideoById(videoId);

        if (!video.getUser().getId().equals(user.getId())) {
            throw new RuntimeException("You don't have permission to delete this video");
        }

        videoStorageService.deleteFile(video.getFilename());
        videoRepository.delete(video);
    }

    private VideoDto convertToDto(Video video) {
        VideoDto dto = new VideoDto();
        dto.setId(video.getId());
        dto.setTitle(video.getTitle());
        dto.setDescription(video.getDescription());
        dto.setFilename(video.getFilename());
        dto.setContentType(video.getContentType());
        dto.setFileSize(video.getFileSize());
        dto.setUploadedAt(video.getUploadedAt());
        dto.setUploaderUsername(video.getUser().getUsername());
        return dto;
    }
}
