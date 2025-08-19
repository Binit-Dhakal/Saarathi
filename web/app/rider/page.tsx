"use client"
import 'leaflet/dist/leaflet.css'

// import Map from "./map"
import dynamic from 'next/dynamic';

const DynamicMap = dynamic(() => import("./map"), {
  loading: () => (
    <div className="w-full h-96 bg-gray-200 flex items-center justify-center">
      <div className="text-center">
        <div className="w-8 h-8 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-2"></div>
        <p>Loading interactive map...</p>
      </div>
    </div>
  ),
  ssr: false, // Fixed: use false instead of !!false
});

function RiderHomeMap() {
  return (
    <div className="w-full h-[100vw]"> {/* Fixed: added specific height */}
      <DynamicMap />
    </div>
  )
}

export default RiderHomeMap
