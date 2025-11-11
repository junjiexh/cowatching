import { useState, useEffect } from 'react';
import type { Route } from './+types/player';
import { VideoPlayer } from '../components/VideoPlayer';
import { VideoList, type Video } from '../components/VideoList';
import { VideoUpload } from '../components/VideoUpload';

const API_URL = 'http://localhost:8080';

export function meta({}: Route.MetaArgs) {
  return [
    { title: 'Video Player - Cowatching' },
    { name: 'description', content: 'Watch and manage your videos' },
  ];
}

export default function Player() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [currentVideo, setCurrentVideo] = useState<Video | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchVideos = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await fetch(`${API_URL}/api/v1/videos`);
      if (!response.ok) {
        throw new Error('Failed to fetch videos');
      }
      const data = await response.json();
      setVideos(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load videos');
      console.error('Error fetching videos:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchVideos();
  }, []);

  const handleVideoSelect = (video: Video) => {
    setCurrentVideo(video);
  };

  const handleVideoDelete = async (video: Video) => {
    if (!confirm(`Are you sure you want to delete "${video.title}"?`)) {
      return;
    }

    try {
      const response = await fetch(`${API_URL}/api/v1/videos/${video.id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Failed to delete video');
      }

      // If the deleted video was currently playing, clear it
      if (currentVideo?.id === video.id) {
        setCurrentVideo(null);
      }

      // Refresh the video list
      await fetchVideos();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete video');
      console.error('Error deleting video:', err);
    }
  };

  const handleUploadSuccess = () => {
    fetchVideos();
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <header className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
        <div className="px-6 py-4">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Cowatching Video Player
          </h1>
        </div>
      </header>

      {/* Main Content */}
      <div className="p-6">
        <div className="flex gap-6 h-[calc(100vh-120px)]">
          {/* Main Player Section */}
          <div className="flex-1 flex flex-col gap-6">
            <div className="flex-1">
              <VideoPlayer
                videoUrl={currentVideo ? `${API_URL}${currentVideo.url}` : null}
                title={currentVideo?.title}
              />
            </div>
          </div>

          {/* Right Sidebar */}
          <div className="w-96 flex flex-col gap-6">
            {/* Upload Section */}
            <VideoUpload onUploadSuccess={handleUploadSuccess} apiUrl={API_URL} />

            {/* Video List Section */}
            <div className="flex-1 bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden flex flex-col">
              <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                  Videos ({videos.length})
                </h2>
              </div>
              <div className="flex-1 overflow-hidden">
                {error ? (
                  <div className="flex items-center justify-center h-full p-4">
                    <div className="text-center text-red-600 dark:text-red-400">
                      <p className="font-medium">Error loading videos</p>
                      <p className="text-sm mt-1">{error}</p>
                      <button
                        onClick={fetchVideos}
                        className="mt-3 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm"
                      >
                        Retry
                      </button>
                    </div>
                  </div>
                ) : (
                  <VideoList
                    videos={videos}
                    currentVideoId={currentVideo?.id || null}
                    onVideoSelect={handleVideoSelect}
                    onVideoDelete={handleVideoDelete}
                    isLoading={isLoading}
                  />
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
