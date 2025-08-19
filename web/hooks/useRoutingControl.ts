import { useEffect, useRef } from "react";
import L from "leaflet";
import "leaflet-routing-machine";


export function useRoutingControl(
  map: L.Map | null,
  source: [number, number] | null,
  destination: [number, number] | null,
  onRouteFound: (summary: { distance: string; time: string }) => void,
  onError?: (message: string) => void
) {
  const routingRef = useRef<L.Routing.Control | null>(null);

  useEffect(() => {
    if (!map) return;

    if (!routingRef.current) {
      // Create control only if it does not exist
      const routingControl = L.Routing.control({
        waypoints: [],
        routeWhileDragging: false,
        addWaypoints: false,
        show: false,
        router: L.Routing.osrmv1({
          serviceUrl: "http://saarathi.com:8080/osrm/route/v1",
          profile: "driving",
        }),
        showAlternatives: false,
        fitSelectedRoutes: true,
      });

      routingControl.on("routesfound", (e: any) => {
        const route = e.routes[0];
        if (!route) return;
        const summary = route.summary;
        onRouteFound({
          distance: (summary.totalDistance / 1000).toFixed(2) + " km",
          time: Math.round(summary.totalTime / 60) + " minutes",
        });
        if (route.coordinates?.length > 0) {
          map.fitBounds(L.latLngBounds(route.coordinates).pad(0.1));
        }
      });

      routingControl.on("routingerror", (e: any) => {
        console.error("Routing error:", e);
        onError?.(e.error?.message || "Could not calculate route");
      });

      routingRef.current = routingControl;

      // Add control to map if not already attached
      routingRef.current.addTo(map);
    }

    if (routingRef.current && (routingRef.current as any)._map && source && destination) {
      routingRef.current.setWaypoints([
        L.latLng(source[0], source[1]),
        L.latLng(destination[0], destination[1]),
      ]);
    }
    return () => {
      if (routingRef.current && map) {
        map.removeControl(routingRef.current);
      }
    };
  }, [source, destination, map]);

  return null
}

