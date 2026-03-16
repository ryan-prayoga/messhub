export const OUTBOX_SYNC_TAG = 'messhub-outbox-sync';

const DB_NAME = 'messhub-pwa';
const DB_VERSION = 1;
const STORE_NAME = 'outbox';

export type QueuedRequest = {
  id: string;
  type: 'create-activity' | 'create-comment';
  url: string;
  method: 'POST';
  headers: Record<string, string>;
  body: string;
  createdAt: string;
};

export async function enqueueRequest(
  item: Omit<QueuedRequest, 'id' | 'createdAt'> & {
    id?: string;
    createdAt?: string;
  }
) {
  const entry: QueuedRequest = {
    ...item,
    id: item.id ?? resolveRequestID(),
    createdAt: item.createdAt ?? new Date().toISOString()
  };

  const database = await openDatabase();
  await requestToPromise(
    database
      .transaction(STORE_NAME, 'readwrite')
      .objectStore(STORE_NAME)
      .put(entry)
  );

  return entry;
}

export async function listQueuedRequests(): Promise<QueuedRequest[]> {
  const database = await openDatabase();
  const items = await requestToPromise<QueuedRequest[]>(
    database
      .transaction(STORE_NAME, 'readonly')
      .objectStore(STORE_NAME)
      .getAll()
  );

  return items.sort((left, right) => left.createdAt.localeCompare(right.createdAt));
}

export async function countQueuedRequests(): Promise<number> {
  const database = await openDatabase();
  return requestToPromise<number>(
    database
      .transaction(STORE_NAME, 'readonly')
      .objectStore(STORE_NAME)
      .count()
  );
}

export async function removeQueuedRequest(id: string) {
  const database = await openDatabase();
  await requestToPromise(
    database
      .transaction(STORE_NAME, 'readwrite')
      .objectStore(STORE_NAME)
      .delete(id)
  );
}

export async function clearQueuedRequests() {
  const database = await openDatabase();
  await requestToPromise(
    database
      .transaction(STORE_NAME, 'readwrite')
      .objectStore(STORE_NAME)
      .clear()
  );
}

export async function flushQueuedRequests() {
  const items = await listQueuedRequests();
  let processed = 0;

  for (const item of items) {
    let response: Response;

    try {
      response = await fetch(item.url, {
        method: item.method,
        headers: item.headers,
        body: item.body,
        credentials: 'same-origin'
      });
    } catch {
      break;
    }

    if (response.ok || shouldDiscardQueuedRequest(response.status)) {
      await removeQueuedRequest(item.id);
      processed += 1;
      continue;
    }

    break;
  }

  return {
    processed,
    remaining: await countQueuedRequests()
  };
}

async function openDatabase() {
  if (!('indexedDB' in globalThis)) {
    throw new Error('Offline queue is not supported in this browser');
  }

  return new Promise<IDBDatabase>((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION);

    request.onupgradeneeded = () => {
      const database = request.result;
      if (!database.objectStoreNames.contains(STORE_NAME)) {
        database.createObjectStore(STORE_NAME, {
          keyPath: 'id'
        });
      }
    };

    request.onsuccess = () => resolve(request.result);
    request.onerror = () =>
      reject(request.error ?? new Error('Failed to open offline queue database'));
  });
}

function requestToPromise<T>(request: IDBRequest<T>) {
  return new Promise<T>((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () =>
      reject(request.error ?? new Error('IndexedDB request failed'));
  });
}

function shouldDiscardQueuedRequest(status: number) {
  return status >= 400 && status < 500 && ![401, 403, 408, 409, 429].includes(status);
}

function resolveRequestID() {
  if ('crypto' in globalThis && 'randomUUID' in crypto) {
    return crypto.randomUUID();
  }

  return `queue-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}
