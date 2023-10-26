<template>
  <Sequence
      v-if="repeats.length-1 > atRepeats"
      v-for="i in num"
      :repeats="repeats"
      :x="xPos(i)"
      :y="yPos(i)"
      :rotation="rot(i)"
      :atRepeats=atRepeats+1
      v-bind="repeats[atRepeats+1]"/>
  <TeamTable
      v-else
      v-for="i in num"
      :x="xPos(i) ?? 0"
      :y="yPos(i) ?? 0"
      :rotation="rot(i) ?? 10"
      :team-id="140+i" />
</template>

<script>

import TeamTable from "@/components/Layout/TeamTable.vue";

const SequenceType = {
  // Area: "area",
  Line: "Line",
  Circle: "Circle",
};

export default {
  components: {TeamTable},
  props: {
    'type': {type: String, required: false, default: SequenceType.Line},
    'x': {type: Number, required: false, default: 0},
    'y': {type: Number, required: false, default: 0},
    'rotation': {type: Number, required: false, default: 0},
    'radius': {type: Number, required: false, default: 100},
    'extra': {type: Number, required: false, default: 0},
    'num': {type: Number, required: false, default: 1},
    'axis': {type: Boolean, required: false, default: true},
    'dir': {type: Boolean, required: false, default: true},
    'separation': {type: Number, required: false, default: 50},
    'repeats': {type: Array, required: true},
    'atRepeats': {type: Number, required: true, default: 0},
    'equivalentSpaced': {type: Boolean, required: false, default: true}
  },
  methods: {
    xPos(i) {
      i--;
      if (this.type === SequenceType.Line) {
        return Math.round(this.x + this.dirVec.x * i);
      } else {
        if (i === 0) {
          return this.x;
        }

        const base = this.x + this.distVec(this.rotation, this.radius, false, this.dir).x

        const rad = (this.rotation+this.dirInt*90 + this.axisInt*this.trueSeparation * i)  * Math.PI/180;
        const offset = Math.cos(rad)*this.radius;

        return base + offset;
      }
    },
    yPos(i) {
      i--;
      if (this.type === SequenceType.Line) {
        return Math.round(this.y + this.dirVec.y * i);
      } else {
        if (i === 0) {
          return this.y;
        }
        const base = this.y + this.distVec(this.rotation, this.radius, false, this.dir).y

        const rad = (this.rotation+this.dirInt*90 + this.axisInt*this.trueSeparation * i)  * Math.PI/180;
        const offset = Math.sin(rad)*this.radius;

        return base + offset;
      }
    },
    rot(i) {
      i--;
      if (this.type === SequenceType.Line) {
        return this.rotation + this.extra;
      } else {
        return this.rotation+this.axisInt*this.trueSeparation*i+this.extra;
      }
    },

    distVec(rotation, sep, axis, dir) {
      let offset = 90;
      if (axis) {
        offset = 0;
      }

      const rad = (rotation+offset) * Math.PI/180;
      let x = Math.cos(rad)
      let y = Math.sin(rad)

      if (dir) {
        x = -x;
        y = -y;
      }

      const mag = Math.sqrt(x*x + y*y)

      return {x: sep*x/mag, y: sep*y/mag, rot: 0};
    }
  },
  computed: {
    dirVec() {
      return this.distVec(this.rotation, this.trueSeparation, this.axis, this.dir)
    },
    axisInt() {
      return this.axis ? 1 : -1;
    },
    dirInt() {
      return this.dir ? 1 : -1;
    },
    trueSeparation() {
      return (this.type !== SequenceType.Circle || !this.equivalentSpaced)
          ? this.separation
          : 360 / Math.max(1, this.num);
    }
  }
}
</script>

<style scoped>
</style>