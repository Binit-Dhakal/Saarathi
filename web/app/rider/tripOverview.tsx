import { Button } from "@/components/ui/button";
import { Drawer, DrawerContent, DrawerFooter, DrawerHeader, DrawerTitle } from "@/components/ui/drawer";
import { CarPackage, Fare, TripStatus } from "@/lib/types";
import { useState } from "react";
import { Car, CarFront, Bus, Truck } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";

type TripOverviewProps = {
  fares: Fare[];
  tripStatus: TripStatus,
  onConfirm: (vehicleType: CarPackage) => void;
  onCancel: () => void;
};

const vehicleIcons: Record<string, React.ElementType> = {
  sedan: Car,
  suv: CarFront,
  van: Bus,
  luxury: Truck,
};

export default function TripOverview({
  fares,
  tripStatus,
  onConfirm,
  onCancel,
}: TripOverviewProps) {
  const [selectedVehicle, setSelectedVehicle] = useState<CarPackage | null>(null);

  return (
    <Drawer open={true} modal={false}>
      <DrawerContent className="rounded-t-2xl shadow-2xl border-t bg-white z-[999]  mx-5">
        <DrawerHeader className="pb-2 text-center">
          <DrawerTitle className="text-lg font-semibold text-gray-800">
            {tripStatus === "selecting" ? "Choose your ride" :
              tripStatus === "waiting" ? "Waiting for a driver" :
                "Trip Completed"
            }
          </DrawerTitle>
          {tripStatus === "selecting" && <p className="text-sm text-gray-500">
            Select a vehicle type to continue
          </p>}
        </DrawerHeader>

        <div className="px-4 py-2 space-y-3 max-h-[50vh] overflow-y-auto">
          {tripStatus === "selecting" && fares.map(({ Package, TotalPrice }) => {
            const isSelected = selectedVehicle === Package;
            const Icon = vehicleIcons[Package.toLowerCase()]
            return (
              <div
                key={Package}
                className={`flex justify-between items-center p-4 rounded-xl border transition-all cursor-pointer shadow-sm 
                ${isSelected
                    ? "border-blue-500 bg-blue-50 shadow-md"
                    : "border-gray-200 hover:bg-gray-50"
                  }`}
                onClick={() => setSelectedVehicle(Package)}
              >
                <div className="flex items-center gap-3">
                  <Icon className="w-5 h-5 text-gray-600" />
                  <span className="capitalize font-medium text-gray-700">
                    {Package.toLowerCase()}
                  </span>
                </div>
                <span className="font-semibold text-gray-800">
                  Rs.{TotalPrice}
                </span>
              </div>
            );
          })}

          {tripStatus === "waiting" && (
            <div className="flex flex-col items-center justify-center py-6 space-y-4">
              <Skeleton className="h-[100px] w-[250px] rounded-xl" />
              <p className="text-gray-600 text-sm text-center">
                Waiting for a driver to accept your trip...
              </p>
            </div>
          )}
        </div>

        <DrawerFooter className="flex justify-between gap-3 border-t pt-4">
          {tripStatus === "selecting" &&
            <Button
              className="w-full"
              disabled={!selectedVehicle}
              onClick={() => selectedVehicle && onConfirm(selectedVehicle)}
            >
              Confirm
            </Button>}
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
}

