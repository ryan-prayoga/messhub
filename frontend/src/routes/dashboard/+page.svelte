<script lang="ts">
  import { invalidateAll } from '$app/navigation';
  import { navigating } from '$app/stores';
  import AppIcon from '$lib/components/AppIcon.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import PullToRefresh from '$lib/components/PullToRefresh.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { PageData } from './$types';
  import type { UserRole } from '$lib/api/types';

  export let data: PageData;

  const roleLabels: Record<UserRole, string> = {
    admin: 'Admin',
    treasurer: 'Bendahara',
    member: 'Anggota'
  };

  function roleBadgeClass(role: UserRole) {
    if (role === 'admin') {
      return 'badge-strong';
    }

    if (role === 'treasurer') {
      return 'badge-success';
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
      return 'Menunggu verifikasi';
    }

    if (status === 'verified') {
      return 'Terverifikasi';
    }

    if (status === 'rejected') {
      return 'Ditolak';
    }

    if (status === 'unpaid') {
      return 'Belum bayar';
    }

    return 'Belum ada tagihan aktif';
  }

  function wifiStatusClass(status: string | null) {
    if (status === 'verified') {
      return 'badge-success';
    }

    if (status === 'pending_verification') {
      return 'badge-warning';
    }

    if (status === 'rejected') {
      return 'badge-danger';
    }

    if (status === 'unpaid') {
      return 'badge-muted';
    }

    return 'badge-info';
  }

  const adminQuickLinks = [
    {
      href: '/admin/import',
      label: 'Pusat impor',
      description: 'Impor anggota dan transaksi kas dari CSV dengan preview sebelum commit.',
      icon: 'lucide:download'
    },
    {
      href: '/settings',
      label: 'Pengaturan',
      description: 'Atur nama mess, nominal wifi, deadline, dan rekening tujuan.',
      icon: 'lucide:settings-2'
    }
  ];

  const quickLinks = [
    {
      href: '/members',
      label: 'Members',
      description: 'Lihat daftar anggota, role, dan status aktif mess.',
      icon: 'lucide:users'
    },
    {
      href: '/wallet',
      label: 'Wallet',
      description: 'Pantau saldo kas dan riwayat transaksi terbaru.',
      icon: 'lucide:wallet'
    },
    {
      href: '/wifi',
      label: 'Wifi',
      description: 'Cek tagihan aktif, deadline, dan status pembayaran.',
      icon: 'lucide:wifi'
    }
  ];

  async function refreshPage() {
    await invalidateAll();
  }
</script>

