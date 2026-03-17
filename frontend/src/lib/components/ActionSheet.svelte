<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import ModalDialog from '$lib/components/ModalDialog.svelte';

  export let open = false;
  export let title: string;
  export let description = '';
  export let icon = 'lucide:panel-right-open';

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  function requestClose() {
    dispatch('close');
  }
</script>

{#if open}
  <ModalDialog
    {open}
    {title}
    {description}
    {icon}
    variant="sheet"
    on:close={requestClose}
  >
    <slot />
  </ModalDialog>
{/if}
