<script lang="ts">
  import { browser } from '$app/environment';
  import { enhance } from '$app/forms';
  import { navigating } from '$app/stores';
  import type { SubmitFunction } from '@sveltejs/kit';
  import { toast } from 'svelte-sonner';
  import ActionButtonGroup from '$lib/components/ActionButtonGroup.svelte';
  import ActionSheet from '$lib/components/ActionSheet.svelte';
  import AppIcon from '$lib/components/AppIcon.svelte';
  import ConfirmDialog from '$lib/components/ConfirmDialog.svelte';
  import FeedbackBanner from '$lib/components/FeedbackBanner.svelte';
  import PageCard from '$lib/components/PageCard.svelte';
  import PageSkeleton from '$lib/components/PageSkeleton.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import { createConfirmableSubmitController } from '$lib/forms/confirmable';
  import type { ActionData, PageData } from './$types';
  import type { MemberUser, UserRole } from '$lib/api/types';

  export let data: PageData;
  export let form: ActionData;

  const roleLabels: Record<UserRole, string> = {
    admin: 'Admin',
    treasurer: 'Bendahara',
    member: 'Anggota'
  };

  let pendingAction: string | null = null;
  let searchTerm = '';
  let roleFilter: 'all' | UserRole = 'all';
  let statusFilter: 'all' | 'active' | 'inactive' | 'archived' = 'all';
  let createSheetOpen = false;
  let editTargetId: string | null = null;
  let passwordTargetId: string | null = null;
  let lastToastKey = '';
  const confirmController = createConfirmableSubmitController({
    setPendingAction: (value) => {
      pendingAction = value;
    }
  });
  const confirmationState = confirmController.state;

  function enhanceWithAction(actionName: string): SubmitFunction {
    return () => {
      pendingAction = actionName;

      return async ({ update }) => {
        await update();
        pendingAction = null;
      };
    };
  }

  function enhanceWithConfirmation(
    actionName: string,
    confirmation: Parameters<typeof confirmController.enhance>[1]
  ): SubmitFunction {
    return confirmController.enhance(actionName, confirmation);
  }

  function roleBadgeClass(role: UserRole) {
    if (role === 'admin') {
      return 'badge-strong';
    }

    if (role === 'treasurer') {
      return 'badge-success';
    }

    return 'badge-muted';
  }

  function memberStatus(member: MemberUser) {
    if (member.archived_at) {
      return 'archived' as const;
    }

    return member.is_active ? 'active' : 'inactive';
  }

  function statusBadgeClass(status: ReturnType<typeof memberStatus>) {
    if (status === 'active') {
      return 'badge-brand';
    }

    if (status === 'archived') {
      return 'badge-strong';
    }

    return 'badge-danger';
  }

  function statusLabel(status: ReturnType<typeof memberStatus>) {
    if (status === 'active') {
      return 'Aktif';
    }

    if (status === 'archived') {
      return 'Arsip';
    }

    return 'Nonaktif';
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

  function normalizeSearch(value: string) {
    return value.trim().toLowerCase();
  }

  function matchesFilters(member: MemberUser) {
    const query = normalizeSearch(searchTerm);
    const inSearch =
      query === '' ||
      normalizeSearch(member.name).includes(query) ||
      normalizeSearch(member.email).includes(query) ||
      normalizeSearch(member.username).includes(query);

    const inRole = roleFilter === 'all' || member.role === roleFilter;
    const currentStatus = memberStatus(member);
    const inStatus = statusFilter === 'all' || currentStatus === statusFilter;

    return inSearch && inRole && inStatus;
  }

  function countByRole(role: UserRole) {
    return data.members.filter((member) => member.role === role).length;
  }

  function createValue(field: string) {
    if (form?.action !== 'createMember') {
      if (field === 'role') {
        return 'member';
      }

      if (field === 'is_active') {
        return 'true';
      }

      return '';
    }

    const values = form.values as Record<string, string> | undefined;
    return values?.[field] ?? '';
  }

  function editValue(member: MemberUser, field: string) {
    if (form?.action === 'updateMember') {
      const values = form.values as Record<string, string> | undefined;
      if (values?.member_id === member.id) {
        return values?.[field] ?? '';
      }
    }

    if (field === 'name') {
      return member.name;
    }

    if (field === 'email') {
      return member.email;
    }

    if (field === 'username') {
      return member.username;
    }

    if (field === 'phone') {
      return member.phone ?? '';
    }

    if (field === 'joined_at') {
      return member.joined_at ? member.joined_at.slice(0, 10) : '';
    }

    if (field === 'role') {
      return member.role;
    }

    if (field === 'is_active') {
      return member.is_active ? 'true' : 'false';
    }

    return '';
  }

  function passwordValue(field: 'new_password' | 'confirm_password') {
    if (form?.action !== 'resetPassword') {
      return '';
    }

    const values = form.values as Record<string, string> | undefined;
    return values?.[field] ?? '';
  }

  function isSheetAction(action: string | undefined) {
    return action === 'createMember' || action === 'updateMember' || action === 'resetPassword';
  }

  function formRequestId() {
    return form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null;
  }

  function openCreateSheet() {
    createSheetOpen = true;
    editTargetId = null;
    passwordTargetId = null;
  }

  function openEditSheet(memberID: string) {
    createSheetOpen = false;
    editTargetId = memberID;
    passwordTargetId = null;
  }

  function openPasswordSheet(memberID: string) {
    createSheetOpen = false;
    editTargetId = null;
    passwordTargetId = memberID;
  }

  function closeSheets() {
    createSheetOpen = false;
    editTargetId = null;
    passwordTargetId = null;
  }

  $: filteredMembers = data.members.filter(matchesFilters);
  $: selectedMember = data.members.find((member) => member.id === editTargetId) ?? null;
  $: passwordTarget = data.members.find((member) => member.id === passwordTargetId) ?? null;
  $: createError = form?.action === 'createMember' ? form.message ?? null : null;
  $: updateError = form?.action === 'updateMember' ? form.message ?? null : null;
  $: passwordError = form?.action === 'resetPassword' ? form.message ?? null : null;
  $: pageError = form?.message && !isSheetAction(form.action) ? form.message : null;
  $: pageSuccess = form?.success && !isSheetAction(form.action) ? form.success : null;
  $: if (form?.action === 'createMember' && form?.message) {
    createSheetOpen = true;
  }
  $: if (form?.action === 'updateMember' && form?.message) {
    const values = form.values as Record<string, string> | undefined;
    editTargetId = values?.member_id ?? editTargetId;
  }
  $: if (form?.action === 'resetPassword' && form?.message) {
    const values = form.values as Record<string, string> | undefined;
    passwordTargetId = values?.member_id ?? passwordTargetId;
  }
  $: if (form?.success) {
    closeSheets();
  }
  $: if (browser && form?.message) {
    const key = `error:${form.action ?? 'unknown'}:${formRequestId() ?? form.message}`;

    if (key !== lastToastKey) {
      toast.error(form.message);
      lastToastKey = key;
    }
  }
  $: if (browser && form?.success) {
    const key = `success:${form.action ?? 'unknown'}:${form.success}`;

    if (key !== lastToastKey) {
      toast.success(form.success);
      lastToastKey = key;
    }
  }
  $: confirmationDialog = $confirmationState.dialog;
  $: confirmationLoading =
    !!confirmationDialog &&
    ($confirmationState.requestingActionKey === confirmationDialog.actionKey || pendingAction === confirmationDialog.actionKey);
</script>

<div class="space-y-5 lg:space-y-6">
  <PageCard
    eyebrow="Members"
    icon="lucide:users"
    title="Anggota Mess"
    description="Kelola penghuni, role, status aktif, dan akses akun dari satu halaman yang lebih rapi dipakai admin."
  >
    <svelte:fragment slot="actions">
      {#if data.canManage}
        <div class="flex flex-wrap gap-3">
          <a href="/admin/import/members" class="btn-secondary px-4 py-3">Impor CSV</a>
          <button type="button" class="btn-primary px-4 py-3" on:click={openCreateSheet}>
            Tambah member
          </button>
        </div>
      {/if}
    </svelte:fragment>

    {#if $navigating?.to?.url.pathname === '/members'}
      <PageSkeleton statCards={4} rows={4} />
    {/if}

    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Akses ditolak"
        message="Daftar anggota hanya dapat dibuka oleh admin dan bendahara."
      />
    {:else if data.loadError}
      <StatePanel tone="error" title="Gagal memuat anggota" message={data.loadError} />
    {:else}
      {#if pageError}
        <StatePanel
          tone="error"
          title="Perubahan belum tersimpan"
          message={pageError}
          requestId={formRequestId()}
        />
      {:else if pageSuccess}
        <FeedbackBanner tone="success" title="Berhasil" message={pageSuccess} />
      {/if}

      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <div class="stat-card bg-slate-950 text-white">
          <p class="helper-label text-slate-300">Total member</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em]">{data.summary.total}</p>
          <p class="mt-2 text-sm text-slate-300">Termasuk anggota aktif dan nonaktif.</p>
        </div>

        <div class="stat-card bg-white">
          <p class="helper-label">Aktif</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.summary.active}</p>
          <p class="mt-2 text-sm text-muted">Masih tercatat tinggal di mess.</p>
        </div>

        <div class="stat-card bg-white">
          <p class="helper-label">Admin & bendahara</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">
            {countByRole('admin') + countByRole('treasurer')}
          </p>
          <p class="mt-2 text-sm text-muted">
            {countByRole('admin')} admin, {countByRole('treasurer')} bendahara.
          </p>
        </div>

        <div class="stat-card bg-white">
          <p class="helper-label">Nonaktif & arsip</p>
          <p class="mt-2 text-3xl font-semibold tracking-[-0.04em] text-ink">{data.summary.inactive + data.summary.archived}</p>
          <p class="mt-2 text-sm text-muted">{data.summary.inactive} nonaktif, {data.summary.archived} arsip.</p>
        </div>
      </div>

      <section class="rounded-[30px] border border-line bg-panel/80 p-5 sm:p-6">
        <div class="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
          <div class="space-y-2">
            <p class="eyebrow">Cari & filter</p>
            <h3 class="section-title text-[1.1rem] sm:text-[1.25rem]">Rapikan daftar anggota yang sedang dibuka</h3>
            <p class="section-subtitle">
              Menampilkan {filteredMembers.length} dari {data.members.length} anggota dengan filter yang aktif saat ini.
            </p>
          </div>

          <div class="grid flex-1 gap-4 sm:grid-cols-2 xl:grid-cols-[minmax(0,1.4fr)_repeat(2,minmax(0,0.7fr))]">
            <label class="sm:col-span-2 xl:col-span-1">
              <span class="field-label">Cari member</span>
              <div class="relative">
                <span class="pointer-events-none absolute inset-y-0 left-4 flex items-center text-dusty">
                  <AppIcon icon="lucide:search" className="h-4 w-4" />
                </span>
                <input
                  class="input-field pl-11"
                  type="search"
                  placeholder="Nama, email, atau username"
                  bind:value={searchTerm}
                />
              </div>
            </label>

            <label>
              <span class="field-label">Filter role</span>
              <select class="input-field" bind:value={roleFilter}>
                <option value="all">Semua role</option>
                <option value="admin">Admin</option>
                <option value="treasurer">Bendahara</option>
                <option value="member">Anggota</option>
              </select>
            </label>

            <label>
              <span class="field-label">Filter status</span>
              <select class="input-field" bind:value={statusFilter}>
                <option value="all">Semua status</option>
                <option value="active">Aktif</option>
                <option value="inactive">Nonaktif</option>
                <option value="archived">Arsip</option>
              </select>
            </label>
          </div>

          {#if searchTerm !== '' || roleFilter !== 'all' || statusFilter !== 'all'}
            <button
              type="button"
              class="btn-secondary shrink-0 px-4 py-3"
              on:click={() => {
                searchTerm = '';
                roleFilter = 'all';
                statusFilter = 'all';
              }}
            >
              Reset filter
            </button>
          {/if}
        </div>
      </section>

      {#if data.members.length === 0}
        <StatePanel
          tone="empty"
          title="Daftar anggota masih kosong"
          message="Tambahkan member pertama atau impor CSV penghuni supaya pengelolaan role, wifi, dan operasional harian bisa langsung dimulai."
          actionHref={data.canManage ? '/admin/import/members' : null}
          actionLabel={data.canManage ? 'Buka impor anggota' : ''}
        >
          {#if data.canManage}
            <div class="state-panel-actions">
              <button type="button" class="btn-primary px-4 py-3" on:click={openCreateSheet}>
                Tambah member pertama
              </button>
            </div>
          {/if}
        </StatePanel>
      {:else if filteredMembers.length === 0}
        <StatePanel
          tone="empty"
          title="Tidak ada member yang cocok"
          message="Coba longgarkan pencarian atau ubah filter role dan status supaya daftar anggota muncul lagi."
          icon="lucide:search-x"
        >
          <div class="state-panel-actions">
            <button
              type="button"
              class="btn-secondary px-4 py-3"
              on:click={() => {
                searchTerm = '';
                roleFilter = 'all';
                statusFilter = 'all';
              }}
            >
              Tampilkan semua member
            </button>
          </div>
        </StatePanel>
      {:else}
        <div class="grid gap-5 xl:grid-cols-2">
          {#each filteredMembers as member}
            <article class="stat-card bg-white">
              <div class="flex flex-col gap-5">
                <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                  <div class="min-w-0 space-y-3">
                    <div class="flex flex-wrap items-start gap-3">
                      <div class="avatar-chip avatar-chip-sm mt-0.5">{member.name.slice(0, 1)}</div>

                      <div class="min-w-0 flex-1">
                        <div class="flex flex-wrap items-center gap-2.5">
                          <h3 class="text-lg font-semibold text-ink">{member.name}</h3>
                          <span class={roleBadgeClass(member.role)}>{roleLabels[member.role]}</span>
                          <span class={statusBadgeClass(memberStatus(member))}>
                            {statusLabel(memberStatus(member))}
                          </span>
                        </div>
                        <p class="mt-2 break-all text-sm text-muted">{member.email}</p>
                        <p class="mt-1 text-xs uppercase tracking-[0.16em] text-dusty">@{member.username}</p>
                      </div>
                    </div>
                  </div>

                  <ActionButtonGroup align="end">
                    {#if data.canManage}
                      <button type="button" class="btn-secondary px-3.5 py-2.5" on:click={() => openEditSheet(member.id)}>
                        <AppIcon icon="lucide:square-pen" className="h-4 w-4" />
                        <span>Edit</span>
                      </button>
                      <button type="button" class="btn-secondary px-3.5 py-2.5" on:click={() => openPasswordSheet(member.id)}>
                        <AppIcon icon="lucide:key-round" className="h-4 w-4" />
                        <span>Reset password</span>
                      </button>
                    {/if}
                  </ActionButtonGroup>
                </div>

                <div class="content-grid-tight">
                  <div class="meta-card">
                    <p class="helper-label">Nomor HP</p>
                    <p class="mt-2 text-sm font-medium text-ink">{member.phone ?? '-'}</p>
                  </div>

                  <div class="meta-card">
                    <p class="helper-label">Mulai tinggal</p>
                    <p class="mt-2 text-sm font-medium text-ink">{formatDate(member.joined_at)}</p>
                  </div>

                  <div class="meta-card">
                    <p class="helper-label">Status keluar</p>
                    <p class="mt-2 text-sm font-medium text-ink">
                      {member.archived_at ? `Diarsipkan ${formatDate(member.archived_at)}` : formatDate(member.left_at)}
                    </p>
                  </div>

                  <div class="meta-card">
                    <p class="helper-label">Terakhir diperbarui</p>
                    <p class="mt-2 text-sm font-medium text-ink">{formatDate(member.updated_at)}</p>
                  </div>
                </div>

                {#if data.canManage}
                  <ActionButtonGroup bordered>
                    {#if !member.archived_at}
                      <form
                        method="POST"
                        action="?/toggleActive"
                        use:enhance={enhanceWithConfirmation(`toggle-${member.id}`, {
                          title: member.is_active ? `Nonaktifkan ${member.name}?` : `Aktifkan kembali ${member.name}?`,
                          description: member.is_active
                            ? 'Akun akan keluar dari siklus aktif mess. Riwayat lama tetap aman dan bisa diaktifkan lagi kapan saja.'
                            : 'Akun akan kembali masuk ke siklus aktif mess dan dipakai lagi untuk operasional berikutnya.',
                          confirmLabel: member.is_active ? 'Nonaktifkan' : 'Aktifkan kembali',
                          icon: member.is_active ? 'lucide:user-x' : 'lucide:user-check',
                          destructive: member.is_active
                        })}
                      >
                        <input type="hidden" name="member_id" value={member.id} />
                        <input type="hidden" name="is_active" value={member.is_active ? 'false' : 'true'} />
                        <button
                          type="submit"
                          class={member.is_active ? 'btn-secondary px-4 py-3' : 'btn-primary px-4 py-3'}
                          disabled={pendingAction === `toggle-${member.id}`}
                        >
                          <AppIcon
                            icon={member.is_active ? 'lucide:user-x' : 'lucide:user-check'}
                            className="h-4 w-4"
                          />
                          <span>
                            {pendingAction === `toggle-${member.id}`
                              ? 'Menyimpan...'
                              : member.is_active
                                ? 'Nonaktifkan akun'
                                : 'Aktifkan kembali'}
                          </span>
                        </button>
                      </form>
                    {/if}

                    {#if member.archived_at}
                      <form
                        method="POST"
                        action="?/reactivateMember"
                        use:enhance={enhanceWithConfirmation(`reactivate-${member.id}`, {
                          title: `Keluarkan ${member.name} dari arsip?`,
                          description: 'Akun akan kembali aktif dan bisa dipakai lagi untuk login, wifi, dan operasional mess berikutnya.',
                          confirmLabel: 'Aktifkan dari arsip',
                          icon: 'lucide:archive-restore'
                        })}
                      >
                        <input type="hidden" name="member_id" value={member.id} />
                        <button
                          type="submit"
                          class="btn-primary px-4 py-3"
                          disabled={pendingAction === `reactivate-${member.id}`}
                        >
                          <AppIcon icon="lucide:archive-restore" className="h-4 w-4" />
                          <span>{pendingAction === `reactivate-${member.id}` ? 'Menyimpan...' : 'Aktifkan dari arsip'}</span>
                        </button>
                      </form>
                    {:else}
                      <form
                        method="POST"
                        action="?/archiveMember"
                        use:enhance={enhanceWithConfirmation(`archive-${member.id}`, {
                          title: `Arsipkan ${member.name}?`,
                          description: 'Akun akan keluar dari lifecycle aktif, tetapi seluruh histori tetap utuh dan masih bisa dipulihkan bila diperlukan.',
                          confirmLabel: 'Arsipkan',
                          icon: 'lucide:archive',
                          destructive: true
                        })}
                      >
                        <input type="hidden" name="member_id" value={member.id} />
                        <button
                          type="submit"
                          class="btn-secondary px-4 py-3"
                          disabled={pendingAction === `archive-${member.id}`}
                        >
                          <AppIcon icon="lucide:archive" className="h-4 w-4" />
                          <span>{pendingAction === `archive-${member.id}` ? 'Mengarsipkan...' : 'Arsipkan akun'}</span>
                        </button>
                      </form>
                    {/if}

                    {#if member.archived_at}
                      <form
                        method="POST"
                        action="?/deletePermanent"
                        use:enhance={enhanceWithConfirmation(`delete-${member.id}`, {
                          title: `Hapus permanen ${member.name}?`,
                          description: 'Tindakan ini hanya akan berhasil jika akun benar-benar belum punya relasi penting. Jika histori masih ada, sistem akan menolak penghapusan demi keamanan data.',
                          confirmLabel: 'Hapus permanen',
                          icon: 'lucide:trash-2',
                          destructive: true
                        })}
                      >
                        <input type="hidden" name="member_id" value={member.id} />
                        <button
                          type="submit"
                          class="btn-danger px-4 py-3"
                          disabled={pendingAction === `delete-${member.id}`}
                        >
                          <AppIcon icon="lucide:trash-2" className="h-4 w-4" />
                          <span>{pendingAction === `delete-${member.id}` ? 'Menghapus...' : 'Hapus permanen'}</span>
                        </button>
                      </form>
                    {/if}

                    <button type="button" class="btn-secondary px-4 py-3" on:click={() => openEditSheet(member.id)}>
                      <AppIcon icon="lucide:shield-ellipsis" className="h-4 w-4" />
                      <span>Ubah role & data</span>
                    </button>
                  </ActionButtonGroup>
                {/if}
              </div>
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>

  {#if data.canManage}
    <ActionSheet
      open={createSheetOpen}
      title="Tambah member baru"
      description="Buat akun penghuni baru, tentukan role, status aktif, dan password awal dari satu panel."
      icon="lucide:user-plus"
      on:close={closeSheets}
    >
      <form method="POST" action="?/createMember" class="space-y-4" use:enhance={enhanceWithAction('createMember')}>
        {#if createError}
          <div class="space-y-2">
            <FeedbackBanner tone="error" title="Anggota belum tersimpan" message={createError} />
            {#if formRequestId()}
              <p class="text-xs text-dusty">Request ID: {formRequestId()}</p>
            {/if}
          </div>
        {/if}

        <div class="grid gap-4 sm:grid-cols-2">
          <label class="sm:col-span-2">
            <span class="field-label">Nama lengkap</span>
            <input class="input-field" type="text" name="name" value={createValue('name')} required />
          </label>

          <label>
            <span class="field-label">Email</span>
            <input class="input-field" type="email" name="email" value={createValue('email')} required />
          </label>

          <label>
            <span class="field-label">Username</span>
            <input class="input-field" type="text" name="username" value={createValue('username')} placeholder="Opsional, akan digenerate bila kosong" />
          </label>

          <label>
            <span class="field-label">Nomor HP</span>
            <input class="input-field" type="text" name="phone" value={createValue('phone')} placeholder="Opsional" />
          </label>

          <label>
            <span class="field-label">Mulai tinggal</span>
            <input class="input-field" type="date" name="joined_at" value={createValue('joined_at')} />
          </label>

          <label>
            <span class="field-label">Role</span>
            <select class="input-field" name="role">
              <option value="admin" selected={createValue('role') === 'admin'}>Admin</option>
              <option value="treasurer" selected={createValue('role') === 'treasurer'}>Bendahara</option>
              <option value="member" selected={createValue('role') === 'member' || createValue('role') === ''}>Anggota</option>
            </select>
          </label>

          <label>
            <span class="field-label">Status akun</span>
            <select class="input-field" name="is_active">
              <option value="true" selected={createValue('is_active') !== 'false'}>Aktif</option>
              <option value="false" selected={createValue('is_active') === 'false'}>Nonaktif</option>
            </select>
          </label>

          <label>
            <span class="field-label">Password awal</span>
            <input class="input-field" type="password" name="password" value={createValue('password')} required />
          </label>

          <label>
            <span class="field-label">Konfirmasi password</span>
            <input class="input-field" type="password" name="confirm_password" value={createValue('confirm_password')} required />
          </label>
        </div>

        <div class="helper-box-brand">
          <p class="helper-label">Catatan</p>
          <p class="mt-2 text-sm leading-6 text-muted">
            Jika username dikosongkan, sistem akan membuat username unik otomatis dari nama atau email anggota.
          </p>
        </div>

        <ActionButtonGroup bordered>
          <button
            type="submit"
            class="btn-primary px-4 py-3"
            disabled={pendingAction === 'createMember'}
          >
            {pendingAction === 'createMember' ? 'Menyimpan...' : 'Simpan member baru'}
          </button>
          <button type="button" class="btn-secondary px-4 py-3" on:click={closeSheets}>
            Batal
          </button>
        </ActionButtonGroup>
      </form>
    </ActionSheet>

    {#if selectedMember}
      <ActionSheet
        open={!!selectedMember}
        title={`Edit ${selectedMember.name}`}
        description="Perbarui data dasar, role, status aktif, dan tanggal mulai tinggal tanpa meninggalkan halaman."
        icon="lucide:user-cog"
        on:close={closeSheets}
      >
        <form method="POST" action="?/updateMember" class="space-y-4" use:enhance={enhanceWithAction(`update-${selectedMember.id}`)}>
          <input type="hidden" name="member_id" value={selectedMember.id} />

          {#if updateError}
            <div class="space-y-2">
              <FeedbackBanner tone="error" title="Perubahan belum tersimpan" message={updateError} />
              {#if formRequestId()}
                <p class="text-xs text-dusty">Request ID: {formRequestId()}</p>
              {/if}
            </div>
          {/if}

          <div class="grid gap-4 sm:grid-cols-2">
            <label class="sm:col-span-2">
              <span class="field-label">Nama lengkap</span>
              <input class="input-field" type="text" name="name" value={editValue(selectedMember, 'name')} required />
            </label>

            <label>
              <span class="field-label">Email</span>
              <input class="input-field" type="email" name="email" value={editValue(selectedMember, 'email')} required />
            </label>

            <label>
              <span class="field-label">Username</span>
              <input class="input-field" type="text" name="username" value={editValue(selectedMember, 'username')} required />
            </label>

            <label>
              <span class="field-label">Nomor HP</span>
              <input class="input-field" type="text" name="phone" value={editValue(selectedMember, 'phone')} placeholder="Kosongkan jika tidak ada" />
            </label>

            <label>
              <span class="field-label">Mulai tinggal</span>
              <input class="input-field" type="date" name="joined_at" value={editValue(selectedMember, 'joined_at')} />
            </label>

            <label>
              <span class="field-label">Role</span>
              <select class="input-field" name="role">
                <option value="admin" selected={editValue(selectedMember, 'role') === 'admin'}>Admin</option>
                <option value="treasurer" selected={editValue(selectedMember, 'role') === 'treasurer'}>Bendahara</option>
                <option value="member" selected={editValue(selectedMember, 'role') === 'member'}>Anggota</option>
              </select>
            </label>

            <label>
              <span class="field-label">Status akun</span>
              <select class="input-field" name="is_active">
                <option value="true" selected={editValue(selectedMember, 'is_active') === 'true'} disabled={!!selectedMember.archived_at}>Aktif</option>
                <option value="false" selected={editValue(selectedMember, 'is_active') === 'false'}>Nonaktif</option>
              </select>
            </label>
          </div>

          <div class="helper-box">
            <p class="helper-label">Keamanan admin</p>
            <p class="mt-2 text-sm leading-6 text-muted">
              Sistem akan menolak jika perubahan ini membuat tidak ada admin aktif yang tersisa, jika admin mencoba menurunkan role akun dirinya sendiri, atau jika akun yang sudah diarsipkan dipaksa aktif lagi dari form edit biasa.
            </p>
          </div>

          <ActionButtonGroup bordered>
            <button
              type="submit"
              class="btn-primary px-4 py-3"
              disabled={pendingAction === `update-${selectedMember.id}`}
            >
              {pendingAction === `update-${selectedMember.id}` ? 'Menyimpan...' : 'Simpan perubahan'}
            </button>
            <button type="button" class="btn-secondary px-4 py-3" on:click={closeSheets}>
              Tutup
            </button>
          </ActionButtonGroup>
        </form>
      </ActionSheet>
    {/if}

    {#if passwordTarget}
      <ActionSheet
        open={!!passwordTarget}
        title={`Reset password ${passwordTarget.name}`}
        description="Masukkan password baru yang akan langsung menggantikan password lama akun anggota."
        icon="lucide:key-round"
        on:close={closeSheets}
      >
        <form method="POST" action="?/resetPassword" class="space-y-4" use:enhance={enhanceWithAction(`password-${passwordTarget.id}`)}>
          <input type="hidden" name="member_id" value={passwordTarget.id} />

          {#if passwordError}
            <div class="space-y-2">
              <FeedbackBanner tone="error" title="Password belum direset" message={passwordError} />
              {#if formRequestId()}
                <p class="text-xs text-dusty">Request ID: {formRequestId()}</p>
              {/if}
            </div>
          {/if}

          <div class="helper-box-brand">
            <p class="helper-label">Perhatian</p>
            <p class="mt-2 text-sm leading-6 text-muted">
              Password baru tidak akan ditampilkan lagi setelah tersimpan. Pastikan admin menyampaikan password ini dengan aman ke anggota terkait.
            </p>
          </div>

          <label>
            <span class="field-label">Password baru</span>
            <input class="input-field" type="password" name="new_password" value={passwordValue('new_password')} required />
          </label>

          <label>
            <span class="field-label">Konfirmasi password baru</span>
            <input class="input-field" type="password" name="confirm_password" value={passwordValue('confirm_password')} required />
          </label>

          <ActionButtonGroup bordered>
            <button
              type="submit"
              class="btn-primary px-4 py-3"
              disabled={pendingAction === `password-${passwordTarget.id}`}
            >
              {pendingAction === `password-${passwordTarget.id}` ? 'Menyimpan...' : 'Reset password'}
            </button>
            <button type="button" class="btn-secondary px-4 py-3" on:click={closeSheets}>
              Batal
            </button>
          </ActionButtonGroup>
        </form>
      </ActionSheet>
    {/if}
  {/if}

  <ConfirmDialog
    open={!!confirmationDialog}
    title={confirmationDialog?.title ?? 'Konfirmasi aksi'}
    description={confirmationDialog?.description ?? ''}
    confirmLabel={confirmationDialog?.confirmLabel ?? 'Lanjutkan'}
    cancelLabel={confirmationDialog?.cancelLabel ?? 'Batal'}
    icon={confirmationDialog?.icon ?? 'lucide:triangle-alert'}
    destructive={confirmationDialog?.destructive ?? false}
    loading={confirmationLoading}
    on:close={confirmController.closeDialog}
    on:confirm={confirmController.confirmDialog}
  />
</div>
