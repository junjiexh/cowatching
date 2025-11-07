package com.cowatching.repository;

import com.cowatching.entity.Video;
import com.cowatching.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface VideoRepository extends JpaRepository<Video, Long> {
    List<Video> findByUserOrderByUploadedAtDesc(User user);
    List<Video> findByUserIdOrderByUploadedAtDesc(Long userId);
}
