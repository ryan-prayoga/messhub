declare global {
  namespace App {
    interface Locals {
      user: {
        id: string;
        email: string;
        name: string;
        role: 'admin' | 'treasurer' | 'member';
      } | null;
    }

    interface PageData {
      user: App.Locals['user'];
    }
  }
}

export {};
