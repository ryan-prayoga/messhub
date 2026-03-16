<script lang="ts">
  import { navigating } from '$app/stores';
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ActionData, PageData } from './$types';
  import type { UserRole } from '$lib/api/types';

  export let data: PageData;
  export let form: ActionData;

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

  function roleValue(memberID: string, fallback: UserRole) {
    if (form?.action !== 'updateRole') {
      return fallback;
    }

    const values = form.values as Record<string, string> | undefined;
    if (values?.member_id !== memberID) {
      return fallback;
    }

    if (values.role === 'admin' || values.role === 'treasurer' || values.role === 'member') {
      return values.role;
    }

    return fallback;
  }
</script>

<div class="space-y-4">
  <PageCard
    title="Members"
    description="Daftar anggota mess dengan admin controls untuk role dan status aktif."
  >
    {#if $navigating?.to?.url.pathname === '/members'}
      <StatePanel tone="loading" title="Loading" message="Memuat ulang data anggota mess..." />
    {/if}

    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Forbidden"
        message={`Role ${data.user?.role} belum diizinkan melihat daftar anggota. Akses tersedia untuk admin dan treasurer.`}
      />
    {:else if data.loadError}
      <StatePanel tone="error" title="Error" message={data.loadError} />
    {:else}
      {#if form?.message}
        <StatePanel
          tone="error"
          title="Error"
          message={form.message}
          requestId={form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null}
        />
      {:else if form?.success}
        <div class="helper-box mb-4 border-emerald-200 bg-emerald-50/80">
          <p class="helper-label text-emerald-700">Success</p>
          <p class="mt-2 text-sm leading-6 text-emerald-800">{form.success}</p>
        </div>
      {/if}

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
        <StatePanel
          tone="empty"
          title="Empty"
          message="Belum ada anggota yang tampil dari backend. Setelah admin menambahkan user baru, daftar ini akan langsung terisi."
        />
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

              <div class="mt-3 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3">
                <p class="helper-label">Phone</p>
                <p class="mt-2 text-sm font-medium text-ink">{member.phone ?? '-'}</p>
              </div>

              {#if data.canManage}
                <div class="mt-4 grid gap-3 sm:grid-cols-[minmax(0,1fr)_auto]">
                  <form method="POST" action="?/updateRole" class="space-y-3">
                    <input type="hidden" name="member_id" value={member.id} />
                    <label>
                      <span class="field-label">Role</span>
                      <select name="role" class="input-field">
                        <option value="admin" selected={roleValue(member.id, member.role) === 'admin'}>Admin</option>
                        <option value="treasurer" selected={roleValue(member.id, member.role) === 'treasurer'}>Treasurer</option>
                        <option value="member" selected={roleValue(member.id, member.role) === 'member'}>Member</option>
                      </select>
                    </label>
                    <button type="submit" class="btn-secondary w-full px-4 py-3">Save role</button>
                  </form>

                  <form method="POST" action="?/toggleActive" class="flex items-end">
                    <input type="hidden" name="member_id" value={member.id} />
                    <input type="hidden" name="is_active" value={member.is_active ? 'false' : 'true'} />
                    <button
                      type="submit"
                      class={member.is_active ? 'btn-secondary w-full px-4 py-3' : 'btn-primary w-full px-4 py-3'}
                    >
                      {member.is_active ? 'Deactivate' : 'Activate'}
                    </button>
                  </form>
                </div>
              {/if}
            </article>
          {/each}
        </div>
      {/if}
    {/if}
  </PageCard>
</div>
