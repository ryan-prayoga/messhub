<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';
  import type { ContributionLeaderboardEntry } from '$lib/api/types';
  import type { PageData } from './$types';

  export let data: PageData;

  let boards: { title: string; description: string; items: ContributionLeaderboardEntry[] }[] = [];

  $: boards = [
    {
      title: 'This Month',
      description: 'Skor kontribusi bulan berjalan.',
      items: data.monthly
    },
    {
      title: 'All Time',
      description: 'Skor kontribusi keseluruhan.',
      items: data.allTime
    }
  ];
</script>

<div class="space-y-4">
  <PageCard
    title="Contributions"
    description="Leaderboard kontribusi yang dihitung dari activity type contribution."
  >
    {#if data.loadError}
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{data.loadError}</p>
      </div>
    {:else}
      <div class="grid gap-4 lg:grid-cols-2">
        {#each boards as board}
          <article class="app-panel p-5">
            <p class="eyebrow">{board.title}</p>
            <h2 class="section-title mt-1">{board.title} leaderboard</h2>
            <p class="section-subtitle mt-2">{board.description}</p>

            {#if board.items.length === 0}
              <div class="empty-state mt-4">Belum ada data leaderboard.</div>
            {:else}
              <div class="mt-4 space-y-3">
                {#each board.items as item}
                  <div class="flex items-center justify-between rounded-3xl border border-slate-200 bg-slate-50 px-4 py-3">
                    <div>
                      <p class="text-sm font-semibold text-ink">#{item.rank} {item.user_name}</p>
                      <p class="mt-1 text-xs text-slate-500">{item.total_activities} aktivitas</p>
                    </div>

                    <div class="text-right">
                      <p class="text-xl font-semibold tracking-[-0.04em] text-ink">{item.total_points}</p>
                      <p class="text-xs uppercase tracking-[0.16em] text-slate-400">points</p>
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </article>
        {/each}
      </div>
    {/if}
  </PageCard>
</div>
