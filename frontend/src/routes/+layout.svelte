<script lang="ts">
  import { onMount } from 'svelte';
  import '../app.css';
  import { page } from '$app/stores';
  import AppShell from '$lib/components/AppShell.svelte';
  import { SHELLLESS_ROUTES } from '$lib/auth/session';
  import { buildPageTitle, getPageMeta } from '$lib/navigation';
  import { initializePwaRuntime } from '$lib/pwa/runtime';
  import { authState } from '$lib/stores/auth';

  export let data: App.PageData;

  $: authState.sync(data.user);
  $: isPublicRoute = SHELLLESS_ROUTES.includes($page.url.pathname);
  $: currentMeta = getPageMeta($page.url.pathname, data.user);

  onMount(() => initializePwaRuntime());
</script>

<svelte:head>
  <title>{buildPageTitle(currentMeta.title)}</title>
  <meta name="description" content={currentMeta.description} />
</svelte:head>

{#if isPublicRoute}
  <slot />
{:else}
  <AppShell
    user={data.user}
    currentPath={$page.url.pathname}
    notificationSummary={data.notificationSummary ?? { unread_count: 0 }}
  >
    <slot />
  </AppShell>
{/if}
