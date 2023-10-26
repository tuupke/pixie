<template>
  <g :a=0 >
    <slot></slot>
  </g>
</template>

<style scoped>

</style>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
  name: "Rotatable"
})
</script>

