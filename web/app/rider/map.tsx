"use client"
import { useState, useCallback, useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from "react-leaflet";
import L from "leaflet";
import { confirmRide, getRoute } from "@/lib/api";
import { CarPackage, FareEstimateResponse, TripStatus } from "@/lib/types";
import { RoutingControl } from "./routing-control";
import { convertCoordinates } from "@/lib/utils";
import TripOverview from "./tripOverview";

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
  const [route, setRoute] = useState<[number, number][]>([]);
  const [fare, setFare] = useState<FareEstimateResponse>();
  const [rideID, setRideID] = useState<string | null>(null);
  const [tripStatus, setTripStatus] = useState<TripStatus>("selecting");

  useEffect(() => {
    const fetchRoute = async () => {
      if (source && destination) {
        try {
          const data = await getRoute(source[0], source[1], destination[0], destination[1])
          setFare(data)
          setRoute(convertCoordinates(data.Geometry.coordinates))
          setTripStatus("selecting")
        } catch (err) {
          console.log("Failed to fetch route: ", err);
        }
      }
    }

    fetchRoute();
  }, [source, destination])

  const handleMapClick = useCallback((lat: number, lng: number) => {
    if (clickMode === 'source') {
      setSource([lat, lng]);
      setClickMode('destination');
    } else {
      setDestination([lat, lng]);
      setClickMode('source');
    }
  }, [clickMode]);

  const handleRideConfirm = async (carPackage: CarPackage) => {
    if (fare?.FareID == undefined) {
      return
    }

    try {
      const data = await confirmRide(fare?.FareID, carPackage)
      setRideID(data.RideID)
      setTripStatus("waiting")
    } catch (err) {
      console.log("Failed to confirm ride: ", err)
    }

    return
  }

  const clearPoints = () => {
    setSource(null);
    setDestination(null);
    setRoute([]);
    setClickMode('source');
  };

  const swapPoints = () => {
    if (source && destination) {
      setSource(destination);
      setDestination(source);
    }
  };

  return (
    <>
      <div className="relative w-full h-full">
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
          {route && <RoutingControl route={route} />}
        </MapContainer>
      </div>
      {fare?.Fares && fare?.Fares?.length > 0 &&
        <TripOverview
          fares={fare.Fares}
          tripStatus={tripStatus}
          onConfirm={handleRideConfirm}
          onCancel={() => console.log("cancelled")}
        />
      }
    </>
  );
};

export default Map;
