import { browser } from '$app/environment';
import { pushApi } from '$lib/api/client';
import { PUSH_VAPID_PUBLIC_KEY } from '$lib/config/env';

export type PushPromptState =
  | 'unsupported'
  | 'missing-key'
  | 'permission-needed'
  | 'denied'
  | 'enabled';

export async function getPushPromptState(): Promise<PushPromptState> {
  if (!browser || !supportsPushNotifications()) {
    return 'unsupported';
  }

  if (!PUSH_VAPID_PUBLIC_KEY) {
    return 'missing-key';
  }

  if (Notification.permission === 'denied') {
    return 'denied';
  }

  const registration = await getServiceWorkerRegistration();
  if (!registration) {
    return 'permission-needed';
  }

  const subscription = await registration.pushManager.getSubscription();

  if (subscription) {
    return 'enabled';
  }

  return 'permission-needed';
}

export async function ensurePushSubscription() {
  if (!browser || !supportsPushNotifications() || !PUSH_VAPID_PUBLIC_KEY) {
    return false;
  }

  if (Notification.permission !== 'granted') {
    return false;
  }

  const registration = await getServiceWorkerRegistration();
  if (!registration) {
    throw new Error('Service worker belum siap untuk push subscription.');
  }

  let subscription = await registration.pushManager.getSubscription();

  if (!subscription) {
    subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: decodeBase64Url(PUSH_VAPID_PUBLIC_KEY)
    });
  }

  const payload = subscription.toJSON();
  const p256dh = payload.keys?.p256dh;
  const auth = payload.keys?.auth;

  if (!payload.endpoint || !p256dh || !auth) {
    throw new Error('Push subscription is missing required keys');
  }

  await pushApi.subscribe({
    endpoint: payload.endpoint,
    keys: {
      p256dh,
      auth
    }
  });

  return true;
}

export async function requestPushPermissionAndSubscribe() {
  if (!browser || !supportsPushNotifications()) {
    return {
      state: 'unsupported' as const,
      message: 'Browser ini belum mendukung push notification.'
    };
  }

  if (!PUSH_VAPID_PUBLIC_KEY) {
    return {
      state: 'missing-key' as const,
      message: 'VAPID public key belum dikonfigurasi.'
    };
  }

  const permission = await Notification.requestPermission();
  if (permission !== 'granted') {
    return {
      state: permission === 'denied' ? ('denied' as const) : ('permission-needed' as const),
      message:
        permission === 'denied'
          ? 'Izin notifikasi diblokir di browser.'
          : 'Izin notifikasi belum diberikan.'
    };
  }

  await ensurePushSubscription();

  return {
    state: 'enabled' as const,
    message: 'Push notification aktif di perangkat ini.'
  };
}

function supportsPushNotifications() {
  return (
    'Notification' in globalThis &&
    'serviceWorker' in navigator &&
    'PushManager' in globalThis
  );
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

function decodeBase64Url(value: string) {
  const padding = '='.repeat((4 - (value.length % 4)) % 4);
  const base64 = (value + padding).replace(/-/g, '+').replace(/_/g, '/');
  const raw = atob(base64);
  const bytes = new Uint8Array(raw.length);

  for (let index = 0; index < raw.length; index += 1) {
    bytes[index] = raw.charCodeAt(index);
  }

  return bytes;
}
