<template>
    <!--    <polyline :points="outline" fill="none"/>-->
    <g ref="slot" transform="">
      <slot></slot>
    </g>
    <circle
        :belongsTo="slot"
        r="5"
        style="z-index: 1000"
        fill="blue"
        :cx="tl.x"
        :cy="tl.y"
    />
</template>

<script setup>

import {computed, onMounted, reactive, ref, toRaw} from 'vue'

let dragging = false
let rotating = false;
let slot = ref(null);

const boundingRect = reactive({
  Top: 0,
  Right: 0,
  Bottom: 0,
  Left: 0,
});

// const outline = computed(() => {
//   return (boundingRect.Left + "," + boundingRect.Top + " " +
//       boundingRect.Right + "," + boundingRect.Top + " " +
//       boundingRect.Right + "," + boundingRect.Bottom + " " +
//       boundingRect.Left + "," + boundingRect.Bottom + " " +
//       boundingRect.Left + "," + boundingRect.Top);
// })

const tl = reactive({
  x: 25,
  y: 55,
})

defineProps({
  x: Number,
  y: Number,
})

onMounted(() => {
  // const root = slot.value.childNodes[1];

  // console.log(root)

  // Find parent svg
  // let svg = root;
  // while (svg && svg.nodeName !== 'svg') {
  //   svg = svg.parentNode;
  // }
  //
  // if (!svg) {
  //   console.log("Draggable failed", root)
  //   return;
  // }
  //
  // const svg_bb = svg.getBoundingClientRect();
  const opts = {
    fill: true,
    stroke: true,
    markers: true,
    clipped: true,
  };

  const el_bb = slot.value.getBBox(opts);
  console.log(slot.value.getBBox(opts), slot.value.childNodes[1].getBBox(opts))
  console.log(slot.value.getBoundingClientRect(), slot.value.childNodes[1].getBoundingClientRect())
  console.log(slot.value.getCTM(), slot.value.childNodes[1].getCTM())
  tl.x = slot.value.getBoundingClientRect().x
  tl.y = slot.value.getBoundingClientRect().y
  // console.log(slot.value.childNodes[1])

  // boundingRect.Left = el_bb.left - svg_bb.left;
  // boundingRect.Right = el_bb.right - svg_bb.left;
  // boundingRect.Top = el_bb.top - svg_bb.top;
  // boundingRect.Bottom = el_bb.bottom - svg_bb.top;
  // console.log(toRaw(boundingRect))
  // console.log(el_bb);
})

</script>

<style scoped>

</style>