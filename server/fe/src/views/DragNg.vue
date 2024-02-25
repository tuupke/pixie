<template>
  <g class="dragNg"
     @mousedown.stop="(e: MouseEvent) => $emit('dragStart', e)"
     @mouseover.stop="(e: MouseEvent) => $emit('maybeShow', e)"
     @click.stop="(e: MouseEvent) => $emit('maybeSelect', e)"
  >
  <slot></slot>
  </g>
</template>

<style scoped>
.dragNg:hover {
  cursor: move;
}

</style>

<script setup lang="ts">

const emits = defineEmits(['dragStart', 'maybeShow', 'maybeSelect'])
</script>
