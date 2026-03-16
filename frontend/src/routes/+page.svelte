<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';

  export let data: App.PageData;
</script>

<div class="space-y-4">
  <PageCard
    title="Dashboard"
    description="Ringkasan awal untuk wallet, wifi, shared expenses, contributions, dan feed aktif."
  >
    <div class="grid gap-3 sm:grid-cols-2">
      <div class="rounded-2xl bg-slate-100 p-4">
        <p class="text-xs uppercase tracking-[0.2em] text-slate-500">Signed in</p>
        <p class="mt-2 text-lg font-bold">{data.user?.name}</p>
        <p class="text-sm text-slate-500">{data.user?.email}</p>
      </div>
      <div class="rounded-2xl bg-orange-50 p-4">
        <p class="text-xs uppercase tracking-[0.2em] text-orange-600">Phase</p>
        <p class="mt-2 text-lg font-bold text-orange-900">Foundation</p>
        <p class="text-sm text-orange-700">Core modules are scaffolded and ready for feature work.</p>
      </div>
    </div>
  </PageCard>
</div>
