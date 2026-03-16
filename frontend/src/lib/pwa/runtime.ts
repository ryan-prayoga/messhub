import { browser } from '$app/environment';
import { invalidateAll } from '$app/navigation';
import { clearQueuedRequests, flushQueuedRequests, OUTBOX_SYNC_TAG } from '$lib/pwa/offline-queue';
import { ensurePushSubscription } from '$lib/pwa/push';

let initialized = false;

export function initializePwaRuntime() {
  if (!browser || initialized) {
    return () => {};
  }

  initialized = true;

  const syncAndRefresh = async () => {
    if (!navigator.onLine) {
      return;
    }

    const result = await flushQueuedRequests().catch(() => ({ processed: 0, remaining: 0 }));
    if (result.processed > 0) {
      await invalidateAll();
    }
  };

  const handleOnline = () => {
    void syncAndRefresh();
  };

  const handleVisibilityChange = () => {
    if (document.visibilityState === 'visible') {
      void syncAndRefresh();
    }
  };

  window.addEventListener('online', handleOnline);
  document.addEventListener('visibilitychange', handleVisibilityChange);

  void syncAndRefresh();
  if (typeof Notification !== 'undefined' && Notification.permission === 'granted') {
    void ensurePushSubscription();
  }

  return () => {
    window.removeEventListener('online', handleOnline);
    document.removeEventListener('visibilitychange', handleVisibilityChange);
    initialized = false;
  };
}

export async function scheduleOutboxSync() {
  if (!browser) {
    return;
  }

  if ('serviceWorker' in navigator) {
    const registration = await getServiceWorkerRegistration();

    if (registration?.sync?.register) {
      try {
        await registration.sync.register(OUTBOX_SYNC_TAG);
        return;
      } catch {
        // Fall back to a foreground flush below when Background Sync is unavailable.
      }
    }

    registration?.active?.postMessage({
      type: 'SYNC_OUTBOX'
    });
  }

  if (navigator.onLine) {
    const result = await flushQueuedRequests().catch(() => ({ processed: 0, remaining: 0 }));
    if (result.processed > 0) {
      await invalidateAll();
    }
  }
}

export async function clearOfflineSessionArtifacts() {
  if (!browser) {
    return;
  }

  await clearQueuedRequests().catch(() => {});

  if (!('serviceWorker' in navigator)) {
    return;
  }

  const registration = await getServiceWorkerRegistration();
  if (!registration) {
    return;
  }

  registration?.active?.postMessage({
    type: 'CLEAR_AUTH_CACHES'
  });
}

async function getServiceWorkerRegistration() {
  const existing = await navigator.serviceWorker.getRegistration();
  if (existing) {
    return existing;
  }

  return Promise.race([
    navigator.serviceWorker.ready,
    new Promise<null>((resolve) => {
      setTimeout(() => resolve(null), 1500);
    })
  ]);
}
