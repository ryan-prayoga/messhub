import { browser } from '$app/environment';
import { readable } from 'svelte/store';

export const networkState = readable(
  {
    online: true
  },
  (set) => {
    if (!browser) {
      return;
    }

    const update = () =>
      set({
        online: navigator.onLine
      });

    update();
    window.addEventListener('online', update);
    window.addEventListener('offline', update);

    return () => {
      window.removeEventListener('online', update);
      window.removeEventListener('offline', update);
    };
  }
);
