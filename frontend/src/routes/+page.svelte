<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';

  export let data: App.PageData;
</script>

<div class="space-y-4">
  <PageCard
    title="Dashboard"
    description="Ringkasan awal untuk wallet, wifi, shared expenses, contributions, dan feed aktif."
  >
    <div class="grid gap-3 md:grid-cols-2">
      <div class="stat-card">
        <p class="helper-label">Signed in</p>
        <p class="mt-2 text-lg font-semibold text-ink">{data.user?.name}</p>
        <p class="mt-1 text-sm text-slate-500">{data.user?.email}</p>
        <div class="mt-3">
          <span class="badge-brand">{data.user?.role}</span>
        </div>
      </div>

      <div class="stat-card bg-accentSoft/70">
        <p class="helper-label text-accentStrong">Phase</p>
        <p class="mt-2 text-lg font-semibold text-ink">Foundation</p>
        <p class="mt-1 text-sm text-slate-600">
          Shell, login, and global UI baseline are ready for dashboard modules.
        </p>
      </div>
    </div>

    <div class="mt-4 empty-state">
      Dashboard masih placeholder. Langkah berikutnya adalah mengisi summary cards, daftar aktivitas terbaru, dan empty/loading states per modul tanpa perlu mengulang styling dasar lagi.
    </div>
  </PageCard>
</div>
