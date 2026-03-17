<script lang="ts">
  import { createEventDispatcher, tick } from 'svelte';
  import AppIcon from '$lib/components/AppIcon.svelte';

  export let open = false;
  export let title: string;
  export let description = '';
  export let icon: string | null = null;
  export let ariaLabel: string | null = null;
  export let variant: 'dialog' | 'sheet' = 'dialog';
  export let closeOnEscape = true;
  export let closeOnBackdrop = true;
  export let dismissible = true;
  export let initialFocusSelector = '';

  const dispatch = createEventDispatcher<{
    close: void;
  }>();

  const titleId = `dialog-title-${Math.random().toString(36).slice(2, 10)}`;
  const descriptionId = `dialog-description-${Math.random().toString(36).slice(2, 10)}`;

  let panelElement: HTMLDivElement | null = null;
  let previouslyFocusedElement: HTMLElement | null = null;
  let lastOpen = false;

  function requestClose() {
    if (!dismissible) {
      return;
    }

    dispatch('close');
  }

  function getFocusableElements() {
    if (!panelElement) {
      return [];
    }

    return Array.from(
      panelElement.querySelectorAll<HTMLElement>(
        'a[href], button:not([disabled]), textarea:not([disabled]), input:not([disabled]), select:not([disabled]), [tabindex]:not([tabindex="-1"])'
      )
    ).filter((element) => !element.hasAttribute('disabled') && element.getAttribute('aria-hidden') !== 'true');
  }

  async function focusInitialElement() {
    await tick();

    if (!panelElement) {
      return;
    }

    const explicitTarget =
      initialFocusSelector !== '' ? panelElement.querySelector<HTMLElement>(initialFocusSelector) : null;
    const target = explicitTarget ?? getFocusableElements()[0] ?? panelElement;
    target.focus();
  }

  function restoreFocus() {
    previouslyFocusedElement?.focus();
    previouslyFocusedElement = null;
  }

  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget && closeOnBackdrop) {
      requestClose();
    }
  }

  function handleWindowKeydown(event: KeyboardEvent) {
    if (open && closeOnEscape && event.key === 'Escape') {
      event.preventDefault();
      requestClose();
    }
  }

  function handlePanelKeydown(event: KeyboardEvent) {
    if (event.key !== 'Tab') {
      return;
    }

    const focusable = getFocusableElements();
    if (focusable.length === 0) {
      event.preventDefault();
      panelElement?.focus();
      return;
    }

    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    const current = document.activeElement as HTMLElement | null;

    if (event.shiftKey && current === first) {
      event.preventDefault();
      last.focus();
      return;
    }

    if (!event.shiftKey && current === last) {
      event.preventDefault();
      first.focus();
    }
  }

  $: if (open && !lastOpen) {
    previouslyFocusedElement = document.activeElement instanceof HTMLElement ? document.activeElement : null;
    void focusInitialElement();
  }

  $: if (!open && lastOpen) {
    restoreFocus();
  }

  $: lastOpen = open;
</script>

<svelte:window on:keydown={handleWindowKeydown} />

{#if open}
  <div
    class={`modal-backdrop ${variant === 'sheet' ? 'modal-backdrop-sheet' : 'modal-backdrop-dialog'}`}
    role="presentation"
    tabindex="-1"
    on:click={handleBackdropClick}
  >
    <div
      bind:this={panelElement}
      class={`modal-panel ${variant === 'sheet' ? 'modal-panel-sheet' : 'modal-panel-dialog'}`}
      role="dialog"
      aria-modal="true"
      aria-label={ariaLabel ?? undefined}
      aria-labelledby={ariaLabel ? undefined : titleId}
      aria-describedby={description ? descriptionId : undefined}
      tabindex="-1"
      on:keydown={handlePanelKeydown}
    >
      {#if variant === 'sheet'}
        <div class="modal-sheet-grabber"></div>
      {/if}

      <header class="modal-header">
        <div class="modal-header-copy">
          {#if icon}
            <div class="dialog-icon" aria-hidden="true">
              <AppIcon {icon} className="h-5 w-5" />
            </div>
          {/if}

          <div class="min-w-0">
            <h2 id={titleId} class="section-title text-[1.2rem] sm:text-[1.35rem]">{title}</h2>
            {#if description}
              <p id={descriptionId} class="section-subtitle mt-2">{description}</p>
            {/if}
          </div>
        </div>

        {#if dismissible}
          <button type="button" class="icon-button" aria-label="Tutup dialog" on:click={requestClose}>
            <AppIcon icon="lucide:x" className="h-5 w-5" />
          </button>
        {/if}
      </header>

      <div class="modal-body">
        <slot />
      </div>

      {#if $$slots.footer}
        <footer class="modal-footer">
          <slot name="footer" />
        </footer>
      {/if}
    </div>
  </div>
{/if}
