<script lang="ts">
  import { APP_NAME } from '$lib/config/env';

  type NavItem = {
    href: string;
    label: string;
  };

  export let user: App.Locals['user'];
  export let currentPath: string;

  const navItems: NavItem[] = [
    { href: '/', label: 'Home' },
    { href: '/members', label: 'Members' },
    { href: '/wallet', label: 'Wallet' },
    { href: '/wifi', label: 'Wifi' },
    { href: '/feed', label: 'Feed' }
  ];

  const utilityItems: NavItem[] = [
    { href: '/shared-expenses', label: 'Shared' },
    { href: '/contributions', label: 'Contrib' },
    { href: '/proposals', label: 'Proposals' },
    { href: '/profile', label: 'Profile' },
    { href: '/settings', label: 'Settings' }
  ];
</script>

<div class="mx-auto flex min-h-screen max-w-md flex-col bg-canvas text-ink shadow-shell">
  <header class="sticky top-0 z-10 border-b border-line/80 bg-panel/95 px-5 pb-4 pt-5 backdrop-blur">
    <div class="flex items-start justify-between gap-4">
      <div>
        <p class="text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">Mess PWA</p>
        <h1 class="mt-1 text-2xl font-bold">{APP_NAME}</h1>
      </div>
      {#if user}
        <div class="rounded-2xl border border-line bg-slate-50 px-3 py-2 text-right text-xs">
          <p class="font-semibold">{user.name}</p>
          <p class="text-slate-500">{user.role}</p>
        </div>
      {/if}
    </div>
  </header>

  <main class="flex-1 px-4 pb-28 pt-5">
    <slot />
  </main>

  <aside class="border-t border-line bg-white px-4 py-4">
    <div class="mb-4 grid grid-cols-5 gap-2">
      {#each navItems as item}
        <a
          href={item.href}
          class={`rounded-2xl px-3 py-2 text-center text-xs font-semibold ${
            currentPath === item.href ? 'bg-ink text-white' : 'bg-slate-100 text-slate-600'
          }`}
        >
          {item.label}
        </a>
      {/each}
    </div>

    <div class="grid grid-cols-5 gap-2">
      {#each utilityItems as item}
        <a
          href={item.href}
          class={`rounded-2xl px-3 py-2 text-center text-xs font-semibold ${
            currentPath === item.href ? 'bg-accent text-white' : 'bg-orange-50 text-orange-700'
          }`}
        >
          {item.label}
        </a>
      {/each}
    </div>

    <form method="POST" action="/logout" class="mt-4">
      <button
        type="submit"
        class="w-full rounded-2xl border border-line bg-slate-50 px-4 py-3 text-sm font-semibold text-slate-700"
      >
        Sign out
      </button>
    </form>
  </aside>
</div>