<PullToRefresh onRefresh={refreshPage}>
<div class="space-y-4">
  {#if $navigating?.to?.url.pathname === '/dashboard'}
    <PageSkeleton statCards={3} rows={2} />
  {/if}

  <PageCard
    eyebrow="Ringkasan Hari Ini"
    icon="lucide:sparkles"
    title={`Halo, ${data.user?.name ?? 'Anggota'}`}
    description="Ringkasan utama untuk kas, wifi, kontribusi, dan akses cepat yang paling sering dipakai."
  >
    <div class="grid gap-3 md:grid-cols-2">
      <div class="stat-card bg-white/80">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="helper-label">Akun aktif</p>
            <p class="mt-2 text-lg font-semibold text-ink">{data.user?.name}</p>
            <p class="mt-1 text-sm text-muted">{data.user?.email}</p>
          </div>

          <div class="nav-link-icon">
            <AppIcon icon="lucide:user-round" className="h-5 w-5" />
          </div>
        </div>

        {#if data.user}
          <div class="mt-4 flex flex-wrap items-center gap-2">
            <span class={roleBadgeClass(data.user.role)}>{roleLabels[data.user.role]}</span>
            <span class="badge-muted">@{data.user.username}</span>
          </div>
        {/if}
      </div>

      <div class="stat-card bg-white">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="helper-label">Ringkasan anggota</p>
            {#if data.memberSummary.state === 'ready'}
              <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.memberSummary.total}</p>
              <p class="mt-2 text-sm text-muted">
                {data.memberSummary.active} aktif, {data.memberSummary.inactive} nonaktif
              </p>
            {:else if data.memberSummary.state === 'restricted'}
              <p class="mt-2 text-lg font-semibold text-ink">Akses terbatas</p>
              <p class="mt-2 text-sm text-muted">
                Ringkasan anggota tersedia untuk admin dan bendahara.
              </p>
            {:else}
              <p class="mt-2 text-lg font-semibold text-ink">Belum tersedia</p>
              <p class="mt-2 text-sm text-muted">{data.memberSummary.message}</p>
            {/if}
          </div>

          <div class="nav-link-icon">
            <AppIcon icon="lucide:users" className="h-5 w-5" />
          </div>
        </div>
      </div>
    </div>

    <div class="mt-4 grid gap-3 sm:grid-cols-3">
      {#each quickLinks as link}
        <a href={link.href} class="stat-card bg-white/80">
          <div class="flex items-start justify-between gap-3">
            <div>
              <p class="helper-label">Menu cepat</p>
              <p class="mt-2 text-base font-semibold text-ink">{link.label}</p>
              <p class="mt-2 text-sm leading-6 text-muted">{link.description}</p>
            </div>

            <div class="nav-link-icon">
              <AppIcon icon={link.icon} className="h-5 w-5" />
            </div>
          </div>
        </a>
      {/each}
    </div>

    <div class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
      <a href="/wallet" class="stat-card bg-white">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="helper-label">Wallet</p>
            <p class="mt-2 text-base font-semibold text-ink">Kantong Duafa</p>
            {#if data.walletSummary.state === 'ready'}
              <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
                {formatCurrency(data.walletSummary.balance)}
              </p>
              <p class="mt-2 text-sm leading-6 text-muted">
                Pemasukan {formatCurrency(data.walletSummary.totalIncome)} / Pengeluaran {formatCurrency(data.walletSummary.totalExpense)}
              </p>
            {:else}
              <p class="mt-2 text-sm leading-6 text-muted">{data.walletSummary.message}</p>
            {/if}
          </div>

          <div class="nav-link-icon">
            <AppIcon icon="lucide:wallet" className="h-5 w-5" />
          </div>
        </div>
      </a>

      <a href="/wifi" class="stat-card bg-white sm:col-span-2 xl:col-span-2">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="helper-label">Wifi</p>
            <p class="mt-2 text-base font-semibold text-ink">Status bulan berjalan</p>
            {#if data.wifiSummary.state === 'ready'}
              <p class="mt-2 text-lg font-semibold text-ink">{data.wifiSummary.monthLabel}</p>
              {#if data.user?.role === 'member'}
                <div class="mt-2">
                  <span class={wifiStatusClass(data.wifiSummary.myStatus)}>
                    {wifiStatusLabel(data.wifiSummary.myStatus)}
                  </span>
                </div>
                <p class="mt-2 text-sm leading-6 text-muted">
                  Jatuh tempo {formatDate(data.wifiSummary.deadline)}.
                </p>
              {:else}
                <p class="mt-2 text-2xl font-semibold tracking-[-0.04em] text-ink">
                  {data.wifiSummary.verified} terverifikasi
                </p>
                <p class="mt-2 text-sm leading-6 text-muted">
                  {data.wifiSummary.unpaid} belum bayar / {data.wifiSummary.pending} menunggu verifikasi
                </p>
              {/if}
            {:else}
              <p class="mt-2 text-sm leading-6 text-muted">{data.wifiSummary.message}</p>
            {/if}
          </div>

          <div class="nav-link-icon">
            <AppIcon icon="lucide:wifi" className="h-5 w-5" />
          </div>
        </div>
      </a>
    </div>
  </PageCard>

  <PageCard
    eyebrow="Kontribusi"
    icon="lucide:trophy"
    title="Papan Kontribusi"
    description="Poin kontribusi bulan berjalan dihitung dari aktivitas bertipe kontribusi."
  >
    {#if data.leaderboardSummary.state === 'ready'}
      <div class="space-y-3">
        {#each data.leaderboardSummary.items as item}
          <div class="stat-card bg-white">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-3">
                <div class="nav-link-icon">
                  <AppIcon
                    icon={item.rank === 1 ? 'lucide:medal' : item.rank === 2 ? 'lucide:award' : 'lucide:star'}
                    className="h-5 w-5"
                  />
                </div>

                <div>
                  <p class="text-sm font-semibold text-ink">#{item.rank} {item.user_name}</p>
                  <p class="mt-1 text-xs text-muted">{item.total_activities} aktivitas tercatat</p>
                </div>
              </div>

              <div class="text-right">
                <p class="text-xl font-semibold tracking-[-0.04em] text-ink">{item.total_points}</p>
                <p class="text-xs uppercase tracking-[0.16em] text-dusty">poin</p>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {:else if data.leaderboardSummary.state === 'empty'}
      <StatePanel tone="empty" title="Belum ada kontribusi" message={data.leaderboardSummary.message ?? 'Belum ada kontribusi bulan ini'} />
    {:else}
      <StatePanel tone="error" title="Belum bisa memuat" message={data.leaderboardSummary.message ?? 'Data kontribusi belum tersedia.'} />
    {/if}
  </PageCard>

  {#if data.user?.role === 'admin'}
    <PageCard
      eyebrow="Admin"
      icon="lucide:shield-check"
      title="Akses Admin"
      description="Shortcut untuk pekerjaan administrasi yang paling sering dipakai di aplikasi."
    >
      <div class="grid gap-3 md:grid-cols-2">
        {#each adminQuickLinks as link}
          <a href={link.href} class="stat-card bg-white">
            <div class="flex items-start justify-between gap-3">
              <div>
                <p class="helper-label">Admin</p>
                <p class="mt-2 text-base font-semibold text-ink">{link.label}</p>
                <p class="mt-2 text-sm leading-6 text-muted">{link.description}</p>
              </div>

              <div class="nav-link-icon">
                <AppIcon icon={link.icon} className="h-5 w-5" />
              </div>
            </div>
          </a>
        {/each}
      </div>
    </PageCard>
  {/if}
</div>
</PullToRefresh>
