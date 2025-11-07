package com.cowatching.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class VideoDto {
    private Long id;
    private String title;
    private String description;
    private String filename;
    private String contentType;
    private Long fileSize;
    private LocalDateTime uploadedAt;
    private String uploaderUsername;
}
