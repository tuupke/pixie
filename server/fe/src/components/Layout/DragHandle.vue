<template>
  <PlacedCrossHairs
      :x=x
      :y=y
      :highlightable="true"
      @mouseover.stop="console.log('over')"
      :scale-multiply=0.7
      color="orange"/>
</template>

<script setup lang="ts">
import {CoordinateInterface, RotationCoordinateInterface} from "../../types.ts"
import PlacedCrossHairs from "../components/Layout/PlacedCrossHairs.vue";

const props = defineProps<CoordinateInterface & { rotateCoord: RotationCoordinateInterface}>()
</script>
