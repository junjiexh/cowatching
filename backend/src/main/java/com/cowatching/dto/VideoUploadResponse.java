package com.cowatching.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class VideoUploadResponse {
    private Long videoId;
    private String message;
    private String filename;
    private Long fileSize;
}
