const CACHE_NAME = "blocks-v1";
const ASSETS = [
  "/",
  "/index.html",
  "/wasm_exec.js",
  "/blocks.wasm",
  "/logo.png",
  "/logo-16.png",
  "/logo-32.png",
  "/logo-152.png",
  "/logo-167.png",
  "/logo-180.png",
  "/logo-192.png",
  "/logo-512.png",
  "/manifest.json",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
  );
});

self.addEventListener("fetch", (event) => {
  // Remove query string from the request URL
  const requestURL = new URL(event.request.url);
  const cleanRequest = new Request(requestURL.origin + requestURL.pathname);

  event.respondWith(
    caches
      .match(cleanRequest)
      .then((response) => response || fetch(event.request))
  );
});
