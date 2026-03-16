<script lang="ts">
  import { onMount } from 'svelte';
  import '../app.css';
  import { page } from '$app/stores';
  import AppShell from '$lib/components/AppShell.svelte';
  import { SHELLLESS_ROUTES } from '$lib/auth/session';
  import { initializePwaRuntime } from '$lib/pwa/runtime';
  import { authState } from '$lib/stores/auth';

  export let data: App.PageData;

  $: authState.sync(data.user);
  $: isPublicRoute = SHELLLESS_ROUTES.includes($page.url.pathname);

  onMount(() => initializePwaRuntime());
</script>

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
