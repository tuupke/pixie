<template>
  <svg style="overflow: visible" class="room">
    <rect v-if="this.shape === 'rectangle'" :x="this.x" :y="this.y" :width="this.width" :height="this.height" />
    <ellipse v-else-if="this.shape === 'circular'" :rx="this.width/2" :ry="this.height/2" :cx="this.x+(this.width/2)" :cy="this.y+(this.height/2)"/>
    <polygon v-else-if="this.shape === 'free'"  points="200,10 250,190 160,210" style="fill:none;stroke:purple;stroke-width:1" />
    <text v-else>Unsupported type {{ this.shape }}</text>
  </svg>
</template>

<style>

svg.room rect, svg.room ellipse, svg.room polygon {
  fill: #eee;
  stroke: #444;
  stroke-width: 1;
}

</style>

<script>

export default {
  props: {
    'shape': {type: String, required: true, default: ""},
    'width': {type: Number, required: false, default: 100},
    'height': {type: Number, required: false, default: 100},
    'x': {type: Number, required: true},
    'y': {type: Number, required: true},
  },
}
</script>