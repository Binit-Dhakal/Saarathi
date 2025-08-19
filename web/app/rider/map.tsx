"use client"
import { useState, useCallback } from "react";
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from "react-leaflet";
import InteractiveRouting from "./simplified-routing";
import L from "leaflet";

// Fix for default markers
const sourceIcon = new L.Icon({
  iconUrl: "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png",
  shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png",
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41]
});

const destinationIcon = new L.Icon({
  iconUrl: "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png",
  shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png",
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41]
});

interface MapClickHandlerProps {
  onMapClick: (lat: number, lng: number) => void;
}

function MapClickHandler({ onMapClick }: MapClickHandlerProps) {
  useMapEvents({
    click: (e) => {
      onMapClick(e.latlng.lat, e.latlng.lng);
    },
  });
  return null;
}

const Map = () => {
  const [source, setSource] = useState<[number, number] | null>(null);
  const [destination, setDestination] = useState<[number, number] | null>(null);
  const [clickMode, setClickMode] = useState<'source' | 'destination'>('source');
  const [routingKey, setRoutingKey] = useState(0); // Force re-render of routing component

  const handleMapClick = useCallback((lat: number, lng: number) => {
    if (clickMode === 'source') {
      setSource([lat, lng]);
      setClickMode('destination');
    } else {
      setDestination([lat, lng]);
      setClickMode('source');
    }
    // Increment key to force routing component to re-render
    setRoutingKey(prev => prev + 1);
  }, [clickMode]);

  const clearPoints = () => {
    setSource(null);
    setDestination(null);
    setClickMode('source');
    setRoutingKey(prev => prev + 1);
  };

  const swapPoints = () => {
    if (source && destination) {
      setSource(destination);
      setDestination(source);
      setRoutingKey(prev => prev + 1);
    }
  };

  return (
    <div className="relative w-full h-full">
      {/* Control Panel */}
      <div className="absolute top-4 left-4 z-[1000] bg-white p-3 rounded-lg shadow-lg">
        <h3 className="font-semibold mb-2">Route Planning</h3>
        <div className="text-sm mb-2">
          <div className={`p-2 rounded ${clickMode === 'source' ? 'bg-green-100 border-green-500' : 'bg-gray-100'} border mb-1`}>
            📍 Click to set {clickMode === 'source' ? 'SOURCE' : 'source'} {source && '✓'}
          </div>
          <div className={`p-2 rounded ${clickMode === 'destination' ? 'bg-red-100 border-red-500' : 'bg-gray-100'} border`}>
            🎯 Click to set {clickMode === 'destination' ? 'DESTINATION' : 'destination'} {destination && '✓'}
          </div>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setClickMode(clickMode === 'source' ? 'destination' : 'source')}
            className="px-3 py-1 bg-blue-500 text-white rounded text-xs hover:bg-blue-600"
          >
            Switch Mode
          </button>
          <button
            onClick={swapPoints}
            disabled={!source || !destination}
            className="px-3 py-1 bg-orange-500 text-white rounded text-xs hover:bg-orange-600 disabled:bg-gray-400"
          >
            Swap
          </button>
          <button
            onClick={clearPoints}
            className="px-3 py-1 bg-red-500 text-white rounded text-xs hover:bg-red-600"
          >
            Clear
          </button>
        </div>
        {source && destination && (
          <div className="mt-2 text-xs text-green-600">
            ✅ Route will be calculated automatically
          </div>
        )}
      </div>

      <MapContainer
        center={[27.7172, 85.3240]} // Kathmandu, Nepal
        zoom={13}
        style={{ height: "100%", width: "100%" }}
        scrollWheelZoom={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />

        <MapClickHandler onMapClick={handleMapClick} />

        {source && (
          <Marker position={source} icon={sourceIcon}>
            <Popup>
              <strong>Source</strong><br />
              Lat: {source[0].toFixed(4)}<br />
              Lng: {source[1].toFixed(4)}
            </Popup>
          </Marker>
        )}

        {destination && (
          <Marker position={destination} icon={destinationIcon}>
            <Popup>
              <strong>Destination</strong><br />
              Lat: {destination[0].toFixed(4)}<br />
              Lng: {destination[1].toFixed(4)}
            </Popup>
          </Marker>
        )}

        {source && destination && (
          <InteractiveRouting
            key={routingKey}
            source={source}
            destination={destination}
          />
        )}
      </MapContainer>
    </div>
  );
};

export default Map;
