<template>
  <g :transform="posCalc()" ref="slot" class="team-area">
    <rect
        class="element background"
        :x="settings.areaX"
        :y="settings.areaY"
        :width="settings.areaWidth"
        :height="settings.areaHeight"/>

    <rect
        class="element"
        :x="settings.tableX"
        :y="settings.tableY"
        :width="settings.tableWidth"
        :height="settings.tableHeight"

        style="overflow: visible;">
    </rect>

    <text
        dominant-baseline="central"
        font-family="sans-serif"
        :x="settings.tableX + settings.tableWidth/2"
        :y="settings.tableY + settings.tableHeight/2"
        :font-size="settings.dFontSize"
        font-style="normal"
        text-anchor="middle"
        font-weight="normal">
      {{ teamId }}
    </text>

    <rect
        class="element"
        v-for="index in settings.seatNum"
        :x="seatX(index-1)"
        :y="settings.seatY"
        :width="settings.seatWidth"
        :height="settings.seatHeight"/>

    <rect
        class="element outline"
        :x="settings.areaX"
        :y="settings.areaY"

        :width="settings.areaWidth"
        :height="settings.areaHeight"/>
    />
  </g>
  <Hatching
      :dist="settings.areaWidth/20"
      :stroke="settings.areaWidth/15"

  />

</template>

<style>

text {
  pointer-events: none;
}

.element {
  fill: none;
  stroke-width: v-bind(settings.strokeWidth);
  stroke: #444;
}

.background {
  fill: url(#diagonalHatch);
}


.element rect {
  stroke: #444;
  stroke-width: v-bind(settings.strokeWidth);
  stroke-dasharray: none;
  stroke-linecap: butt;
  stroke-dashoffset: 0;
  stroke-linejoin: miter;
  stroke-miterlimit: 4;
  fill-rule: nonzero;
  opacity: 1;

  fill: #fff;
}

.teamtable rect {
  stroke: #444;
  stroke-width: v-bind(settings.strokeWidth);
  stroke-dasharray: none;
  stroke-linecap: butt;
  stroke-dashoffset: 0;
  stroke-linejoin: miter;
  stroke-miterlimit: 4;
  fill-rule: nonzero;
  opacity: 1;

  fill: #fff;
}

.selectedteam rect {
  stroke-width: 2;
}

.teamtable.found rect {
  fill: orange;
}

.teamtable.found.exists rect {
  fill: lightgreen;
}

.teamtable.double rect {
  opacity: 0.2;
}

.teamtable.noteam rect {
  opacity: 0.3;
}

.teamtable:hover rect {
  stroke-width: v-bind(settings.strokeWidth);
}

rect.outline {
  fill: none;
  stroke-width: v-bind(settings.strokeWidth);
  stroke: #476cff;
}

</style>

<script setup>

// import {mapStores} from 'pinia'
import {teamareaStore} from "../../stores/teamarea";
import {defineProps, onMounted, ref} from "vue";
import {storeToRefs} from "pinia";
import Hatching from "./Hatching.vue";

const props = defineProps({
  'x': {type: Number, required: true},
  'y': {type: Number, required: true},
  'rotation': {type: Number, required: false, default: 0},
  'teamId': {type: String, required: false},
  'relevant-room-element': {type: Array, required: true},
})

let group = ref(null)

function posCalc() {
  return `translate(${props.x}, ${props.y}) rotate(${props.rotation})`
}

function seatX(i) {
  return settings.tableX + settings.seatPadding
      + i * (settings.seatWidth + settings.seatSep + settings.seatPadding);
}

const settings = teamareaStore()
const {strokeWidth} = storeToRefs(settings)

onMounted(() => {
  settings.registerTeamId(props.teamId)
})

</script>
