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

  function formatCurrency(value: number | null) {
    if (value === null) {
      return '-';
    }

    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(value);
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

  function wifiStatusLabel(status: string | null) {
    if (status === 'pending_verification') {
      return 'Pending';
    }

    if (status === 'verified') {
      return 'Verified';
    }

    if (status === 'rejected') {
      return 'Rejected';
    }

    if (status === 'unpaid') {
      return 'Unpaid';
    }

    return 'No active bill';
  }

  function wifiStatusClass(status: string | null) {
    if (status === 'verified') {
      return 'badge bg-emerald-100 text-emerald-700';
    }

    if (status === 'pending_verification') {
      return 'badge bg-amber-100 text-amber-700';
    }

    if (status === 'rejected') {
      return 'badge bg-rose-100 text-rose-700';
    }

    if (status === 'unpaid') {
      return 'badge-muted';
    }

    return 'badge bg-slate-100 text-slate-600';
  }
</script>

<div class="space-y-4">
  <PageCard
    title={`Halo, ${data.user?.name ?? 'Member'}`}
    description="Dashboard sekarang menampilkan fondasi auth, members, wallet, dan wifi billing untuk STEP 3."
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

      <a href="/wallet" class="stat-card bg-white transition hover:border-sky-300 hover:bg-sky-50/50">
        <p class="helper-label">Wallet</p>
        <p class="mt-2 text-base font-semibold text-ink">Wallet</p>
        {#if data.walletSummary.state === 'ready'}
          <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
            {formatCurrency(data.walletSummary.balance)}
          </p>
          <p class="mt-2 text-sm leading-6 text-slate-500">
            Income {formatCurrency(data.walletSummary.totalIncome)} / Expense {formatCurrency(data.walletSummary.totalExpense)}
          </p>
        {:else}
          <p class="mt-2 text-sm leading-6 text-slate-500">{data.walletSummary.message}</p>
        {/if}
      </a>

      <a href="/wifi" class="stat-card bg-white transition hover:border-slate-300 hover:bg-slate-50">
        <p class="helper-label">Wifi</p>
        <p class="mt-2 text-base font-semibold text-ink">Wifi</p>
        {#if data.wifiSummary.state === 'ready'}
          <p class="mt-2 text-lg font-semibold text-ink">{data.wifiSummary.monthLabel}</p>
          {#if data.user?.role === 'member'}
            <div class="mt-2">
              <span class={wifiStatusClass(data.wifiSummary.myStatus)}>
                {wifiStatusLabel(data.wifiSummary.myStatus)}
              </span>
            </div>
            <p class="mt-2 text-sm leading-6 text-slate-500">
              Deadline {formatDate(data.wifiSummary.deadline)}.
            </p>
          {:else}
            <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
              {data.wifiSummary.verified} verified
            </p>
            <p class="mt-2 text-sm leading-6 text-slate-500">
              {data.wifiSummary.unpaid} unpaid / {data.wifiSummary.pending} pending
            </p>
          {/if}
        {:else}
          <p class="mt-2 text-sm leading-6 text-slate-500">{data.wifiSummary.message}</p>
        {/if}
      </a>
    </div>
  </PageCard>

  <PageCard
    title="Foundation Scope"
    description="Fondasi yang sudah aktif setelah STEP 2 wallet dasar selesai ditambahkan."
  >
    <div class="grid gap-3 md:grid-cols-2">
      <div class="helper-box-brand">
        <p class="helper-label text-sky-700">Backend</p>
        <p class="mt-2 text-sm leading-6 text-slate-700">
          Login, me, auth middleware, role middleware, member API, wallet summary plus transaction API, serta wifi billing with audit log.
        </p>
      </div>

      <div class="helper-box">
        <p class="helper-label">Frontend</p>
        <p class="mt-2 text-sm leading-6 text-slate-600">
          Login end-to-end, route protection, dashboard summary, wallet list/form, members list, dan halaman wifi mobile-first.
        </p>
      </div>
    </div>
  </PageCard>
</div>
