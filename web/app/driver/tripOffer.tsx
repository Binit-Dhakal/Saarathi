import { TripOffer } from "@/lib/types";
import { useEffect, useState } from "react";
import { Drawer, DrawerContent, DrawerHeader, DrawerFooter, DrawerTitle } from "@/components/ui/drawer"
import { Button } from "@/components/ui/button"
import { Progress } from "@/components/ui/progress"

const TripOfferDrawer = ({ offer, setOffer, onAccept, onReject }: {
  offer: TripOffer | null,
  setOffer: React.Dispatch<React.SetStateAction<TripOffer | null>>,
  onAccept: () => void,
  onReject: () => void
}) => {
  const [secondsLeft, setSecondsLeft] = useState<number>(0);
  const [totalSeconds, setTotalSeconds] = useState<number>(0);

  useEffect(() => {
    if (!offer) return;

    const expiry = new Date(offer.expiresAt).getTime()
    const total = Math.floor((expiry - Date.now()) / 1000);

    setTotalSeconds(total);
    setSecondsLeft(total);

    const interval = setInterval(() => {
      const remaining = Math.max(0, Math.floor((expiry - Date.now()) / 1000));
      setSecondsLeft(remaining);
      if (remaining == 0) {
        clearInterval(interval)
        setOffer(null)
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
          <p>Pickup: {offer.pickUp[1].toFixed(4)}, {offer.pickUp[0].toFixed(4)}</p>
          <p>Dropoff: {offer.dropOff[1].toFixed(4)}, {offer.dropOff[0].toFixed(4)}</p>
          <p>Price: ${0}</p>
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
