<template>
  <g :transform="posCalc()" class="team-area">
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
        class="element"
        :x="settings.areaX"
        :y="settings.areaY"

        :width="settings.areaWidth"
        :height="settings.areaHeight"/>
    />
  </g>

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
  fill: v-bind(fill);
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
  opacity: 0.1;

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

<script setup lang="ts">

import {teamareaStore} from "../../stores/teamarea";
import {computed, onMounted, ref} from "vue";
import {RotationCoordinateInterface} from "../../types.ts";

interface StatusBooleans{
  team?: boolean
  host?: boolean
  used?: boolean
  seen?: boolean
}

const props = withDefaults(defineProps<RotationCoordinateInterface & StatusBooleans & { 'teamId': string }>(), {
  rotation: 0
});


const fill = computed(() => {
  if ((props.team && props.host) || !props.used) {
    return "white"
  }

  // url(#selectedHatching)

  return "orange"
});


function posCalc(): string {
  return `translate(${props.x}, ${props.y}) rotate(${props.rotation})`
}

function seatX(i: number): number {
  return settings.tableX + settings.seatPadding
      + i * (settings.seatWidth + settings.seatSep + settings.seatPadding);
}

const settings = teamareaStore()
// const {strokeWidth} = storeToRefs(settings)

onMounted(() => {
  settings.registerTeamId(props.teamId)
})

</script>
