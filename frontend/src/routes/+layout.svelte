<script lang="ts">
  import '../app.css';
  import { page } from '$app/stores';
  import AppShell from '$lib/components/AppShell.svelte';
  import { authUser } from '$lib/stores/auth';

  export let data: App.PageData;

  $: authUser.set(data.user);
  $: isPublicRoute = $page.url.pathname === '/login';
</script>

{#if isPublicRoute}
  <slot />
{:else}
  <AppShell user={data.user} currentPath={$page.url.pathname}>
    <slot />
  </AppShell>
{/if}
