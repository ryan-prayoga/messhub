<script lang="ts">
  import { navigating } from '$app/stores';
  import PageCard from '$lib/components/PageCard.svelte';
  import type { PageData } from './$types';
  import type { UserRole } from '$lib/api/types';

  export let data: PageData;

  const roleLabels: Record<UserRole, string> = {
    admin: 'Admin',
    treasurer: 'Treasurer',
    member: 'Member'
  };

  function roleBadgeClass(role: UserRole) {
    if (role === 'admin') {
      return 'badge bg-slate-950 text-white';
    }

    if (role === 'treasurer') {
      return 'badge bg-emerald-100 text-emerald-700';
    }

    return 'badge-muted';
  }

  function statusBadgeClass(isActive: boolean) {
    return isActive ? 'badge-brand' : 'badge bg-rose-100 text-rose-700';
  }

  function formatDate(value: string | null) {
    if (!value) {
      return '-';
    }

    return new Intl.DateTimeFormat('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    }).format(new Date(value));
  }
</script>

<div class="space-y-4">
  <PageCard
    title="Members"
    description="Daftar anggota mess untuk STEP 1. Endpoint backend yang dipakai: `GET /api/v1/users`."
  >
    {#if $navigating?.to?.url.pathname === '/members'}
      <div class="helper-box mb-4">
        <p class="helper-label">Loading</p>
        <p class="mt-2 text-sm text-slate-600">Memuat ulang data anggota mess...</p>
      </div>
    {/if}

    {#if data.accessDenied}
      <div class="empty-state">
        Role <strong>{data.user?.role}</strong> belum diizinkan melihat daftar anggota. Akses page ini
        tersedia untuk admin dan treasurer.
      </div>
    {:else if data.loadError}
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Error</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">{data.loadError}</p>
      </div>
    {:else}
      <div class="grid gap-3 sm:grid-cols-3">
        <div class="stat-card">
          <p class="helper-label">Total</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.summary.total}</p>
        </div>

        <div class="stat-card bg-white">
          <p class="helper-label">Aktif</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.summary.active}</p>
        </div>

        <div class="stat-card bg-white">
          <p class="helper-label">Nonaktif</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.summary.inactive}</p>
        </div>
      </div>

      {#if data.members.length === 0}
        <div class="mt-4 empty-state">
          Belum ada anggota yang tampil dari backend. Setelah admin menambahkan user baru, daftar ini akan
          langsung terisi.
        </div>
      {:else}
        <div class="mt-4 space-y-3">
          {#each data.members as member}
            <article class="stat-card bg-white">
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div class="min-w-0">
                  <h3 class="text-base font-semibold text-ink">{member.name}</h3>
                  <p class="mt-1 break-all text-sm text-slate-500">{member.email}</p>
                </div>

                <div class="flex flex-wrap gap-2">
                  <span class={roleBadgeClass(member.role)}>{roleLabels[member.role]}</span>
                  <span class={statusBadgeClass(member.is_active)}>
                    {member.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
              </div>

              <div class="mt-4 grid gap-2 sm:grid-cols-2">
                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Joined</p>
                  <p class="mt-2 text-sm font-medium text-ink">{formatDate(member.joined_at)}</p>
                </div>

                <div class="rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                  <p class="helper-label">Left</p>
                  <p class="mt-2 text-sm font-medium text-ink">{formatDate(member.left_at)}</p>
                </div>
              </div>
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
