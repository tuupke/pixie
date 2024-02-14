<template>
  <g :transform="posCalc()" ref="slot" class="team-area">
    <rect
        class="background"
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
        dominant-baseline="middle"
        alignment-baseline="central"
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
        class="outline"
        :x="settings.areaX"
        :y="settings.areaY"

        :width="settings.areaWidth"
        :height="settings.areaHeight"/>
    />
  </g>
</template>

<style>

.background {
  fill: none;
}

rect.outline {
  fill: none;
  stroke-width: 1px;
  stroke: #476cff;
}

.team-area:hover .outline {
  stroke-width: 3px;
}

text {
  pointer-events: none;
}

rect.element {
  fill: #fff;
  stroke: #444;
}

.area .element {
  stroke: #444;
  stroke-width: 1;
  stroke-dasharray: none;
  stroke-linecap: butt;
  stroke-dashoffset: 0;
  stroke-linejoin: miter;
  stroke-miterlimit: 4;
  fill-rule: nonzero;
  opacity: 1;

  fill: #fff;
}

.area:hover, .area .selected {
  stroke-width: 2;
}

.noteam {

}


.teamtable rect {
  stroke: #444;
  stroke-width: 1;
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
  stroke-width: 3px;
}

</style>

<script setup>

// import {mapStores} from 'pinia'
import {teamareaStore} from "@/stores/teamarea";
import {onMounted, defineProps, ref} from "vue";

const props = defineProps({
    'x': {type: Number, required: true},
    'y': {type: Number, required: true},
    'rotation': {type: Number, required: false, default: 0},
    'teamId': {type: String, required: false},
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

onMounted(() => {
  settings.registerTeamId(props.teamId)
})


</script>
