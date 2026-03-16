<script lang="ts">
  import { APP_NAME } from '$lib/config/env';

  type NavItem = {
    href: string;
    label: string;
  };

  export let user: App.Locals['user'];
  export let currentPath: string;
  export let notificationSummary: {
    unread_count: number;
  };

  const navItems: NavItem[] = [
    { href: '/dashboard', label: 'Dashboard' },
    { href: '/members', label: 'Members' },
    { href: '/wallet', label: 'Wallet' },
    { href: '/wifi', label: 'Wifi' },
    { href: '/feed', label: 'Feed' }
  ];

  const utilityBaseItems: NavItem[] = [
    { href: '/shared-expenses', label: 'Shared' },
    { href: '/contributions', label: 'Contrib' },
    { href: '/proposals', label: 'Proposals' },
    { href: '/profile', label: 'Profile' }
  ];

  const shellOnlyItems: NavItem[] = [{ href: '/notifications', label: 'Notifications' }];
  let utilityItems: NavItem[] = utilityBaseItems;

  const isCurrentPath = (href: string) =>
    currentPath === href || (href !== '/' && currentPath.startsWith(`${href}/`));

  $: utilityItems = user?.role === 'admin'
    ? [...utilityBaseItems, { href: '/settings', label: 'Settings' }]
    : utilityBaseItems;
  $: currentItem = [...navItems, ...utilityItems, ...shellOnlyItems].find((item) =>
    isCurrentPath(item.href)
  );
  $: unreadCount = notificationSummary?.unread_count ?? 0;
</script>

<div class="app-shell">
  <div class="mx-auto flex min-h-screen w-full max-w-4xl flex-col">
    <header class="sticky top-0 z-20 border-b border-line/80 bg-white/90 backdrop-blur">
      <div class="mx-auto flex w-full max-w-4xl flex-col gap-4 px-4 pb-4 pt-4 sm:px-6">
        <div class="flex items-start justify-between gap-4">
          <div class="min-w-0">
            <p class="eyebrow">Internal Mess App</p>
            <h1 class="mt-1 text-2xl font-semibold tracking-[-0.03em] text-ink">{APP_NAME}</h1>
            <p class="mt-1 text-sm text-slate-500">
              {currentItem?.label || 'Workspace'} untuk operasional harian mess.
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
                <p class="text-slate-500">{user.role}</p>
              </div>
            {/if}

            <form method="POST" action="/logout">
              <button type="submit" class="btn-secondary px-3 py-2 text-xs">Sign out</button>
            </form>
          </div>
        </div>

        {#if user}
          <div class="sm:hidden">
            <div class="inline-flex rounded-full border border-line bg-slate-50 px-3 py-1.5 text-xs text-slate-600">
              <span class="font-semibold text-ink">{user.name}</span>
              <span class="mx-2 text-slate-300">/</span>
              <span>{user.role}</span>
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
            {item.label}
          </a>
        {/each}
      </div>
    </nav>
  </div>
</div>
