"use client"
import { useState, useEffect, useRef } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L from "leaflet";
import { useWS } from "@/context/WebSocketContext";
import { TripOffer } from "@/lib/types";
import TripOfferDrawer from "./tripOffer";

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

const driverIcon = new L.Icon({
  iconUrl: "./car-icon.svg",
  iconSize: [25, 41],
})


const Map = () => {
  const [loc, setLoc] = useState<[number, number] | null>(null);
  const locRef = useRef<[number, number] | null>(null);
  const { sendMessage, lastMessage } = useWS()
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const [offer, setOffer] = useState<TripOffer | null>(null);

  useEffect(() => {
    if (!lastMessage) return;

    if (lastMessage.event == "TRIP_OFFER_REQUEST") {
      const data = lastMessage.data as TripOffer;
      setOffer(data);
    }

  }, [lastMessage])


  useEffect(() => {
    const watchId = navigator.geolocation.watchPosition(
      (position) => {
        const newLoc: [number, number] = [position.coords.latitude, position.coords.longitude]
        setLoc(newLoc)
        locRef.current = newLoc;
      },
      (error) => {
        const fallback: [number, number] = [27.6922, 85.3344]
        console.log(error)
        setLoc(fallback)
        locRef.current = fallback;
      },
      { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 })

    const sendLocation = () => {
      if (!locRef.current) { return }
      sendMessage({
        event: "DRIVER_LOCATION_UPDATE",
        data: {
          "latitude": locRef.current[0],
          "longitude": locRef.current[1]
        }
      })
    }

    sendLocation()
    intervalRef.current = setInterval(sendLocation, 10000)

    return () => {
      navigator.geolocation.clearWatch(watchId)
      if (intervalRef.current) clearInterval(intervalRef.current);
    }
  }, [])

  const handleAccept = () => {
    if (!offer) return;
    sendMessage({
      event: "TRIP_ACCEPTED",
      data: { tripID: offer.tripID, offerID: offer.offerID }
    })
  }

  const handleReject = () => {
    if (!offer) return;
    sendMessage({
      event: "TRIP_REJECTED",
      data: { tripID: offer.tripID, offerID: offer.offerID }
    })
  }

  return (
    <div className="relative w-full h-full">
      <MapContainer
        center={[27.7172, 85.3240]}
        zoom={13}
        style={{ height: "100%", width: "100%" }}
        scrollWheelZoom={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        {loc && (
          <Marker position={loc} icon={driverIcon}>
            <Popup>
              <strong>Source</strong><br />
              Lat: {loc[0].toFixed(4)}<br />
              Lng: {loc[1].toFixed(4)}
            </Popup>
          </Marker>
        )}

        {offer && (
          <>
            <Marker position={[offer.pickUp[1], offer.pickUp[0]]} icon={sourceIcon}>
              <Popup>Pickup</Popup>
            </Marker>


            <Marker position={[offer.dropOff[1], offer.dropOff[0]]} icon={destinationIcon}>
              <Popup>Dropoff</Popup>
            </Marker>
          </>
        )}
      </MapContainer>
      <TripOfferDrawer offer={offer} setOffer={setOffer} onAccept={handleAccept} onReject={handleReject} />
    </div>
  );
};

export default Map;
