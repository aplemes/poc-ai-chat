<script setup lang="ts">
interface Props {
  open: boolean
  hasNewMessage: boolean
  label: string
}

defineProps<Props>()
defineEmits<{ click: [] }>()
</script>

<template>
  <button class="fab" :class="{ open }" :aria-label="label" @click="$emit('click')">
    <Transition name="icon" mode="out-in">
      <svg v-if="!open" key="open" viewBox="0 0 24 24" fill="currentColor">
        <path
          d="M4.913 2.658c2.075-.27 4.19-.408 6.337-.408 2.147 0 4.262.139 6.337.408 1.922.25 3.291 1.861 3.405 3.727a4.403 4.403 0 00-1.032-.211 50.89 50.89 0 00-8.42 0c-2.358.196-4.04 2.19-4.04 4.434v4.286a4.47 4.47 0 002.433 3.984L7.28 21.53A.75.75 0 016 21v-4.03a48.527 48.527 0 01-1.087-.128C2.905 16.58 1.5 14.833 1.5 12.862V6.638c0-1.97 1.405-3.718 3.413-3.979z"
        />
        <path
          d="M15.75 7.5c-1.376 0-2.739.057-4.086.169C10.124 7.797 9 9.103 9 10.609v4.285c0 1.507 1.128 2.814 2.67 2.94 1.243.102 2.5.157 3.768.165l2.782 2.781a.75.75 0 001.28-.53v-2.39l.33-.026c1.542-.125 2.67-1.433 2.67-2.94v-4.286c0-1.505-1.125-2.811-2.664-2.94A49.392 49.392 0 0015.75 7.5z"
        />
      </svg>
      <svg
        v-else
        key="close"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.5"
        stroke-linecap="round"
      >
        <path d="M18 6L6 18M6 6l12 12" />
      </svg>
    </Transition>
    <span v-if="hasNewMessage" class="fab-badge" aria-hidden="true"></span>
  </button>
</template>

<style scoped>
.fab {
  position: relative;
  width: 3.25rem;
  height: 3.25rem;
  border-radius: 50%;
  background: var(--gradient-primary);
  color: #fff;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-xl), var(--shadow-primary);
  transition: all var(--transition-normal);
  flex-shrink: 0;
}

.fab:hover {
  transform: scale(1.07);
  box-shadow:
    0 12px 32px rgba(0, 0, 0, 0.15),
    0 6px 16px rgba(0, 135, 74, 0.4);
}
.fab.open {
  box-shadow: var(--shadow-md);
}

.fab:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 3px;
}

.fab svg {
  width: 1.35rem;
  height: 1.35rem;
}

.fab-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #ef4444;
  border: 2px solid #fff;
  animation: badge-pop 0.3s var(--ease-out);
}

@keyframes badge-pop {
  from {
    transform: scale(0);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.icon-enter-active,
.icon-leave-active {
  transition:
    opacity var(--transition-fast),
    transform var(--transition-fast);
}
.icon-enter-from,
.icon-leave-to {
  opacity: 0;
  transform: rotate(25deg) scale(0.75);
}

@media (prefers-reduced-motion: reduce) {
  .fab,
  .fab-badge,
  .icon-enter-active,
  .icon-leave-active {
    transition: none;
    animation: none;
  }
}
</style>
