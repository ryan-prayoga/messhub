<script lang="ts">
  export let statCards = 3;
  export let rows = 3;
</script>

<div class="page-skeleton">
  <div class="page-skeleton-hero">
    <div class="skeleton-line h-4 w-24"></div>
    <div class="skeleton-line h-8 w-56"></div>
    <div class="skeleton-line h-4 w-full max-w-xl"></div>
  </div>

  <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
    {#each Array(statCards) as _, index}
      <div class={`skeleton-card ${index === 0 ? 'xl:col-span-1' : ''}`}>
        <div class="skeleton-line h-3 w-20"></div>
        <div class="skeleton-line mt-4 h-8 w-28"></div>
        <div class="skeleton-line mt-4 h-4 w-full"></div>
      </div>
    {/each}
  </div>

  <div class="space-y-3">
    {#each Array(rows) as _}
      <div class="skeleton-card">
        <div class="flex items-center justify-between gap-4">
          <div class="w-full">
            <div class="skeleton-line h-4 w-32"></div>
            <div class="skeleton-line mt-3 h-4 w-full"></div>
          </div>
          <div class="skeleton-line h-10 w-20"></div>
        </div>
      </div>
    {/each}
  </div>
</div>
