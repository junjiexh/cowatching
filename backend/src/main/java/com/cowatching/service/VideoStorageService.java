package com.cowatching.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.UUID;

@Service
public class VideoStorageService {

    private final Path storageLocation;

    public VideoStorageService(@Value("${video.storage.location}") String storageLocation) {
        this.storageLocation = Paths.get(storageLocation).toAbsolutePath().normalize();
        try {
            Files.createDirectories(this.storageLocation);
        } catch (IOException e) {
            throw new RuntimeException("Could not create storage directory", e);
        }
    }

    public String storeFile(MultipartFile file, String username) {
        if (file.isEmpty()) {
            throw new RuntimeException("Failed to store empty file");
        }

        String originalFilename = file.getOriginalFilename();
        String fileExtension = "";
        if (originalFilename != null && originalFilename.contains(".")) {
            fileExtension = originalFilename.substring(originalFilename.lastIndexOf("."));
        }

        String newFilename = username + "_" + UUID.randomUUID() + fileExtension;

        try {
            Path targetLocation = this.storageLocation.resolve(newFilename);
            Files.copy(file.getInputStream(), targetLocation, StandardCopyOption.REPLACE_EXISTING);
            return newFilename;
        } catch (IOException e) {
            throw new RuntimeException("Failed to store file", e);
        }
    }

    public Path loadFile(String filename) {
        return storageLocation.resolve(filename).normalize();
    }

    public void deleteFile(String filename) {
        try {
            Path file = loadFile(filename);
            Files.deleteIfExists(file);
        } catch (IOException e) {
            throw new RuntimeException("Failed to delete file", e);
        }
    }

    public Path getStorageLocation() {
        return storageLocation;
    }
}
