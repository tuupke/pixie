<template>
  <div class="flex-1 flex-column">

  </div>
  <div class="grid my-5 mx-1">
    <div></div>
    <div class="my-3 flex flex-column align-content-end" style="grid-column: span 2">
        <span class="py-2" :style="'margin-right: '+2*scale +'px;'">
            <Slider :tooltips=false :lazy=false v-model="settings.areaOffsetX"
                    :min=0 :max=100 :step=sliderStep
                    :style="'width: ' +scale*settings.areaWidth+ 'px;'"
                    class='slider-blue'/>
        </span>
      <span class="py-1">
            <Slider v-model=tableX :tooltips=false :lazy=false
                    :min=0 :max=settings.areaWidth :step=sliderStep
                    :style="'width: ' +scale*settings.areaWidth+ 'px;'"
                    class="slider-blue"/>
      </span>
    </div>

    <div class="flex flex-row justify-content-end flex-wrap">
      <div class="px-1">
        <Slider :tooltips=false :lazy=false v-model="settings.areaOffsetY"
                :min=0 :max=100 orientation="vertical" :step=sliderStep
                :style="'height: ' +scale*settings.areaHeight+ 'px;'"
                class="slider-blue"/>
      </div>
      <div class="px-2" :style="'margin-bottom: '+2*scale +'px;'">
        <Slider v-model=tableY :tooltips=false :lazy=false
                orientation="vertical" :step=sliderStep
                :min=0 :max=settings.areaHeight
                :style="'height: ' +scale*settings.areaHeight+ 'px;'"
                class="slider-blue"/>
      </div>
    </div>
    <div ref="svgDiv">
      <svg width="100%" height="100%" class="" style="margin-bottom: 5rem;">
        <g :transform="'scale('+scale+') translate('+(-settings.areaX)+','+(-settings.areaY)+')'">
          <TeamTable :x="0.5" :rotation="0" :y="0.5" team-id="100"/>
          <circle :cx=0.5 :cy=0.5 r="3" fill="orange"/>
        <g transform="translate(0, -11) rotate(45)">
          <svg xmlns="http://www.w3.org/2000/svg" class="translate">
            <circle class="stroked" fill="white" stroke="blue" cx="8" cy="8" r="6"/>
            <path class="stroked" fill="none" stroke="blue" d="M 8 0 L 8 6.5"/>
            <path class="stroked" fill="none" stroke="blue" d="M 0 8 L 6.5 8"/>
            <path class="stroked" fill="none" stroke="blue" d="M 8 9.5 L 8 16"/>
            <path class="stroked" fill="none" stroke="blue" d="M 9.5 8 L 16 8"/>
          </svg>
        </g>
        </g>

      </svg>
    </div>

    <div class="mx-4 flex flex-column flex-grow-1">
      <Panel header="Area Settings">
        <div class="flex-row m-2">
          <label>Area height</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" :min="settings.tableHeight+settings.areaPaddingY+settings.seatHeight+settings.seatDist"
                         v-model="settings.areaHeight"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Area width</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" :min="settings.tableWidth+settings.areaPaddingX" v-model="settings.areaWidth"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Area middle relative to area</label>
          <div class="p-inputgroup">
            <InputNumber prefix="x: " step="0.1" suffix=" %" :min=0 :max="100" v-model="settings.areaOffsetX"/>
            <InputNumber prefix="y: " step="0.1" suffix=" %" :min=0 :max="100" v-model="settings.areaOffsetY"/>
          </div>
        </div>
      </Panel>
      <Panel header="Table Settings" class="my-3">
        <div class="flex-row m-2">
          <label>Table width</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" v-model="settings.tableWidth"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Table height</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" v-model="settings.tableHeight"/>
          </div>
        </div>
      </Panel>
      <Panel header="Seat Settings" class="my-3">
        <div class="flex-row m-2">
          <label>Number of seats</label>
          <div class="p-inputgroup">
            <InputNumber :min=0 :suffix="' seat'+(settings.seatNum===1 ? '' : 's')" v-model="settings.seatNum"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Seat height</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" v-model="settings.seatHeight"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Distance between seats</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" :min=0 v-model="settings.seatSep"/>
          </div>
        </div>
        <div class="flex-row m-2">
          <label>Distance to table</label>
          <div class="p-inputgroup">
            <InputNumber :suffix="' '+settings.distanceUnit" v-model="settings.seatDist"/>
          </div>
        </div>
      </Panel>
    </div>

  </div>
</template>

<style scoped>
.grid {
  display: inline-grid;
  width: 100%;
  grid-template-columns: max-content auto min-content;
}

svg rect {
  color: lightgray;
}

.slider-blue {
  --slider-connect-bg: none;
}

</style>

<script setup>

import {teamareaStore} from "@/stores/teamarea";
import TeamTable from "../components/Layout/TeamTable.vue";
import {computed, onMounted, ref} from "vue";
import {useKeyModifier} from '@vueuse/core';

const settings = teamareaStore()
const control = useKeyModifier('Control')

const middleShown = computed(() => showMid.value || activeTab.value === 1 || true);
const showMid = ref(false);
const activeTab = ref(false);
const scale = ref(15)

const sliderStep = computed(() => control.value ? 1 : -1)

const tableY = computed({
  get() {
    return [
      settings.areaPaddingY,
      settings.areaPaddingY + settings.tableHeight,
      settings.areaPaddingY + settings.tableHeight + settings.seatDist,
      settings.areaPaddingY + settings.tableHeight + settings.seatDist + settings.seatHeight,
    ]
  },
  set(n) {
    settings.areaPaddingY = n[0]
    settings.tableHeight = n[1] - n[0]
    settings.seatDist = n[2] - n[1]
    settings.seatHeight = n[3] - n[2]
  }
})

const tableX = computed({
  get() {
    return [
      settings.areaPaddingX,
      settings.areaPaddingX + settings.tableWidth,
    ]
  },
  set(n) {
    settings.areaPaddingX = n[0]
    settings.tableWidth = n[1] - n[0]
  }
})

const svgDiv = ref()

onMounted(() => {
  const ro = new ResizeObserver((e) => {

    // Consider changing this to retrieving the size
    e = e[0].contentRect
    const newScale = svgDiv.value.offsetWidth / (settings.areaWidth + 4)
    if (Math.abs(scale.value - newScale) < 0.001) {
      return
    }
    scale.value = newScale
    console.log("Setting scale to: "+newScale)
  })
  ro.observe(svgDiv.value)
})



</script>
