<script lang="ts">
  import AppIcon from '$lib/components/AppIcon.svelte';
  import Spinner from '$lib/components/Spinner.svelte';

  export let tone: 'loading' | 'error' | 'empty' | 'forbidden' = 'empty';
  export let title = '';
  export let message: string;
  export let requestId: string | null = null;

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
</script>

<div class={classes[tone]}>
  {#if tone === 'loading' || tone === 'error' || tone === 'forbidden'}
    <div class="feedback-banner-icon">
      {#if tone === 'loading'}
        <Spinner className="h-5 w-5" />
      {:else}
        <AppIcon icon={icons[tone]} className="h-5 w-5" />
      {/if}
    </div>
  {/if}

  <div class="min-w-0">
    {#if title}
      <p class={tone === 'empty' ? 'helper-label' : 'feedback-banner-title'}>{title}</p>
    {/if}

    <p class={tone === 'empty' ? 'mt-2 text-sm leading-6 text-muted' : 'feedback-banner-message'}>
      {message}
    </p>
  </div>
</div>
