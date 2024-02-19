<template>
  <div class="flex flow-column">
    <span class="flex-1">
      &nbsp;
    </span>
    <span class="flex-1">
    vertical settings

    </span>
  </div>

  <div class="h-full flex flex-1 flex-row my-3" style="min-height: 600px;">
    <span class="mx-3" >
      <Slider v-model="areaHeight" :min="-600/scale" :max="0" orientation="vertical" class='h-full' />
    </span>
    <span class="mx-3" >
      <Slider v-model="areaOffsetY" :min=-100 :max=0 :style="'min-height: ' +scale*settings.areaHeight+ 'px;'" orientation="vertical" />
    </span>
    <span class="mx-2">
      <Slider v-model=table
              orientation="vertical" />
<!--    <Slider v-model="areaOffsetY" range orientation="vertical" />-->
<!--    <Slider v-model="areaOffsetY" range orientation="vertical" />-->
    </span>
    <div class="flex flex-1 mx-2">
      <svg id="layoutsvg">
        <g :transform="'scale('+scale+') translate('+(-settings.areaX)+','+(-settings.areaY)+')'">
          <TeamTable :x="1" :rotation="0" :y="1" team-id="100"/>
          <circle v-if="middleShown" :cx=1 :cy=1 r="3" fill="orange"/>

          <!--          <g v-if="middleShown" transform="translate(0, -11) rotate(45)">-->
<!--            <svg xmlns="http://www.w3.org/2000/svg" @mousedown="dragStart" @click="select" class="translate">-->
<!--              <circle class="stroked" fill="white" stroke="blue" cx="8" cy="8" r="6"/>-->
<!--              <path class="stroked" fill="none" stroke="blue" d="M 8 0 L 8 6.5"/>-->
<!--              <path class="stroked" fill="none" stroke="blue" d="M 0 8 L 6.5 8"/>-->
<!--              <path class="stroked" fill="none" stroke="blue" d="M 8 9.5 L 8 16"/>-->
<!--              <path class="stroked" fill="none" stroke="blue" d="M 9.5 8 L 16 8"/>-->
<!--            </svg>-->
<!--          </g>-->
        </g>
      </svg>
    </div>
<!--    <div class="flex-1 p-2 m-2">-->
<!--      <Accordion v-model:activeIndex="activeTab">-->
<!--        <AccordionTab header="Area sizing" class="flex flex-0">-->
<!--          <div class="p-2">-->
<!--            Area Width {{ settings.areaWidth }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.areaWidth"/>-->
<!--          </div>-->
<!--          <div class="p-2">-->
<!--            Area Height {{ settings.areaHeight }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.areaHeight"/>-->
<!--          </div>-->

<!--          <div class="p-2">-->
<!--            Area paddingX {{ settings.areaPaddingX }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.areaPaddingX"/>-->
<!--          </div>-->
<!--          <div class="p-2">-->
<!--            Area paddingY {{ settings.areaPaddingY }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.areaPaddingY"/>-->
<!--          </div>-->
<!--        </AccordionTab>-->
<!--        <AccordionTab>-->

<!--          <template #header>-->
<!--            <span class="flex align-items-center gap-2 w-full">-->
<!--                <span class="font-bold white-space-nowrap">Area positioning</span>-->
<!--                <ToggleButton class="ml-auto" size="small" v-model="showMid" onIcon="pi pi-eye" onLabel="middle visible"-->
<!--                              offIcon="pi pi-eye-slash" offLabel="middle hidden" @click.stop/>-->
<!--            </span>-->
<!--          </template>-->

<!--          <div class="p-2">-->
<!--            Area offset-x compared to center {{ settings.areaOffsetX }}%-->
<!--            <Slider v-model="settings.areaOffsetX" :min=0 :max=100-->
<!--            />-->
<!--          </div>-->
<!--          <div class="p-2">-->
<!--            Area offset-y compared to center {{ settings.areaOffsetY }}%-->
<!--            <Slider v-model="settings.areaOffsetY" :min=0 :max=100-->
<!--            />-->
<!--          </div>-->

<!--        </AccordionTab>-->
<!--        <AccordionTab header="Table sizing">-->
<!--          <div class="p-2">-->
<!--            table offset x {{ settings.tableOffsetX }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.tableOffsetX" :min="-settings.areaPaddingX" :max="settings.areaPaddingX"/>-->
<!--          </div>-->
<!--          <div class="p-2">-->
<!--            table offset y {{ settings.tableOffsetY }} {{ settings.distanceUnit }}-->
<!--            <Slider v-model="settings.tableOffsetY" :min="-settings.areaPaddingY" :max="settings.areaPaddingY"/>-->
<!--          </div>-->
<!--        </AccordionTab>-->
<!--        <AccordionTab></AccordionTab>-->
<!--      </Accordion>-->


<!--      <div class="p-2">-->
<!--        Separation between the seats-->
<!--        <Slider v-model="settings.seatSep"/>-->
<!--      </div>-->
<!--      <div class="p-2">-->

<!--      </div>-->
<!--      <div class="p-2">-->
<!--        Seat Distance to the table-->
<!--        <Slider v-model="settings.seatDist" :min=0 :max="settings.areaHeight - settings.seatHeight"/>-->
<!--      </div>-->

<!--      <div class="p-2">-->
<!--        Num Seats-->
<!--        <Slider v-model="settings.seatNum"/>-->
<!--      </div>-->
<!--      <div class="p-2">-->
<!--        Seat height-->
<!--        <Slider v-model="settings.seatHeight"/>-->
<!--      </div>-->
<!--      <div class="p-2">-->
<!--        Seat padding-->
<!--        <Slider v-model="settings.seatPadding"/>-->
<!--      </div>-->
<!--    </div>-->
  </div>
</template>

<style scoped>
svg rect {
  color: lightgray;
}

</style>

<script setup>

import {teamareaStore} from "@/stores/teamarea";
import TeamTable from "../components/Layout/TeamTable.vue";
import {computed, ref, toRaw, markRaw} from "vue";
import {Slider as fSlider} from '@vueform/slider'

const settings = teamareaStore()

const middleShown = computed(() => showMid.value || activeTab.value === 1 || true);
const showMid = ref(false);
const activeTab = ref(false);
const scale = ref(4)

const areaHeight = computed({get() {return -settings.areaHeight }, set(n) {settings.areaHeight = -n}})
const areaOffsetY = computed({get() {return -settings.areaOffsetY }, set(n) {settings.areaOffsetY = -n}})

const table = computed({
  get() {
    return [settings.areaPaddingY, 80]
  },
  set(n) {
    console.log(n)
    settings.areaPaddingY = n[0]
  }
})

</script>
