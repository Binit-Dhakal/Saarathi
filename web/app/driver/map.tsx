"use client"
import { useState, useEffect, useRef } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import L from "leaflet";
import { useWS } from "@/context/WebSocketContext";
import TripOfferDrawer from "./tripOffer";
import { DriverUpdatePayload, DriverUpdatePayloadSchema, TripOffer, TripOfferSchema } from "@/gen/driverspb/drivers_messages_pb";
import { decodeProtoMessage } from "@/lib/proto-utils";
import TripDetailsDrawer from "./tripDetail";

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
  const [acceptedTrip, setAcceptedTrip] = useState<DriverUpdatePayload | null>(null);

  useEffect(() => {
    if (!lastMessage) return;

    if (lastMessage.event == "TRIP_OFFER_REQUEST") {
      const base64 = lastMessage.data as string;
      const data = decodeProtoMessage(TripOfferSchema, base64)
      setOffer(data);
    } else if (lastMessage.event == "TRIP_ACCEPT_PAYLOAD") {
      const base64Msg = lastMessage.data as string
      const data = decodeProtoMessage(DriverUpdatePayloadSchema, base64Msg)
      console.log(data)
      setAcceptedTrip(data);
      setOffer(null);
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
      data: { tripID: offer.tripId, offerID: offer.offerId }
    })
  }

  const handleReject = () => {
    if (!offer) return;
    sendMessage({
      event: "TRIP_REJECTED",
      data: { tripID: offer.tripId, offerID: offer.offerId }
    })
  }

  const handleStartTrip = () => {
    if (!acceptedTrip) return;
    sendMessage({
      event: "TRIP_STARTED",
      data: { tripID: acceptedTrip.tripId }
    });
  }

  const handleEndTrip = () => {
    if (!acceptedTrip) return;
    sendMessage({
      event: "TRIP_COMPLETED",
      data: { tripID: acceptedTrip.tripId }
    });
    setAcceptedTrip(null);
  };

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
            {offer.pickUp &&
              <Marker position={[offer.pickUp?.lat, offer.pickUp?.lng]} icon={sourceIcon}>
                <Popup>Pickup</Popup>
              </Marker>}


            {offer.dropOff &&
              <Marker position={[offer.dropOff?.lat, offer.dropOff?.lng]} icon={destinationIcon}>
                <Popup>Dropoff</Popup>
              </Marker>
            }
          </>
        )}

        {acceptedTrip && (
          <>
            {
              acceptedTrip.pickUp &&
              <Marker position={[acceptedTrip.pickUp?.lat, acceptedTrip.pickUp?.lng]} icon={sourceIcon}>
                <Popup>Pickup</Popup>
              </Marker>
            }

            {
              acceptedTrip.dropOff &&
              <Marker position={[acceptedTrip.dropOff?.lat, acceptedTrip.dropOff?.lng]} icon={destinationIcon}>
                <Popup>Dropoff</Popup>
              </Marker>

            }

          </>
        )}
      </MapContainer>
      <TripOfferDrawer offer={offer} setOffer={setOffer} onAccept={handleAccept} onReject={handleReject} />
      <TripDetailsDrawer trip={acceptedTrip} onStartTrip={handleStartTrip} onEndTrip={handleEndTrip} />
    </div>
  );
};

export default Map;
