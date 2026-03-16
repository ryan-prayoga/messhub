<script lang="ts">
  export let disabled = false;
  export let threshold = 72;
  export let onRefresh: () => Promise<void>;

  let startY = 0;
  let pullDistance = 0;
  let pulling = false;
  let refreshing = false;

  function handleTouchStart(event: TouchEvent) {
    if (disabled || refreshing || event.touches.length !== 1 || window.scrollY > 0) {
      pulling = false;
      return;
    }

    startY = event.touches[0].clientY;
    pulling = true;
  }

  function handleTouchMove(event: TouchEvent) {
    if (!pulling || refreshing) {
      return;
    }

    const delta = event.touches[0].clientY - startY;
    if (delta <= 0) {
      pullDistance = 0;
      return;
    }

    pullDistance = Math.min(delta * 0.42, 96);
    if (pullDistance > 0) {
      event.preventDefault();
    }
  }

  async function handleTouchEnd() {
    if (!pulling) {
      return;
    }

    pulling = false;

    if (pullDistance < threshold) {
      pullDistance = 0;
      return;
    }

    refreshing = true;
    pullDistance = 56;

    try {
      await onRefresh?.();
    } finally {
      refreshing = false;
      pullDistance = 0;
    }
  }

  $: hint =
    refreshing
      ? 'Menyegarkan...'
      : pullDistance >= threshold
        ? 'Lepas untuk refresh'
        : 'Tarik ke bawah untuk refresh';
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="relative"
  on:touchstart={handleTouchStart}
  on:touchmove={handleTouchMove}
  on:touchend={handleTouchEnd}
  on:touchcancel={handleTouchEnd}
>
  <div
    class="pointer-events-none absolute inset-x-0 top-0 z-10 flex justify-center transition-opacity duration-200"
    style={`opacity: ${pullDistance > 0 || refreshing ? 1 : 0}; transform: translateY(${Math.max(
      0,
      pullDistance - 28
    )}px);`}
  >
    <div class="rounded-full border border-line bg-panel/95 px-4 py-2 text-xs font-semibold tracking-[0.12em] text-muted shadow-sm">
      {hint}
    </div>
  </div>

  <div
    class="transition-transform duration-200 ease-out"
    style={`transform: translateY(${pullDistance}px);`}
  >
    <slot />
  </div>
</div>
