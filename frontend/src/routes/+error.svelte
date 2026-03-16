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
      eyebrow="Kendala"
      icon="lucide:shield-alert"
      title="Terjadi Kendala"
      description="Halaman tetap dijaga agar tidak kosong saat ada gangguan pemuatan atau proses data."
    >
      <StatePanel
        tone={status === 403 ? 'forbidden' : 'error'}
        title={status === 403 ? 'Akses ditolak' : `Status ${status}`}
        message={error?.message ?? 'Terjadi gangguan yang belum dapat diproses.'}
        requestId={error?.requestId ?? null}
      />

      <div class="mt-4 flex flex-wrap gap-3">
        <a href="/dashboard" class="btn-primary px-4 py-3">Kembali ke dashboard</a>
        <a href={$page.url.pathname} class="btn-secondary px-4 py-3">Coba lagi</a>
      </div>
    </PageCard>
  </main>
</div>
