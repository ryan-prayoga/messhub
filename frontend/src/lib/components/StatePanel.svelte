<script lang="ts">
  import AppIcon from '$lib/components/AppIcon.svelte';
  import Spinner from '$lib/components/Spinner.svelte';

  export let tone: 'loading' | 'error' | 'empty' | 'forbidden' = 'empty';
  export let title = '';
  export let message: string;
  export let requestId: string | null = null;
  export let icon: string | null = null;
  export let actionHref: string | null = null;
  export let actionLabel = '';
  export let actionVariant: 'primary' | 'secondary' = 'secondary';

  $: requestId;

  const classes = {
    loading: 'feedback-banner feedback-banner-info',
    error: 'feedback-banner feedback-banner-error',
    empty: 'empty-state',
    forbidden: 'feedback-banner feedback-banner-warning'
  } as const;

  const icons = {
    loading: 'loading',
    error: 'lucide:circle-alert',
    empty: 'lucide:inbox',
    forbidden: 'lucide:shield-alert'
  } as const;

  $: resolvedIcon = icon ?? icons[tone];
</script>

<div class={classes[tone]}>
  {#if tone === 'empty'}
    <div class="state-panel-empty-icon">
      <AppIcon icon={resolvedIcon} className="h-6 w-6" />
    </div>
  {:else if tone === 'loading' || tone === 'error' || tone === 'forbidden'}
    <div class="feedback-banner-icon">
      {#if tone === 'loading'}
        <Spinner className="h-5 w-5" />
      {:else}
        <AppIcon icon={resolvedIcon} className="h-5 w-5" />
      {/if}
    </div>
  {/if}

  <div class="min-w-0">
    {#if title}
      <p class={tone === 'empty' ? 'state-panel-empty-title' : 'feedback-banner-title'}>{title}</p>
    {/if}

    <p class={tone === 'empty' ? 'state-panel-empty-copy' : 'feedback-banner-message'}>
      {message}
    </p>

    {#if requestId && tone !== 'empty'}
      <p class="mt-2 text-xs text-muted">Request ID: {requestId}</p>
    {/if}

    {#if actionHref && actionLabel}
      <div class="state-panel-actions">
        <a href={actionHref} class={actionVariant === 'primary' ? 'btn-primary px-4 py-3' : 'btn-secondary px-4 py-3'}>
          {actionLabel}
        </a>
      </div>
    {/if}

    <slot />
  </div>
</div>
