<script lang="ts">
  import { navigating, page } from '$app/stores';
  import { fly, fade } from 'svelte/transition';
  import { APP_NAME } from '$lib/config/env';
  import AppIcon from '$lib/components/AppIcon.svelte';
  import PwaControlBar from '$lib/components/PwaControlBar.svelte';
  import { clearOfflineSessionArtifacts } from '$lib/pwa/runtime';
  import { getCurrentNavigationItem, getPageMeta, getVisibleNavigation, isPathActive } from '$lib/navigation';

  export let user: App.Locals['user'];
  export let currentPath: string;
  export let notificationSummary: {
    unread_count: number;
  };

  const roleLabels: Record<string, string> = {
    admin: 'Admin',
    treasurer: 'Bendahara',
    member: 'Anggota'
  };

  let mobileMenuOpen = false;
  let userMenuOpen = false;

  $: actualPath = $page.url.pathname || currentPath;
  $: displayPath = $navigating?.to?.url.pathname ?? actualPath;
  $: navGroups = getVisibleNavigation(user);
  $: currentItem = getCurrentNavigationItem(displayPath, user);
  $: currentMeta = getPageMeta(displayPath, user);
  $: unreadCount = notificationSummary?.unread_count ?? 0;
  $: if (actualPath) {
    mobileMenuOpen = false;
    userMenuOpen = false;
  }

  function toggleMobileMenu() {
    mobileMenuOpen = !mobileMenuOpen;
    if (mobileMenuOpen) {
      userMenuOpen = false;
    }
  }

  function toggleUserMenu() {
    userMenuOpen = !userMenuOpen;
    if (userMenuOpen) {
      mobileMenuOpen = false;
    }
  }

  function handleSignOut() {
    void clearOfflineSessionArtifacts();
  }

  function closeMenus() {
    mobileMenuOpen = false;
    userMenuOpen = false;
  }

  function handleWindowKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      closeMenus();
    }
  }
</script>

<svelte:window on:keydown={handleWindowKeydown} />

