<script lang="ts">
  import { APP_NAME } from '$lib/config/env';
  import PwaControlBar from '$lib/components/PwaControlBar.svelte';
  import { clearOfflineSessionArtifacts } from '$lib/pwa/runtime';

  type NavItem = {
    href: string;
    label: string;
    icon:
      | 'dashboard'
      | 'feed'
      | 'wallet'
      | 'wifi'
      | 'profile'
      | 'members'
      | 'contrib'
      | 'notifications'
      | 'shared'
      | 'proposals'
      | 'settings';
  };

  export let user: App.Locals['user'];
  export let currentPath: string;
  export let notificationSummary: {
    unread_count: number;
  };

  const navItems: NavItem[] = [
    { href: '/dashboard', label: 'Dashboard', icon: 'dashboard' },
    { href: '/feed', label: 'Feed', icon: 'feed' },
    { href: '/wallet', label: 'Kas', icon: 'wallet' },
    { href: '/wifi', label: 'Wifi', icon: 'wifi' },
    { href: '/profile', label: 'Profil', icon: 'profile' }
  ];

  const utilityBaseItems: NavItem[] = [
    { href: '/members', label: 'Anggota', icon: 'members' },
    { href: '/notifications', label: 'Inbox', icon: 'notifications' },
    { href: '/contributions', label: 'Kontribusi', icon: 'contrib' },
    { href: '/shared-expenses', label: 'Patungan', icon: 'shared' },
    { href: '/proposals', label: 'Usulan', icon: 'proposals' }
  ];

  let utilityItems: NavItem[] = utilityBaseItems;
  const roleLabels: Record<string, string> = {
    admin: 'Admin',
    treasurer: 'Bendahara',
    member: 'Anggota'
  };

  const isCurrentPath = (href: string) =>
    currentPath === href || (href !== '/' && currentPath.startsWith(`${href}/`));

  $: utilityItems = user?.role === 'admin'
    ? [
        ...utilityBaseItems,
        { href: '/admin/import', label: 'Impor', icon: 'settings' },
        { href: '/settings', label: 'Pengaturan', icon: 'settings' }
      ]
    : utilityBaseItems;
  $: currentItem = [...navItems, ...utilityItems].find((item) =>
    isCurrentPath(item.href)
  );
  $: unreadCount = notificationSummary?.unread_count ?? 0;

  function handleSignOut() {
    void clearOfflineSessionArtifacts();
  }
</script>

<div class="app-shell">
  <div class="mx-auto flex min-h-screen w-full max-w-4xl flex-col">
    <header class="sticky top-0 z-20 border-b border-line/80 bg-white/90 backdrop-blur">
      <div class="mx-auto flex w-full max-w-4xl flex-col gap-4 px-4 pb-4 pt-[calc(1rem+env(safe-area-inset-top))] sm:px-6">
        <div class="flex items-start justify-between gap-4">
          <div class="min-w-0">
            <p class="eyebrow">Operasional Mess</p>
            <h1 class="mt-1 text-2xl font-semibold tracking-[-0.03em] text-ink">{APP_NAME}</h1>
            <p class="mt-1 text-sm text-slate-500">
              {currentItem?.label || 'Beranda'} untuk operasional harian mess.
            </p>
          </div>

          <div class="flex shrink-0 items-start gap-2">
            <a
              href="/notifications"
              class={`relative inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-slate-200 bg-white text-slate-600 transition hover:border-slate-300 hover:text-slate-950 ${
                isCurrentPath('/notifications') ? 'border-sky-200 bg-sky-50 text-sky-700' : ''
              }`}
              aria-label="Notifications"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                class="h-5 w-5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M14.857 17H9.143M18 17V11.5a6 6 0 0 0-12 0V17l-1.2 1.6A1 1 0 0 0 5.6 20h12.8a1 1 0 0 0 .8-1.4L18 17ZM13.73 20a2 2 0 0 1-3.46 0"
                />
              </svg>

              {#if unreadCount > 0}
                <span class="absolute -right-1 -top-1 inline-flex min-w-[1.35rem] items-center justify-center rounded-full bg-rose-500 px-1.5 py-0.5 text-[10px] font-semibold text-white">
                  {unreadCount > 9 ? '9+' : unreadCount}
                </span>
              {/if}
            </a>

            {#if user}
              <div class="hidden rounded-2xl border border-line bg-slate-50 px-3 py-2 text-right text-xs sm:block">
                <p class="font-semibold text-ink">{user.name}</p>
                <p class="text-slate-500">{roleLabels[user.role] ?? user.role}</p>
              </div>
            {/if}

            <form method="POST" action="/logout" on:submit={handleSignOut}>
              <button type="submit" class="btn-secondary px-3 py-2 text-xs">Keluar</button>
            </form>
          </div>
        </div>

        {#if user}
          <div class="sm:hidden">
            <div class="inline-flex rounded-full border border-line bg-slate-50 px-3 py-1.5 text-xs text-slate-600">
              <span class="font-semibold text-ink">{user.name}</span>
              <span class="mx-2 text-slate-300">/</span>
              <span>{roleLabels[user.role] ?? user.role}</span>
            </div>
          </div>
        {/if}

        <div class="flex gap-2 overflow-x-auto pb-1">
          {#each utilityItems as item}
            <a
              href={item.href}
              class={`nav-chip ${isCurrentPath(item.href) ? 'nav-chip-active' : ''}`}
            >
              {item.label}
            </a>
          {/each}
        </div>
      </div>
    </header>

    <PwaControlBar />

    <main class="page-container flex-1">
      <slot />
    </main>

    <nav class="sticky bottom-0 z-20 border-t border-line/80 bg-white/95 px-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] pt-3 backdrop-blur">
      <div class="mx-auto grid w-full max-w-4xl grid-cols-5 gap-2">
        {#each navItems as item}
          <a
            href={item.href}
            class={`bottom-nav-link ${
              isCurrentPath(item.href) ? 'bottom-nav-link-active' : 'bg-white hover:bg-slate-50'
            }`}
          >
            <span class="bottom-nav-icon" aria-hidden="true">
              {#if item.icon === 'dashboard'}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 11.5 12 5l8 6.5V20a1 1 0 0 1-1 1h-4.5v-5.5h-5V21H5a1 1 0 0 1-1-1v-8.5Z" />
                </svg>
              {:else if item.icon === 'feed'}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 6.5h12M6 12h12M6 17.5h7" />
                </svg>
              {:else if item.icon === 'wallet'}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 8.5A2.5 2.5 0 0 1 6.5 6H18a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6.5A2.5 2.5 0 0 1 4 15.5v-7Z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M15 12h5M17 12h.01" />
                </svg>
              {:else if item.icon === 'wifi'}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 9.5a11 11 0 0 1 15 0M7.5 12.5a7 7 0 0 1 9 0M10.5 15.5a3 3 0 0 1 3 0" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 19.5h.01" />
                </svg>
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8ZM5 20a7 7 0 0 1 14 0" />
                </svg>
              {/if}
            </span>
            <span>{item.label}</span>
          </a>
        {/each}
      </div>
    </nav>
  </div>
</div>
