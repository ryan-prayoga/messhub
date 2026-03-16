<script lang="ts">
  import '../app.css';
  import { page } from '$app/stores';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';

  export let error: App.Error & { requestId?: string };
  export let status: number;
</script>

<div class="app-shell">
  <main class="page-container">
    <PageCard
      title="Something went wrong"
      description="UI tetap dijaga agar tidak blank saat backend, load function, atau route action gagal."
    >
      <StatePanel
        tone={status === 403 ? 'forbidden' : 'error'}
        title={status === 403 ? 'Forbidden' : `Error ${status}`}
        message={error?.message ?? 'Unexpected application failure'}
        requestId={error?.requestId ?? null}
      />

      <div class="mt-4 flex flex-wrap gap-3">
        <a href="/dashboard" class="btn-primary px-4 py-3">Back to dashboard</a>
        <a href={$page.url.pathname} class="btn-secondary px-4 py-3">Try again</a>
      </div>
    </PageCard>
  </main>
</div>
