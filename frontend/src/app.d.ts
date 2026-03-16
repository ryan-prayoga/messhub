declare global {
  namespace App {
    interface Locals {
      token: string | null;
      user: import('$lib/api/types').SessionUser | null;
    }

    interface PageData {
      user: App.Locals['user'];
    }
  }
}

export {};
