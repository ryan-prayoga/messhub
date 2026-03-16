/// <reference types="@sveltejs/kit" />

import { build, files, version } from '$service-worker';
import {
  clearQueuedRequests,
  flushQueuedRequests,
  OUTBOX_SYNC_TAG
} from '$lib/pwa/offline-queue';

declare const self: ServiceWorkerGlobalScope;

const STATIC_CACHE = `messhub-static-${version}`;
const PAGE_CACHE = `messhub-pages-${version}`;
const API_CACHE = `messhub-api-${version}`;
const OFFLINE_URL = '/offline';
const PRECACHE_URLS = [...build, ...files, OFFLINE_URL];
const STATIC_ASSETS = new Set([...build, ...files]);
const SAFE_PAGE_PATHS = ['/dashboard', '/feed', '/wallet', '/wifi', '/profile'];
const SAFE_API_PREFIXES = [
  '/api/v1/activities',
  '/api/v1/contributions/leaderboard',
  '/api/v1/wallet',
  '/api/v1/wallet/transactions',
  '/api/v1/wifi/active',
  '/api/v1/wifi/my'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(STATIC_CACHE);
      await cache.addAll(PRECACHE_URLS);
      await self.skipWaiting();
    })()
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    (async () => {
      const cacheKeys = await caches.keys();
      await Promise.all(
        cacheKeys
          .filter((cacheKey) => ![STATIC_CACHE, PAGE_CACHE, API_CACHE].includes(cacheKey))
          .map((cacheKey) => caches.delete(cacheKey))
      );
      await self.clients.claim();
    })()
  );
});

self.addEventListener('message', (event) => {
  const message = event.data as
    | {
        type?: string;
      }
    | undefined;

  if (message?.type === 'SYNC_OUTBOX') {
    event.waitUntil(flushQueuedRequests());
    return;
  }

  if (message?.type === 'CLEAR_AUTH_CACHES') {
    event.waitUntil(clearAuthenticatedState());
  }
});

self.addEventListener('sync', (event) => {
  if (event.tag !== OUTBOX_SYNC_TAG) {
    return;
  }

  event.waitUntil(flushQueuedRequests());
});

self.addEventListener('push', (event) => {
  if (!event.data) {
    return;
  }

  const payload = event.data.json() as {
    title?: string;
    body?: string;
    icon?: string;
    badge?: string;
    tag?: string;
    data?: {
      url?: string;
      type?: string;
      entity_id?: string;
    };
  };

  event.waitUntil(
    self.registration.showNotification(payload.title ?? 'MessHub', {
      body: payload.body ?? 'Ada update baru di MessHub.',
      icon: payload.icon ?? '/icons/icon-192.png',
      badge: payload.badge ?? '/icons/icon-192.png',
      tag: payload.tag ?? 'messhub',
      data: payload.data ?? {
        url: '/notifications'
      }
    })
  );
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();

  const targetPath = String(event.notification.data?.url ?? '/notifications');
  const targetURL = new URL(targetPath, self.location.origin).href;

  event.waitUntil(focusOrOpenWindow(targetURL));
});

self.addEventListener('fetch', (event) => {
  const { request } = event;
  if (request.method !== 'GET') {
    return;
  }

  const url = new URL(request.url);
  if (url.origin !== self.location.origin) {
    return;
  }

  if (request.mode === 'navigate') {
    event.respondWith(handleNavigationRequest(request));
    return;
  }

  if (STATIC_ASSETS.has(url.pathname)) {
    event.respondWith(staleWhileRevalidate(request, STATIC_CACHE));
    return;
  }

  if (isSafePageDataRequest(url)) {
    event.respondWith(networkFirst(request, PAGE_CACHE));
    return;
  }

  if (isSafeAPIRequest(url)) {
    event.respondWith(staleWhileRevalidate(request, API_CACHE));
  }
});

async function handleNavigationRequest(request: Request) {
  const url = new URL(request.url);
  if (!SAFE_PAGE_PATHS.includes(url.pathname)) {
    try {
      return await fetch(request);
    } catch {
      return (await caches.match(OFFLINE_URL)) ?? Response.error();
    }
  }

  const response = await fetch(request)
    .then(async (networkResponse) => {
      if (shouldCacheResponse(request.url, networkResponse)) {
        const cache = await caches.open(PAGE_CACHE);
        await cache.put(request, networkResponse.clone());
      }

      return networkResponse;
    })
    .catch(() => null);

  if (response) {
    return response;
  }

  const cachedResponse = await caches.match(request);
  if (cachedResponse) {
    return cachedResponse;
  }

  return (await caches.match(OFFLINE_URL)) ?? Response.error();
}

async function staleWhileRevalidate(request: Request, cacheName: string) {
  const cache = await caches.open(cacheName);
  const cachedResponse = await cache.match(request);
  const networkPromise = fetch(request)
    .then(async (networkResponse) => {
      if (shouldCacheResponse(request.url, networkResponse)) {
        await cache.put(request, networkResponse.clone());
      }

      return networkResponse;
    })
    .catch(() => null);

  if (cachedResponse) {
    void networkPromise;
    return cachedResponse;
  }

  const networkResponse = await networkPromise;
  return networkResponse ?? Response.error();
}

async function networkFirst(request: Request, cacheName: string) {
  const cache = await caches.open(cacheName);

  try {
    const networkResponse = await fetch(request);
    if (shouldCacheResponse(request.url, networkResponse)) {
      await cache.put(request, networkResponse.clone());
    }

    return networkResponse;
  } catch {
    const cachedResponse = await cache.match(request);
    return cachedResponse ?? Response.error();
  }
}

function isSafePageDataRequest(url: URL) {
  if (!url.pathname.endsWith('/__data.json')) {
    return false;
  }

  return SAFE_PAGE_PATHS.some((path) => url.pathname === `${path}/__data.json`);
}

function isSafeAPIRequest(url: URL) {
  return SAFE_API_PREFIXES.some((prefix) => url.pathname.startsWith(prefix));
}

function shouldCacheResponse(requestURL: string, response: Response) {
  if (!response.ok) {
    return false;
  }

  if (response.redirected) {
    const redirectedURL = new URL(response.url, self.location.origin);
    if (redirectedURL.pathname === '/login') {
      return false;
    }
  }

  const url = new URL(requestURL);
  return url.pathname !== '/api/v1/auth/me' && url.pathname !== '/api/v1/notifications';
}

async function focusOrOpenWindow(targetURL: string) {
  const windowClients = await self.clients.matchAll({
    type: 'window',
    includeUncontrolled: true
  });

  for (const client of windowClients) {
    if ('focus' in client) {
      await client.navigate(targetURL);
      return client.focus();
    }
  }

  return self.clients.openWindow(targetURL);
}

async function clearAuthenticatedState() {
  await Promise.all([caches.delete(PAGE_CACHE), caches.delete(API_CACHE), clearQueuedRequests()]);
}
