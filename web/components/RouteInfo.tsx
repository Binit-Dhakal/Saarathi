import React from "react";

interface RouteInfoProps {
  summary: { distance: string; time: string };
  onClose?: () => void;
  onBookRide?: () => void;
}

export const RouteInfo: React.FC<RouteInfoProps> = ({ summary, onClose, onBookRide }) => {
  return (
    <div
      className="fixed bottom-4 left-4 z-50 bg-white p-4 rounded-lg shadow-lg border-l-4 border-green-500 animate-fade-in"
      style={{ minWidth: "240px" }}
    >
      <div className="flex items-start">
        <div className="text-green-600 mr-3 text-xl">✅</div>
        <div className="flex-1">
          <h4 className="font-semibold text-gray-800 mb-1">Route Calculated</h4>
          <div className="text-sm text-gray-600 mb-2">
            <div><strong>Distance:</strong> {summary.distance}</div>
            <div><strong>Time:</strong> {summary.time}</div>
          </div>

          {onBookRide && (
            <button
              onClick={onBookRide}
              className="w-full px-3 py-2 bg-blue-500 text-white rounded text-sm hover:bg-blue-600"
            >
              🚕 Book Ride
            </button>
          )}
        </div>
        <button
          onClick={onClose}
          className="ml-3 text-gray-400 hover:text-gray-600"
        >
          ✕
        </button>
      </div>
    </div>
  );
};

