import { Trash2 } from 'lucide-react';

export interface Video {
  id: number;
  title: string;
  url: string;
  size: number;
  contentType?: string;
  uploadedAt: string;
}

interface VideoListProps {
  videos: Video[];
  currentVideoId: number | null;
  onVideoSelect: (video: Video) => void;
  onVideoDelete: (video: Video) => void;
  isLoading?: boolean;
}

export function VideoList({
  videos,
  currentVideoId,
  onVideoSelect,
  onVideoDelete,
  isLoading = false,
}: VideoListProps) {
  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-gray-400">Loading videos...</div>
      </div>
    );
  }

  if (videos.length === 0) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center text-gray-400">
          <svg
            className="mx-auto h-16 w-16 mb-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"
            />
          </svg>
          <p>No videos yet</p>
          <p className="text-sm mt-1">Upload a video to get started</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-y-auto">
      <div className="space-y-2 p-2">
        {videos.map((video) => (
          <div
            key={video.id}
            className={`
              group relative rounded-lg border-2 transition-all cursor-pointer
              ${
                currentVideoId === video.id
                  ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                  : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
              }
            `}
            onClick={() => onVideoSelect(video)}
          >
            <div className="p-3">
              <div className="flex items-start justify-between">
                <div className="flex-1 min-w-0">
                  <h3 className="font-medium text-sm text-gray-900 dark:text-white truncate">
                    {video.title}
                  </h3>
                  <div className="mt-1 space-y-1">
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                      {formatFileSize(video.size)}
                    </p>
                    <p className="text-xs text-gray-400 dark:text-gray-500">
                      {formatDate(video.uploadedAt)}
                    </p>
                  </div>
                </div>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onVideoDelete(video);
                  }}
                  className="ml-2 p-1.5 rounded-md opacity-0 group-hover:opacity-100 transition-opacity
                           hover:bg-red-100 dark:hover:bg-red-900/30 text-red-600 dark:text-red-400"
                  title="Delete video"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
