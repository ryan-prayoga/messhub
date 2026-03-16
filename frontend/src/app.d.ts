declare global {
  interface BeforeInstallPromptEvent extends Event {
    prompt: () => Promise<void>;
    userChoice: Promise<{
      outcome: 'accepted' | 'dismissed';
      platform: string;
    }>;
  }

  interface SyncManager {
    register: (tag: string) => Promise<void>;
  }

  interface ServiceWorkerRegistration {
    sync?: SyncManager;
  }

  namespace App {
    interface Locals {
      token: string | null;
      user: import('$lib/api/types').SessionUser | null;
    }

    interface PageData {
      user: App.Locals['user'];
      notificationSummary?: import('$lib/api/types').NotificationList;
    }
  }
}

export {};
