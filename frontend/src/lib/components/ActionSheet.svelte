<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import AppIcon from '$lib/components/AppIcon.svelte';

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

  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      requestClose();
    }
  }

  function handleBackdropKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      requestClose();
    }
  }

  function handleWindowKeydown(event: KeyboardEvent) {
    if (open && event.key === 'Escape') {
      requestClose();
    }
  }
</script>

<svelte:window on:keydown={handleWindowKeydown} />

{#if open}
  <div
    class="modal-backdrop"
    role="presentation"
    tabindex="-1"
    on:click={handleBackdropClick}
    on:keydown={handleBackdropKeydown}
  >
    <div class="modal-sheet" role="dialog" aria-modal="true" aria-label={title}>
      <div class="modal-sheet-grabber"></div>

      <header class="modal-sheet-header">
        <div class="flex min-w-0 items-start gap-4">
          <div class="nav-link-icon mt-0.5">
            <AppIcon {icon} className="h-5 w-5" />
          </div>

          <div class="min-w-0">
            <h2 class="section-title text-[1.35rem]">{title}</h2>
            {#if description}
              <p class="section-subtitle mt-2">{description}</p>
            {/if}
          </div>
        </div>

        <button type="button" class="icon-button" aria-label="Tutup panel" on:click={requestClose}>
          <AppIcon icon="lucide:x" className="h-5 w-5" />
        </button>
      </header>

      <div class="modal-sheet-body">
        <slot />
      </div>
    </div>
  </div>
{/if}
