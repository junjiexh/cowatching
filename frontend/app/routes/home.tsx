import type { Route } from "./+types/home";
import { Link } from "react-router";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Cowatching - Home" },
    { name: "description", content: "Welcome to Cowatching!" },
  ];
}

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800">
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-3xl mx-auto text-center">
          <h1 className="text-6xl font-bold text-gray-900 dark:text-white mb-6">
            Cowatching
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-12">
            Watch videos together with friends in real-time
          </p>

          <div className="space-y-4">
            <Link
              to="/player"
              className="inline-block px-8 py-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold text-lg rounded-lg shadow-lg transition-colors"
            >
              Go to Video Player
            </Link>

            <div className="mt-8 text-gray-500 dark:text-gray-400">
              <p className="text-sm">Upload and play videos in sync with others</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
