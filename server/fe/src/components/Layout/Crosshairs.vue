<template>
  <g :transform="'scale('+scale+')'">
  <g :transform="'translate(0, '+scale*-11+') rotate(45)'">

    <svg xmlns="http://www.w3.org/2000/svg">
      <circle class="stroked" fill="white" stroke="blue" :cx="8" :cy="8" :r="6"/>
      <path class="stroked" fill="none" stroke="blue" d="M 8 0 L 8 6.5"/>
      <path class="stroked" fill="none" stroke="blue" d="M 0 8 L 6.5 8"/>
      <path class="stroked" fill="none" stroke="blue" d="M 8 9.5 L 8 16"/>
      <path class="stroked" fill="none" stroke="blue" d="M 9.5 8 L 16 8"/>
    </svg>
  </g>
  </g>
</template>

<script setup lang="ts">

withDefaults(defineProps<{
  scale: number
}>(), {
  scale: 1
})


</script>
