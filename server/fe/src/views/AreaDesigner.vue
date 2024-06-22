<template>
  <div class="grid">
    <div class="col-3">
      <Accordion>
        <AccordionPanel value="0">
          <AccordionHeader>Area settings</AccordionHeader>
          <AccordionContent>
            <div class="flex-row m-2">
              <label>Area height</label>
              <div class="p-inputgroup">
                <InputNumber :suffix="' '+settings.distanceUnit"
                             :min="settings.tableHeight+settings.areaPaddingY+settings.seatHeight+settings.seatDist"
                             v-model="settings.areaHeight"/>
              </div>
            </div>
            <div class="flex-row m-2">
              <label>Area width</label>
              <div class="p-inputgroup">
                <InputNumber :suffix="' '+settings.distanceUnit" :min="settings.tableWidth+settings.areaPaddingX"
                             v-model="settings.areaWidth"/>
              </div>
            </div>
            <div class="flex-row m-2">
              <label>Area middle relative to area</label>
              <div class="p-inputgroup">
                <InputNumber prefix="x: " :step=0.1 suffix=" %" :min=0 :max="100" v-model="settings.areaOffsetX"/>
                <InputNumber prefix="y: " :step=0.1 suffix=" %" :min=0 :max="100" v-model="settings.areaOffsetY"/>
              </div>
            </div>
            <div class="flex-row m-2">
              <label>Border stroke width</label>
              <div class="p-inputgroup">
                <InputNumber suffix=" px" :min="1"
                             v-model="settings.strokeWidth"/>
              </div>
            </div>
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="1">
          <AccordionHeader>Table settings</AccordionHeader>
          <AccordionContent>
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
          </AccordionContent>
        </AccordionPanel>
        <AccordionPanel value="2">
          <AccordionHeader>Seat settings</AccordionHeader>
          <AccordionContent>
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
          </AccordionContent>
        </AccordionPanel>
      </Accordion>
    </div>
    <div class="col-9">
      <div class="gridLocal">
        <div></div>
        <div class="my-3 flex flex-column align-content-end">
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
          <svg width="100%" height="100%">
            <g :transform="'scale('+scale+') translate('+(-settings.areaX)+','+(-settings.areaY)+')'">
              <TeamTable :x="0.5" :rotation="0" :y="0.5" team-id="100" />
              <!--              <circle :cx=0.5 :cy=0.5 r="3" fill="orange"/>-->
              <CrossHairs/>
            </g>
          </svg>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.gridLocal {
  display: inline-grid;
  grid-template-columns: max-content max-content;
}

.slider-blue {
  --slider-connect-bg: none;
}

</style>

<script setup lang="ts">

import {teamareaStore} from "../stores/teamarea";
import {computed, onMounted, provide, ref, watch} from "vue";
import {useKeyModifier} from '@vueuse/core';
import {storeToRefs} from "pinia";
import CrossHairs from "../components/Layout/CrossHairs.vue";
import '@vueform/slider/themes/default.css'
import Slider from '@vueform/slider'

const settings = teamareaStore()
const control = useKeyModifier('Control')

provide('editing', false)

const scale = ref(1)

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
  (new ResizeObserver(rescale)).observe(svgDiv.value.parentNode.parentNode)
})

const {areaWidth, strokeWidth} = storeToRefs(settings)

watch(areaWidth, rescale)
watch(strokeWidth, rescale)

async function rescale() {
  const svg = svgDiv.value
  const availableWidth = svg.parentNode.parentNode.offsetWidth -
      svg.previousSibling.offsetWidth

  const newScale = availableWidth / (settings.areaWidth + 20)
  if (Math.abs(scale.value - newScale) < 0.001) {
    return
  }

  scale.value = newScale
  console.log("Setting scale to: " + newScale)
}

</script>