<div class="app-shell">
  <div class="shell-layout">
    <aside class="shell-sidebar">
      <a href="/dashboard" class="shell-brand" data-sveltekit-preload-data="tap" on:click={closeMenus}>
        <div class="shell-brand-mark">M</div>
        <div class="min-w-0">
          <p class="shell-brand-kicker">Mess Traspac Menyala</p>
          <h1 class="shell-brand-title">{APP_NAME}</h1>
          <p class="shell-brand-copy">Operasional harian yang rapi, ringan, dan nyaman dibuka dari HP maupun desktop.</p>
        </div>
      </a>

      <section class="shell-sidebar-card">
        <p class="helper-label">Sedang dibuka</p>
        <h2 class="mt-3 text-xl font-semibold text-ink">{currentItem?.label ?? currentMeta.title}</h2>
        <p class="mt-2 text-sm leading-6 text-muted">{currentMeta.description}</p>
      </section>

      <nav class="nav-section">
        <p class="nav-section-label">Utama</p>
        <div class="nav-stack">
          {#each navGroups.primary as item}
            <a
              href={item.href}
              class={`nav-link ${isPathActive(displayPath, item.href) ? 'nav-link-active' : ''}`}
              data-sveltekit-preload-data="tap"
              on:click={closeMenus}
            >
              <span class="nav-link-icon">
                <AppIcon icon={item.icon} className="h-5 w-5" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="nav-link-label">{item.label}</span>
                <span class="nav-link-copy">{item.description}</span>
              </span>
              <AppIcon icon="lucide:chevron-right" className="nav-link-chevron" />
            </a>
          {/each}
        </div>
      </nav>

      <nav class="nav-section">
        <p class="nav-section-label">Ruang Kerja</p>
        <div class="nav-stack">
          {#each [...navGroups.workspace, ...navGroups.admin] as item}
            <a
              href={item.href}
              class={`nav-link ${isPathActive(displayPath, item.href) ? 'nav-link-active' : ''}`}
              data-sveltekit-preload-data="tap"
              on:click={closeMenus}
            >
              <span class="nav-link-icon">
                <AppIcon icon={item.icon} className="h-5 w-5" />
              </span>
              <span class="min-w-0 flex-1">
                <span class="nav-link-label">{item.label}</span>
                <span class="nav-link-copy">{item.description}</span>
              </span>
              {#if item.href === '/notifications' && unreadCount > 0}
                <span class="nav-badge">{unreadCount > 9 ? '9+' : unreadCount}</span>
              {:else}
                <AppIcon icon="lucide:chevron-right" className="nav-link-chevron" />
              {/if}
            </a>
          {/each}
        </div>
      </nav>

      {#if user}
        <section class="shell-sidebar-card shell-user-card">
          <div class="avatar-chip">{user.name.slice(0, 1)}</div>
          <div class="min-w-0">
            <p class="text-sm font-semibold text-ink">{user.name}</p>
            <p class="mt-1 truncate text-sm text-muted">@{user.username}</p>
            <p class="mt-1 text-xs uppercase tracking-[0.16em] text-dusty">{roleLabels[user.role] ?? user.role}</p>
          </div>
        </section>
      {/if}
    </aside>

    <div class="shell-content">
      <header class="shell-header">
        <div class="shell-header-intro">
          <div class="lg:hidden">
            <p class="shell-mobile-kicker">{APP_NAME}</p>
            <h1 class="shell-mobile-title">{currentItem?.label ?? currentMeta.title}</h1>
          </div>

          <div class="hidden lg:block">
            <p class="eyebrow">Operasional Mess</p>
            <h1 class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{currentMeta.title}</h1>
            <p class="mt-2 max-w-2xl text-sm leading-6 text-muted">{currentMeta.description}</p>
          </div>
        </div>

        <div class="shell-header-actions">
          <a
            href="/notifications"
            class={`icon-button ${isPathActive(displayPath, '/notifications') ? 'icon-button-active' : ''}`}
            aria-label="Buka notifikasi"
            data-sveltekit-preload-data="tap"
            on:click={closeMenus}
          >
            <AppIcon icon="lucide:bell-ring" className="h-5 w-5" />
            {#if unreadCount > 0}
              <span class="notification-dot">{unreadCount > 9 ? '9+' : unreadCount}</span>
            {/if}
          </a>

          <button
            type="button"
            class="icon-button lg:hidden"
            aria-expanded={mobileMenuOpen}
            aria-label="Buka menu"
            on:click={toggleMobileMenu}
          >
            <AppIcon icon={mobileMenuOpen ? 'lucide:x' : 'lucide:menu'} className="h-5 w-5" />
          </button>

          {#if user}
            <div class="relative">
              <button
                type="button"
                class={`profile-trigger ${userMenuOpen ? 'profile-trigger-active' : ''}`}
                aria-expanded={userMenuOpen}
                aria-label="Buka menu akun"
                on:click={toggleUserMenu}
              >
                <div class="avatar-chip avatar-chip-sm">{user.name.slice(0, 1)}</div>
                <div class="hidden min-w-0 text-left sm:block">
                  <p class="truncate text-sm font-semibold text-ink">{user.name}</p>
                  <p class="truncate text-xs text-muted">@{user.username}</p>
                </div>
                <AppIcon icon="lucide:chevrons-up-down" className="hidden h-4 w-4 text-muted sm:block" />
              </button>

              {#if userMenuOpen}
                <div class="menu-popover" transition:fly={{ y: -6, duration: 160 }}>
                  <div class="menu-profile">
                    <div class="avatar-chip">{user.name.slice(0, 1)}</div>
                    <div class="min-w-0">
                      <p class="truncate text-sm font-semibold text-ink">{user.name}</p>
                      <p class="mt-1 truncate text-sm text-muted">{user.email}</p>
                      <p class="mt-1 text-xs uppercase tracking-[0.16em] text-dusty">
                        {roleLabels[user.role] ?? user.role}
                      </p>
                    </div>
                  </div>

                  <a href="/profile" class="menu-link" data-sveltekit-preload-data="tap" on:click={closeMenus}>
                    <AppIcon icon="lucide:user-round" className="h-4 w-4" />
                    <span>Profile</span>
                  </a>

                  {#if user.role === 'admin'}
                    <a href="/settings" class="menu-link" data-sveltekit-preload-data="tap" on:click={closeMenus}>
                      <AppIcon icon="lucide:settings-2" className="h-4 w-4" />
                      <span>Settings</span>
                    </a>
                  {/if}

                  <form method="POST" action="/logout" on:submit={handleSignOut}>
                    <button type="submit" class="menu-link menu-link-danger w-full">
                      <AppIcon icon="lucide:log-out" className="h-4 w-4" />
                      <span>Keluar</span>
                    </button>
                  </form>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </header>

      {#if mobileMenuOpen}
        <section class="mobile-menu" transition:fade={{ duration: 160 }}>
          <div class="mobile-menu-grid">
            {#each [...navGroups.primary, ...navGroups.workspace, ...navGroups.admin] as item}
              <a
                href={item.href}
                class={`mobile-menu-link ${isPathActive(displayPath, item.href) ? 'mobile-menu-link-active' : ''}`}
                data-sveltekit-preload-data="tap"
                on:click={closeMenus}
              >
                <AppIcon icon={item.icon} className="h-5 w-5" />
                <div class="min-w-0">
                  <p class="text-sm font-semibold">{item.label}</p>
                  <p class="mt-1 text-xs leading-5 text-muted">{item.description}</p>
                </div>
              </a>
            {/each}
          </div>
        </section>
      {/if}

      <PwaControlBar />

      <main class="page-container">
        <div class="page-stack">
          {#key currentPath}
            <div class="page-transition-frame" transition:fly={{ y: 12, duration: 180 }}>
              <slot />
            </div>
          {/key}
        </div>
      </main>

      <nav class="bottom-nav-bar lg:hidden">
        <div class="bottom-nav-grid">
          {#each navGroups.bottom as item}
            <a
              href={item.href}
              class={`bottom-nav-link ${isPathActive(displayPath, item.href) ? 'bottom-nav-link-active' : ''}`}
              data-sveltekit-preload-data="tap"
              on:click={closeMenus}
            >
              <span class="bottom-nav-icon">
                <AppIcon icon={item.icon} className="h-5 w-5" />
              </span>
              <span>{item.label}</span>
            </a>
          {/each}
        </div>
      </nav>
    </div>
  </div>
</div>
