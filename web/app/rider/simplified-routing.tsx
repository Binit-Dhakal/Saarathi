// components/SimplifiedRouting.tsx
"use client";
import { useState } from "react";
import { useMap } from "react-leaflet";
import { useRoutingControl } from "@/hooks/useRoutingControl";
import { RouteInfo } from "@/components/RouteInfo";

interface SimplifiedRoutingProps {
  source: [number, number];
  destination: [number, number];
}

export default function SimplifiedRouting({ source, destination }: SimplifiedRoutingProps) {
  const map = useMap();
  const [summary, setSummary] = useState<{ distance: string; time: string } | null>(null);
  const [showInfo, setShowInfo] = useState(true);

  useRoutingControl(
    map,
    source,
    destination,
    (s) => {
      setSummary(s);
      setShowInfo(true);
    },
    (errMsg) => {
      alert(`Routing Error: ${errMsg}\n\nPlease check:\n- OSRM server\n- Valid locations`);
    }
  );

  if (!summary || !showInfo) return null;

  return <RouteInfo summary={summary} onClose={() => setShowInfo(false)} />;
}

