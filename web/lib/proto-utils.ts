import { fromBinary, type DescMessage, type MessageShape } from "@bufbuild/protobuf";

/**
 * Decode base64 protobuf data to typed message
 */
export function decodeProtoMessage<T extends DescMessage>(
  schema: T,
  base64Data: string
): MessageShape<T> {
  // convert base64 -> Uint8Array (works in browser and Node)
  let bytes: Uint8Array;
  if (typeof window !== "undefined" && typeof window.atob === "function") {
    // browser
    bytes = Uint8Array.from(atob(base64Data), (c) => c.charCodeAt(0));
  } else {
    // Node / SSR (Next.js server-side)
    bytes = Uint8Array.from(Buffer.from(base64Data, "base64"));
  }

  return fromBinary(schema, bytes);
}

