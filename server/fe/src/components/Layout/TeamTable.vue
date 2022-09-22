<template>
  <g :transform="posCalc()" class="teamtable">

    <rect :x="-areaWidth/2" :y="-areaHeight/2" style='stroke-width: 1px; stroke: #0b7ad1; fill: none;' :width="areaWidth" :height="areaHeight"></rect>
    <rect :x="-width/2" :y="-2*height/3" :width="width" :height="height"/>
    <text alignment-baseline="middle" font-family="sans-serif" :font-size="height" font-style="normal"
          text-anchor="middle"
          font-weight="normal">
      <slot></slot>
    </text>
    <g>
      <rect v-for="(n,index) in seatNum" :x="seatX(index)" :y="height / 3 + seatDist" :width="seatWidth" :height="seatHeight"/>
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
  stroke-width: 3;
}

</style>

<script>

export default {
  props: {
    'x': {type: Number, required: true},
    'y': {type: Number, required: true},

    'seatSep': {type: Number, required: false, default: 3},
    'seatDist': {type: Number, required: false, default: 3},
    'seatNum': {type: Number, required: false, default: 3},

    'margin': {type: Number, required: false, default: 0},

    'areaWidth': {type: Number, required: false, default: 50},
    'areaHeight': {type: Number, required: false, default: 30},

    'height': {type: Number, required: false, default: 15},
    'seatHeight': {type: Number, required: false, default: 3},
    'rotation': {type: Number, required: false, default: 0},
  },
  computed: {
    width() {
      return this.areaWidth - 2*this.margin;
    },
    seatWidth() {
      return this.width / this.seatNum - this.seatSep;
    },
    totalHeight() {
      return this.areaHeight - 2*this.margin;
    }
  },
  methods: {
    posCalc: function () {
      return `translate(${this.x}, ${this.y}) rotate(${this.rotation})`
    },
    seatX: function(i) {
      let x = -this.width / 2 + this.seatSep / 2;
      // if (this.seatNum % 2 !== 0) {
      //   // Odd, shift everything one over
      //   x += this.seatWidth / 4;
      // }

      x += i * (this.seatWidth + this.seatSep)

      return x
    }
  },
}
</script>