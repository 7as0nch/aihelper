<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue';

const props = withDefaults(
  defineProps<{
    /** Stagger `.scroll-reveal-stagger-item` children (queue-anim “bottom”) */
    stagger?: boolean;
    threshold?: number;
  }>(),
  { stagger: false, threshold: 0.12 },
);

const root = ref<HTMLElement | null>(null);
const visible = ref(false);
let observer: IntersectionObserver | null = null;

onMounted(() => {
  const el = root.value;
  if (!el) {
    return;
  }
  observer = new IntersectionObserver(
    (entries) => {
      const hit = entries.some((e) => e.isIntersecting);
      if (hit) {
        visible.value = true;
        observer?.disconnect();
        observer = null;
      }
    },
    { threshold: props.threshold, rootMargin: '0px 0px -6% 0px' },
  );
  observer.observe(el);
});

onBeforeUnmount(() => {
  observer?.disconnect();
  observer = null;
});
</script>

<template>
  <div
    ref="root"
    class="scroll-reveal"
    :class="{
      'scroll-reveal--visible': visible,
      'scroll-reveal--stagger': stagger,
    }"
  >
    <slot />
  </div>
</template>

<style scoped>
/* Whole block: OverPack-style fade + slide */
.scroll-reveal:not(.scroll-reveal--stagger) {
  opacity: 0;
  transform: translate3d(0, 40px, 0);
  transition:
    opacity 0.55s cubic-bezier(0.25, 0.46, 0.45, 0.94),
    transform 0.55s cubic-bezier(0.25, 0.46, 0.45, 0.94);
  will-change: opacity, transform;
}

.scroll-reveal--visible:not(.scroll-reveal--stagger) {
  opacity: 1;
  transform: translate3d(0, 0, 0);
}

/* Stagger: container is neutral; only items animate */
.scroll-reveal--stagger {
  opacity: 1;
}

.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item) {
  opacity: 0;
  transform: translate3d(0, 32px, 0);
  transition:
    opacity 0.48s cubic-bezier(0.25, 0.46, 0.45, 0.94),
    transform 0.48s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item) {
  opacity: 1;
  transform: translate3d(0, 0, 0);
}

.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(1)) {
  transition-delay: 0.04s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(2)) {
  transition-delay: 0.1s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(3)) {
  transition-delay: 0.16s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(4)) {
  transition-delay: 0.22s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(5)) {
  transition-delay: 0.28s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(6)) {
  transition-delay: 0.34s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(7)) {
  transition-delay: 0.4s;
}
.scroll-reveal--visible.scroll-reveal--stagger :deep(.scroll-reveal-stagger-item:nth-child(8)) {
  transition-delay: 0.46s;
}
</style>
