<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import ActionButtonGroup from '$lib/components/ActionButtonGroup.svelte';
  import ModalDialog from '$lib/components/ModalDialog.svelte';

  export let open = false;
  export let title: string;
  export let description = '';
  export let confirmLabel = 'Konfirmasi';
  export let cancelLabel = 'Batal';
  export let destructive = false;
  export let loading = false;
  export let icon: string | null = null;

  const dispatch = createEventDispatcher<{
    close: void;
    confirm: void;
  }>();

  function handleClose() {
    if (loading) {
      return;
    }

    dispatch('close');
  }

  function handleConfirm() {
    if (loading) {
      return;
    }

    dispatch('confirm');
  }

  $: resolvedIcon = icon ?? (destructive ? 'lucide:triangle-alert' : 'lucide:badge-alert');
</script>

<ModalDialog
  {open}
  {title}
  {description}
  icon={resolvedIcon}
  variant="dialog"
  dismissible={!loading}
  closeOnEscape={!loading}
  closeOnBackdrop={!loading}
  initialFocusSelector="[data-dialog-cancel]"
  on:close={handleClose}
>
  <slot />

  <svelte:fragment slot="footer">
    <ActionButtonGroup align="end">
      <button
        type="button"
        class="btn-secondary min-w-[7.5rem]"
        data-dialog-cancel
        disabled={loading}
        on:click={handleClose}
      >
        {cancelLabel}
      </button>
      <button
        type="button"
        class={`${destructive ? 'btn-danger' : 'btn-primary'} min-w-[8.5rem]`}
        disabled={loading}
        on:click={handleConfirm}
      >
        {loading ? 'Memproses...' : confirmLabel}
      </button>
    </ActionButtonGroup>
  </svelte:fragment>
</ModalDialog>
