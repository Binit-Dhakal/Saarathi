import { useEffect, useState } from "react";
import { TripOffer } from "@/gen/driverspb/drivers_messages_pb";
import { Drawer, DrawerContent, DrawerHeader, DrawerFooter, DrawerTitle } from "@/components/ui/drawer"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"
import { timestampDate } from "@bufbuild/protobuf/wkt";


const TripOfferDrawer = ({ offer, setOffer, onAccept, onReject, onTimeout }: {
  offer: TripOffer | null,
  setOffer: React.Dispatch<React.SetStateAction<TripOffer | null>>,
  onAccept: () => void,
  onReject: () => void,
  onTimeout: () => void
}) => {
  const [secondsLeft, setSecondsLeft] = useState<number>(0);
  const [totalSeconds, setTotalSeconds] = useState<number>(0);


  useEffect(() => {
    if (!offer) return;
    if (!offer.expiresAt) return;

    const expiry = timestampDate(offer.expiresAt).getTime()
    const total = Math.floor((expiry - Date.now()) / 1000);

    setTotalSeconds(total);
    setSecondsLeft(total);

    const interval = setInterval(() => {
      const remaining = Math.max(0, Math.floor((expiry - Date.now()) / 1000));
      setSecondsLeft(remaining);
      if (remaining == 0) {
        clearInterval(interval)
        setOffer(null)
        onTimeout()
      }
    }, 1000)

    return () => clearInterval(interval)
  }, [offer])

  if (!offer) return null;

  const progressValue = totalSeconds > 0 ? (secondsLeft / totalSeconds) * 100 : 0;

  return (
    <Drawer open={!!offer} modal={false}>
      <DrawerContent className="rounded-t-2xl shadow-2xl border-t bg-white z-[999]  mx-5">
        <DrawerHeader>
          <DrawerTitle>🚖 New Trip Offer</DrawerTitle>
          <p>Pickup: {offer.pickUp?.lat.toFixed(4)}, {offer.pickUp?.lng.toFixed(4)}</p>
          <p>Dropoff: {offer.dropOff?.lat.toFixed(4)}, {offer.dropOff?.lng.toFixed(4)}</p>
          <p>Price: ${offer.price}</p>
          <div className="mt-3">
            <Progress value={progressValue} className="h-2" />
            <p className="text-sm text-red-500 font-bold mt-1">
              Expires in {secondsLeft}s
            </p>
          </div>
        </DrawerHeader>
        <DrawerFooter>
          <Button onClick={onAccept} className="bg-green-500">Accept</Button>
          <Button onClick={onReject} variant="destructive">Reject</Button>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  )
}

export default TripOfferDrawer;
