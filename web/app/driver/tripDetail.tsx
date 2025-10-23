"use client";

import { Drawer, DrawerContent, DrawerHeader, DrawerTitle, DrawerFooter } from "@/components/ui/drawer";
import { Button } from "@/components/ui/button";
import { DriverUpdatePayload } from "@/gen/driverspb/drivers_messages_pb";

interface TripDetailsDrawerProps {
  trip: DriverUpdatePayload | null;
  onStartTrip: () => void;
  onEndTrip: () => void;
}

const TripDetailsDrawer = ({ trip, onStartTrip, onEndTrip }: TripDetailsDrawerProps) => {
  if (!trip) return null;

  return (
    <Drawer open={!!trip} modal={false}>
      <DrawerContent className="rounded-t-2xl shadow-2xl border-t bg-white z-[999] mx-5">
        <DrawerHeader>
          <DrawerTitle>🚘 Trip Details</DrawerTitle>
          <div className="text-gray-700 space-y-1 mt-2">
            {trip.riderName && <p><strong>Passenger:</strong> {trip.riderName}</p>}
            {trip.riderNumber && <p><strong>Phone:</strong> {trip.riderNumber}</p>}
            {trip.pickUp && (
              <p>
                <strong>Pickup:</strong> {trip.pickUp.lat.toFixed(4)}, {trip.pickUp.lng.toFixed(4)}
              </p>
            )}
            {trip.dropOff && (
              <p>
                <strong>Dropoff:</strong> {trip.dropOff.lat.toFixed(4)}, {trip.dropOff.lng.toFixed(4)}
              </p>
            )}
            {trip.price && <p><strong>Fare:</strong> ${trip.price.toFixed(2)}</p>}
          </div>
        </DrawerHeader>

        <DrawerFooter className="flex flex-col gap-2">
          <Button onClick={onStartTrip} className="bg-green-600 hover:bg-green-700">
            Start Trip
          </Button>
          <Button onClick={onEndTrip} variant="destructive">
            End Trip
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

export default TripDetailsDrawer;

