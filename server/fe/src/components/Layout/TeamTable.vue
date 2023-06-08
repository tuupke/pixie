<template>
  <g :transform="posCalc()" class="teamtable">
    <rect :x="this.teamtableStore.areaX" :y="this.teamtableStore.areaY" style='stroke-width: 1px; stroke: #0b7ad1; fill: none;' :width="this.teamtableStore.areaWidth" :height="this.teamtableStore.areaHeight"></rect>
    <svg :x="this.teamtableStore.tableX" :y="this.teamtableStore.tableY" :width="this.teamtableStore.tableWidth" :height="this.teamtableStore.tableHeight" style="overflow: visible">

      <rect :width="this.teamtableStore.tableWidth" :height="this.teamtableStore.tableHeight" />
      <text
          x="50%" y="50%"
          dominant-baseline="middle"
          alignment-baseline="central"
          font-family="sans-serif"
          :font-size="this.teamtableStore.dFontSize"
          font-style="normal"
          text-anchor="middle"
          font-weight="normal">
        {{ teamId }}
      </text>
    </svg>
    <g>
      <rect v-for="(n,index) in this.teamtableStore.seatNum" :x="seatX(index)" :y="this.teamtableStore.seatY" :width="this.teamtableStore.seatWidth" :height="this.teamtableStore.seatHeight"/>
    </g>
  </g>
</template>

<style>
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

<script>
import {mapStores} from 'pinia'
import {teamtableStore} from "@/stores/teamtable";

export default {
  props: {
    'x': {type: Number, required: true},
    'y': {type: Number, required: true},
    'rotation': {type: Number, required: false, default: 0},
    'teamId': {type: Number, required: true},
  },

  computed: mapStores(teamtableStore),
  methods: {
    posCalc: function () {
      return `translate(${this.x}, ${this.y}) rotate(${this.rotation})`
    },
    seatX: function(i) {
      return -this.teamtableStore.tableWidth / 2 + this.teamtableStore.seatSep / 2 + i * (this.teamtableStore.seatWidth + this.teamtableStore.seatSep) + this.teamtableStore.offsetX;
    }
  },

  mounted() {
    this.teamtableStore.registerTeamId(this.teamId)
  }
}
</script>