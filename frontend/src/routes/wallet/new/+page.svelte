<script lang="ts">
  import PageCard from '$lib/components/PageCard.svelte';
  import type { ActionData, PageData } from './$types';

  export let data: PageData;
  export let form: ActionData;
</script>

<div class="space-y-4">
  <PageCard
    title="New Wallet Transaction"
    description="Form pencatatan Kantong Duafa. Pembayaran tetap dilakukan di luar aplikasi, halaman ini hanya mencatat transaksi."
  >
    {#if data.accessDenied}
      <div class="empty-state">
        Role <strong>{data.user?.role}</strong> hanya punya akses baca. Pencatatan transaksi wallet tersedia untuk admin dan treasurer.
      </div>
    {:else}
      <form method="POST" class="space-y-4">
        {#if form?.message}
          <div class="helper-box-brand">
            <p class="helper-label text-sky-700">Error</p>
            <p class="mt-2 text-sm leading-6 text-slate-700">{form.message}</p>
          </div>
        {/if}

        <div class="grid gap-4 sm:grid-cols-2">
          <label>
            <span class="field-label">Type</span>
            <select name="type" class="input-field" required>
              <option value="income" selected={(form?.values?.type ?? 'income') === 'income'}>Income</option>
              <option value="expense" selected={(form?.values?.type ?? '') === 'expense'}>Expense</option>
            </select>
          </label>

          <label>
            <span class="field-label">Amount</span>
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
          <span class="field-label">Category</span>
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
          <span class="field-label">Description</span>
          <textarea
            class="input-field min-h-32 resize-y"
            name="description"
            placeholder="Catatan singkat transaksi"
            required
          >{form?.values?.description ?? ''}</textarea>
        </label>

        <div class="flex items-center justify-between gap-3">
          <a href="/wallet" class="btn-secondary px-4 py-3">Back to wallet</a>
          <button type="submit" class="btn-primary px-4 py-3">Save transaction</button>
        </div>
      </form>
    {/if}
  </PageCard>
</div>
