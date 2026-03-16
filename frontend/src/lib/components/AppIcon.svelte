<script lang="ts">
  import Icon from '@iconify/svelte';

export let icon: string;
export let className = '';
export let title: string | undefined = undefined;
</script>

<Icon {icon} class={className} aria-hidden={title ? undefined : 'true'} aria-label={title} />
