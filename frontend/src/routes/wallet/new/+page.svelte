<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';
  import StatePanel from '$lib/components/StatePanel.svelte';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;
</script>

<div class="space-y-4">
  <PageCard
    title="Tambah Transaksi Kas"
    description="Form pencatatan Kantong Duafa. Pembayaran tetap dilakukan di luar aplikasi, halaman ini hanya mencatat transaksi."
  >
    {#if data.accessDenied}
      <StatePanel
        tone="forbidden"
        title="Akses ditolak"
        message="Pencatatan transaksi kas hanya tersedia untuk admin dan bendahara."
      />
    {:else}
      <form method="POST" class="space-y-4">
        {#if form?.message}
          <StatePanel
            tone="error"
            title="Gagal memproses"
            message={form.message}
            requestId={form && 'requestId' in form && typeof form.requestId === 'string' ? form.requestId : null}
          />
        {/if}

        <div class="grid gap-4 sm:grid-cols-2">
          <label>
            <span class="field-label">Jenis transaksi</span>
            <select name="type" class="input-field" required>
              <option value="income" selected={(form?.values?.type ?? 'income') === 'income'}>Pemasukan</option>
              <option value="expense" selected={(form?.values?.type ?? '') === 'expense'}>Pengeluaran</option>
            </select>
          </label>

          <label>
            <span class="field-label">Nominal</span>
            <input
              class="input-field"
              type="number"
              name="amount"
              min="1"
              step="1"
              inputmode="numeric"
              placeholder="20000"
              value={form?.values?.amount ?? ''}
              required
            />
          </label>
        </div>

        <label>
          <span class="field-label">Kategori</span>
          <input
            class="input-field"
            type="text"
            name="category"
            placeholder="Hibah Orang Baik"
            value={form?.values?.category ?? ''}
            required
          />
        </label>

        <label>
          <span class="field-label">Deskripsi</span>
          <textarea
            class="input-field min-h-32 resize-y"
            name="description"
            placeholder="Catatan singkat transaksi"
            required
          >{form?.values?.description ?? ''}</textarea>
        </label>

        <div class="flex items-center justify-between gap-3">
          <a href="/wallet" class="btn-secondary px-4 py-3">Kembali</a>
          <button type="submit" class="btn-primary px-4 py-3">Simpan transaksi</button>
        </div>
      </form>
    {/if}
  </PageCard>
</div>
