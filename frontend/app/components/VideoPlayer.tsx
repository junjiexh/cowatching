import { useEffect, useRef } from 'react';

interface VideoPlayerProps {
  videoUrl: string | null;
  title?: string;
}

export function VideoPlayer({ videoUrl, title }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    if (videoRef.current && videoUrl) {
      videoRef.current.load();
    }
  }, [videoUrl]);

  return (
    <div className="w-full h-full flex flex-col bg-black rounded-lg overflow-hidden">
      {videoUrl ? (
        <>
          <div className="relative flex-1 bg-black">
            <video
              ref={videoRef}
              className="w-full h-full"
              controls
              controlsList="nodownload"
            >
              <source src={videoUrl} type="video/mp4" />
              Your browser does not support the video tag.
            </video>
          </div>
          {title && (
            <div className="bg-gray-900 px-4 py-3 text-white">
              <h2 className="text-lg font-semibold">{title}</h2>
            </div>
          )}
        </>
      ) : (
        <div className="flex-1 flex items-center justify-center text-gray-400">
          <div className="text-center">
            <svg
              className="mx-auto h-24 w-24 mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <p className="text-xl">Select a video to play</p>
          </div>
        </div>
      )}
    </div>
  );
}
