<template>
  <Button
      v-bind="$attrs"
      :outlined="!confirming"
      v-on:click.stop="maybeConfirm"
  />
</template>
<script setup lang="ts">

import {defineEmits, onUnmounted, ref} from "vue";

const emit = defineEmits<{ confirmed: [] }>()
const confirming = ref<number>(0)

function maybeConfirm() {
  if (!confirming.value) {
    confirming.value = window.setTimeout(clearConfirm, 1000)
    return
  }

  clearConfirm(true)
  emit('confirmed')
}

function clearConfirm(andTimer: boolean = false) {
  if (andTimer) {
    window.clearTimeout(confirming.value)
  }
  confirming.value = 0
}

onUnmounted(() => clearConfirm(true))

</script>
