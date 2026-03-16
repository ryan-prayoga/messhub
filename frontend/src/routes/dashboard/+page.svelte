<script lang="ts">
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
</script>

<div class="space-y-4">
  <PageCard
    title={`Halo, ${data.user?.name ?? 'Member'}`}
    description="STEP 1 aktif: auth, session, endpoint me, role guard, dan fondasi user management sudah tersambung."
  >
    <div class="grid gap-3 md:grid-cols-3">
      <div class="stat-card">
        <p class="helper-label">Signed in as</p>
        <p class="mt-2 text-lg font-semibold text-ink">{data.user?.name}</p>
        <p class="mt-1 text-sm text-slate-500">{data.user?.email}</p>
        {#if data.user}
          <div class="mt-3">
            <span class={roleBadgeClass(data.user.role)}>{roleLabels[data.user.role]}</span>
          </div>
        {/if}
      </div>

      <div class="stat-card bg-white">
        <p class="helper-label">Member count</p>
        {#if data.memberSummary.state === 'ready'}
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.memberSummary.total}</p>
          <p class="mt-2 text-sm text-slate-500">
            {data.memberSummary.active} aktif, {data.memberSummary.inactive} nonaktif
          </p>
        {:else if data.memberSummary.state === 'restricted'}
          <p class="mt-2 text-lg font-semibold text-ink">Restricted</p>
          <p class="mt-2 text-sm text-slate-500">
            Ringkasan anggota tersedia untuk admin dan treasurer.
          </p>
        {:else}
          <p class="mt-2 text-lg font-semibold text-ink">Unavailable</p>
          <p class="mt-2 text-sm text-slate-500">{data.memberSummary.message}</p>
        {/if}
      </div>

      <div class="stat-card bg-sky-50/80">
        <p class="helper-label text-sky-700">Auth status</p>
        <p class="mt-2 text-lg font-semibold text-ink capitalize">{data.authStatus}</p>
        <p class="mt-2 text-sm text-slate-600">
          Session diverifikasi lewat token cookie dan endpoint <code>/api/v1/auth/me</code>.
        </p>
      </div>
    </div>

    <div class="mt-4 grid gap-3 sm:grid-cols-3">
      <a href="/members" class="stat-card bg-white transition hover:border-sky-300 hover:bg-sky-50/50">
        <p class="helper-label">Menu cepat</p>
        <p class="mt-2 text-base font-semibold text-ink">Members</p>
        <p class="mt-2 text-sm leading-6 text-slate-500">
          Lihat daftar anggota, role, dan status aktif mess.
        </p>
      </a>

      <a href="/wallet" class="stat-card bg-white transition hover:border-slate-300 hover:bg-slate-50">
        <p class="helper-label">Placeholder</p>
        <p class="mt-2 text-base font-semibold text-ink">Wallet</p>
        <p class="mt-2 text-sm leading-6 text-slate-500">
          Fondasi auth sudah siap untuk modul kas dan transaksi.
        </p>
      </a>

      <a href="/wifi" class="stat-card bg-white transition hover:border-slate-300 hover:bg-slate-50">
        <p class="helper-label">Placeholder</p>
        <p class="mt-2 text-base font-semibold text-ink">Wifi</p>
        <p class="mt-2 text-sm leading-6 text-slate-500">
          Billing wifi bulanan bisa dibangun di atas session dan role guard ini.
        </p>
      </a>
    </div>
  </PageCard>

  <PageCard
    title="Foundation Scope"
    description="Hal yang sudah aktif di STEP 1 tanpa membuka scope wallet, wifi, contribution, atau feed lebih jauh."
  >
    <div class="grid gap-3 md:grid-cols-2">
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Backend</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">
          Login, me, auth middleware, role middleware, users schema, seed admin, dan member API dasar.
        </p>
      </div>

      <div class="helper-box">
        <p class="helper-label">Frontend</p>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Login end-to-end, auth state sederhana, route protection, dashboard awal, dan members list mobile-first.
        </p>
      </div>
    </div>
  </PageCard>
</div>
