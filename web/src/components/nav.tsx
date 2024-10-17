import React from 'react';

// Access environment variables
// @ts-ignore
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;
// @ts-ignore
const WEBSOCKET_BASE_URL = import.meta.env.VITE_WEBSOCKET_BASE_URL;

interface NavBarProps {
  onHome: () => void;
  isSocketConnected: boolean;
}

const NavBar: React.FC<NavBarProps> = ({ onHome, isSocketConnected }) => {
  return (
    <nav className="bg-gray-800 text-white p-4">
      <div className="container mx-auto flex justify-between items-center">
        <h1
          className="text-xl font-bold cursor-pointer hover:text-gray-300"
          onClick={onHome}
        >
          Sentinel
        </h1>
        <div className="text-sm flex items-center space-x-4">
          <div>
            <p>API: {API_BASE_URL}</p>
          </div>
          <div className="flex items-center">
            <p>WebSocket: {WEBSOCKET_BASE_URL}</p>
            <span
              className={`ml-2 h-3 w-3 rounded-full ${isSocketConnected ? 'bg-green-500' : 'bg-red-500'
                }`}
            ></span>
          </div>
          {/* GitHub Link */}
          <a
            href="https://github.com/EntropyLabsAI/sentinel"
            target="_blank"
            rel="noopener noreferrer"
            className="ml-4"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6 text-white hover:text-gray-300"
              viewBox="0 0 24 24"
              fill="currentColor"
            >
              <path
                fillRule="evenodd"
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.089.636-1.339-2.22-.253-4.555-1.113-4.555-4.947 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.378.202 2.397.1 2.65.64.7 1.028 1.595 1.028 2.688 0 3.842-2.339 4.69-4.566 4.936.359.31.678.921.678 1.856 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.523 2 12 2z"
                clipRule="evenodd"
              />
            </svg>
          </a>
        </div>
      </div>
    </nav>
  );
};

export default NavBar;
