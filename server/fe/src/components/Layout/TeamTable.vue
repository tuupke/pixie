<template>
  <g :transform="posCalc()" class="team-area">
    <rect
        class="element"
        style="fill: white"
        :x="settings.areaX"
        :y="settings.areaY"
        :width="settings.areaWidth"
        :height="settings.areaHeight"/>

    <rect
        class="background element"
        :x="settings.areaX"
        :y="settings.areaY"
        :width="settings.areaWidth"
        :height="settings.areaHeight"/>

    <rect
        class="element"
        :x="settings.tableX"
        :y="settings.tableY"
        :width="settings.tableWidth"
        :height="settings.tableHeight"/>

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
  stroke-width: v-bind(outline);
  opacity: v-bind(opacity);
  stroke: #000;
}

.background {
  fill: v-bind(fill);
  opacity: v-bind(backgroundOpacity);
}

.element rect {
  stroke: #000;
  stroke-width: v-bind(outline);
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
  stroke: #000;
  stroke-width: v-bind(outline);
  stroke-dasharray: none;
  stroke-linecap: butt;
  stroke-dashoffset: 0;
  stroke-linejoin: miter;
  stroke-miterlimit: 4;
  fill-rule: nonzero;
}

.selectedteam rect {
  stroke-width: 2;
}

.teamtable:hover rect {
  stroke-width: v-bind(outline);
}

</style>

<script setup lang="ts">

import {teamareaStore} from "../../stores/teamarea";
import {computed, inject} from "vue";
import {RotationCoordinateInterface} from "../../types.ts";

// Three use cases, and multiple failure-modes, when showing tables exist.
// 1. Constructing the map (handled in EditableTeamTable)
//    1.1 The table is not attached to a path. -> lhs blue
//    1.2 The assigned table number is not unique. -> lhs orange
// 2. Host and team assignments
//    2.1 The table is not attached to a team. -> lhs light blue
//    2.2 The table is not attached to a host. -> rhs light orange
//    2.3 A single team assigned to multiple tables. -> lhs dark blue
//    2.4 A single host assigned to multiple tables. -> rhs dark orange
//    2.5 Host attached, but not seen.
// 3. Printing
//    No failure modes, colored team tables are not acceptable.


const props = withDefaults(defineProps<RotationCoordinateInterface & {
  teamId: string
  fill?: string
  hidden? : boolean
  highlighted? : boolean
}>(), {
  rotation: 0,
  hidden: false,
  highlighted: false,
  fill: 'white',
});

const outline = computed(() => (props.highlighted ? 5 : 1) * settings.strokeWidth)

const editing = inject<boolean>("editing")! ?? false

const opacity = computed<number>(() => {
  return props.hidden ? (editing ? 0.1 : 0) : 1
})

const backgroundOpacity = computed<number>(() => {
  return props.hidden ? (editing ? 0 : 0) : 0.6
})

function posCalc(): string {
  return `translate(${props.x}, ${props.y}) rotate(${props.rotation})`
}

function seatX(i: number): number {
  return settings.tableX + settings.seatPadding + i * (settings.seatWidth + settings.seatSep + settings.seatPadding);
}

const settings = teamareaStore()

</script>
