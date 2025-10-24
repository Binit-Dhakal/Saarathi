"use client";

import { Drawer, DrawerContent, DrawerHeader, DrawerTitle, DrawerFooter } from "@/components/ui/drawer";
import { Button } from "@/components/ui/button";
import { MapPin, User, CarFront } from "lucide-react";
import { RiderUpdatePayload } from "@/gen/riderspb/riders_messages_pb";

type DriverAssignedDrawerProps = {
  data: RiderUpdatePayload;
  onCancel: () => void;
};

export default function DriverAssignedDrawer({ data, onCancel }: DriverAssignedDrawerProps) {
  return (
    <Drawer open={true} modal={false}>
      <DrawerContent className="rounded-t-2xl shadow-2xl border-t bg-white z-[999] mx-5">
        <DrawerHeader className="pb-2 text-center">
          <DrawerTitle className="text-lg font-semibold text-gray-800">
            Your driver is on the way 🚗
          </DrawerTitle>
          <p className="text-sm text-gray-500">Track driver’s progress in real-time</p>
        </DrawerHeader>

        <div className="px-5 py-3 space-y-4 max-h-[50vh] overflow-y-auto">
          {/* Driver info */}
          <div className="flex items-center justify-between p-3 rounded-xl border border-gray-200 shadow-sm bg-gray-50">
            <div className="flex items-center gap-3">
              <User className="w-5 h-5 text-gray-600" />
              <div>
                <p className="font-semibold text-gray-800">{data.driverName}</p>
                <p className="text-sm text-gray-500">Driver</p>
              </div>
            </div>
            <div className="text-sm text-gray-700 text-right">
              <p>{data.vehicleMake} {data.vehicleModel}</p>
              <p className="font-medium">{data.vehicleNumber}</p>
            </div>
          </div>

          {/* Trip info */}
          <div className="p-3 rounded-xl border border-gray-200 shadow-sm bg-gray-50 space-y-2">
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-green-500" />
              <p className="text-sm text-gray-700">Pickup: {data.pickup?.lat?.toFixed(4)}, {data.pickup?.lng?.toFixed(4)}</p>
            </div>
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-red-500" />
              <p className="text-sm text-gray-700">Dropoff: {data.dropoff?.lat?.toFixed(4)}, {data.dropoff?.lng?.toFixed(4)}</p>
            </div>
            <div className="flex justify-between pt-2 border-t border-gray-200">
              <p className="text-gray-600 text-sm">Distance</p>
              <p className="font-medium">{data.distance.toFixed(2)} km</p>
            </div>
            <div className="flex justify-between">
              <p className="text-gray-600 text-sm">Fare</p>
              <p className="font-semibold text-green-600"> {data.price}</p>
            </div>
          </div>
        </div>

        <DrawerFooter className="border-t pt-4 flex justify-center">
          <Button variant="destructive" className="w-full" onClick={onCancel}>
            Cancel Ride
          </Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}

