<template>
  <Toolbar class="flex-initial flex m-2">
    <template #start>
      <Button label="New Room" icon="pi pi-plus" class="mr-2" severity="success"/>
      <Button label="New Room element" icon="pi pi-plus" class="mr-2" severity="success"/>
      <ToggleButton on-label="Snapped" off-label="free-form" v-model="snapToGrid"/>
    </template>
  </Toolbar>
  <div class="settings-view-wrapper">
    <div class="col-3 left-nav">
      <Accordion v-if="selectedRoom == null">
        <AccordionTab>
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">Info</span>
                <Button class="ml-auto" size="small" icon="pi pi-plus" severity="success" rounded
                        outlined v-on:click.stop=""/>
            </span>
          </template>
          <div>
            Number of rooms: {{ map.placements.length }}
          </div>
          <div>
            Number of tables/areas: {{
              map.placements.reduce((totalTables, roomPlacement): number => {
                return totalTables + roomPlacement.room.elements.reduce((previous, element): number => {
                  return previous + element.repeats.reduce((p, r) => {
                    return p * r.num
                  }, 1)
                }, 0)
              }, 0)
            }}
          </div>
          <div v-if="coordinateBeingRotated !== null">
          </div>
        </AccordionTab>

        <AccordionTab v-for="placement in map.placements">
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                          <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ placement.room.name }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-pencil" severity="warning" rounded
                        outlined v-on:click.stop="selectedRoom = placement.room; rescale()"/>
            </span>
            </span>
          </template>
          <div>
            Number of placed elements: {{ placement.room.elements.length }}
          </div>
          <div>
            Number of tables/areas: {{
              placement.room.elements.reduce((previous, element): number => {
                return previous + element.repeats.reduce((p, r) => {
                  return p * r.num
                }, 1)
              }, 0)
            }}
          </div>

        </AccordionTab>
      </Accordion>
      <Accordion v-else>
        <AccordionTab>
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ selectedRoom.name }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-backward" severity="success" rounded
                        outlined v-on:click.stop="selectedRoom = null; ; rescale()"/>
            </span>
          </template>
          <div>
            Number of placed elements: {{ selectedRoom!.elements.length }}
          </div>
          <div>
            Number of tables/areas: {{
              selectedRoom!.elements.reduce((previous, element): number => {
                return previous + element.repeats.reduce((p, r) => {
                  return p * r.num
                }, 1)
              }, 0)
            }}
          </div>
        </AccordionTab>
        <!--        <AccordionTab-->
        <!--            v-for="(repeat, k) in selectedElement.repeats">-->
        <!--          <template #header>-->
        <!--            <span class="flex align-items-center gap-2 w-full">-->
        <!--                <span class="font-bold white-space-nowrap">{{ repeat.type }}</span>-->
        <!--                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded-->
        <!--                        aria-label="confirm deletion"-->
        <!--                        v-if="confirmdelete[Math.pow(2, selectedRoom[0]) + Math.pow(3, selectedRoom[1]) + Math.pow(5, k)]"-->
        <!--                        v-on:click.stop="deleteRepeats(k)"/>-->
        <!--                <Button class="ml-auto" size="small" icon="pi pi-times" severity="danger" rounded aria-label="Delete"-->
        <!--                        outlined v-else v-on:click.stop="toggleDelete(k)"/>-->
        <!--            </span>-->
        <!--          </template>-->
        <!--          <div class="flex flex-column">-->
        <!--            <div class="flex flex-row">-->
        <!--              <div class="flex flex-shrink-0 flex-column m-2">-->
        <!--                <label :for="'type'+k">Shape type</label>-->
        <!--                <SelectButton-->
        <!--                    :allowEmpty=false-->
        <!--                    v-model="repeat.type" :options="[SequenceType.Line, SequenceType.Circle]" :inputId="'type'+k"/>-->
        <!--              </div>-->
        <!--            </div>-->

        <!--            <div class="flex-row m-2" v-if="repeat.type === SequenceType.Circle">-->
        <!--              <label :for="'repeat'+k">Radius</label>-->
        <!--              <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons class="p-inputgroup"/>-->
        <!--            </div>-->

        <!--            <div class="flex flex-shrink-0 flex-column m-2">-->
        <!--              &lt;!&ndash;              <label v-if="repeat.type === SequenceType.Line" :for="k+'axis'">Orientation</label>&ndash;&gt;-->
        <!--              &lt;!&ndash;              <label v-else :for="k+'axis'"> clockwise </label>&ndash;&gt;-->
        <!--              <label :for="k+'axis'"> Direction </label>-->
        <!--              <SelectButton-->
        <!--                  :allowEmpty=false-->
        <!--                  :optionLabel="(data: SequenceAxis) => axisMap(repeat.type, data)"-->

        <!--                  v-model="repeat.axis" :options="[SequenceAxis.Horizontal, SequenceAxis.Vertical]" :inputId="k+'axis'"/>-->
        <!--            </div>-->


        <!--            <div class="flex flex-shrink-0 flex-column m-2">-->
        <!--              <label v-if="repeat.type === SequenceType.Line" :for="k+'dir'" class="ml-2">-->
        <!--                Direction-->
        <!--              </label>-->
        <!--              <label v-else :for="k+'dir'" class="ml-2"> backs facing</label>-->

        <!--              <SelectButton-->
        <!--                  :allowEmpty=false-->
        <!--                  v-model="repeat.dir"-->
        <!--                  :options="[SequenceDirection.Negative, SequenceDirection.Positive]"-->
        <!--                  :inputId="k+'dir'"-->
        <!--                  :optionLabel="(data: SequenceDirection) => directionMap(repeat.type, repeat.axis, data)"-->
        <!--              />-->
        <!--            </div>-->

        <!--            <div class="flex-row m-2">-->
        <!--              <label :for="'repeat'+k">Repeat for #times</label>-->
        <!--              <InputNumber-->
        <!--                  class="p-inputgroup"-->
        <!--                  id="'repeat'+k"-->
        <!--                  show-buttons-->
        <!--                  v-model="repeat.num"/>-->
        <!--            </div>-->

        <!--            <div class="flex-row m-2">-->
        <!--              <label :for="'separation'+k">Separation</label>-->

        <!--              <div class="p-inputgroup">-->
        <!--              <span v-if="repeat.type ===  SequenceType.Circle" class="p-inputgroup-addon">-->
        <!--                  <Checkbox v-model="repeat.equivalentSpaced" :binary="true"/>-->
        <!--              </span>-->
        <!--                <InputNumber-->
        <!--                    :inputId="'separation'+k"-->
        <!--                    :disabled="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle"-->
        <!--                    :modelValue="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle ? 360/Math.max(1, repeat.num) : repeat.separation"-->
        <!--                    @update:model-value="(newValue: number) => repeat.separation = newValue"-->
        <!--                    showButtons-->
        <!--                    :suffix="repeat.type ===  SequenceType.Circle ? '°':' cm'"/>-->
        <!--              </div>-->
        <!--            </div>-->
        <!--          </div>-->
        <!--        </AccordionTab>-->
      </Accordion>
    </div>
    <div class="col-9">
      <svg style="border: 1px solid red" id="layoutsvg" ref="svgRef" @wheel="scroll" width="100%" height="100%"
           @mousemove="maybeTranslateRotate" @mouseup="resetTranslateRotate"
           @mousedown="(e: MouseEvent) => dragStart({coord: topLeft, event: e}, false)">
        <g ref="innerSvgRef" :transform="'scale('+scale+') translate('+topLeft.x+','+topLeft.y+')'">
          <DragNg
              v-if="selectedRoom==null"
              v-for="placement in map.placements"
              :coord=placement.coord
              :transform=true
              @dragStart="dragStart"
              @rotateStart="rotateStart">
            <Room v-bind="placement.room"         :x=0
                  :y=0
                  :rotation=0
            />
          </DragNg>

          <Room v-else v-bind="selectedRoom" :translating=true
                @dragStart="dragStart"
                @rotateStart="rotateStart"
          />
        </g>
        <Hatching colorA="orange" colorB="lightgreen" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"
                  id="noTeamHatching"/>
        <Hatching colorA="orange" colorB="white" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"
                  id="noHostHatching"/>
        <Hatching colorA="lightgray" colorB="white" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"
                  id="selectedHatching"/>
      </svg>
    </div>
  </div>
</template>

<style scoped>
svg rect {
  color: lightgray;
}

.settings-view-wrapper {
  display: flex;
  flex: 1;
}

.left-nav {
  max-height: calc(100vh - 9rem);
  overflow-y: auto;
}
</style>

<script setup lang="ts">

import Room from "../components/Layout/Room.vue";
import Hatching from "../components/Layout/Hatching.vue";
import Crosshairs from "../components/Layout/Crosshairs.vue";
import {computed, onMounted, provide, reactive, ref} from "vue";
import {useKeyModifier} from '@vueuse/core';
import {
  CoordinateInterface,
  DragStartEvent,
  RoomInterface,
  RotationCoordinateInterface,
  RotationStartEvent,
  Vector
} from "../types.ts";

// import Hatching from "../components/Layout/Hatching.vue";
import {teamareaStore} from "../stores/teamarea";
import {mapStore} from "../stores/map";
import DragNg from "./DragNg.vue";
import Sequence from "../components/Layout/Sequence.vue";

provide('toInnerCoordinates', toInnerCoordinates)

const settings = teamareaStore()
const map = mapStore()

function scroll(state: WheelEvent) {
  // Before calculating scale, see where the cursor currently is. Needs to be kept 'constant'.
  const innerCoords = toInnerCoordinates(state)

  if ((scale.value >= 3 && state.deltaY < 0) || (scale.value <= 0.03 && state.deltaY > 0)) {
    return;
  }

  const dt = 0.06
  const scrollDt = 1 + (state.deltaY < 0 ? dt : -dt)

  scale.value *= scrollDt

  // Calculate the new top-left offset
  const newInner = toInnerCoordinates(state)
  topLeft.x += newInner.x - innerCoords.x
  topLeft.y += newInner.y - innerCoords.y
}

// function toggleDelete(k: number) {
//   const key = Math.pow(2, selectedRoom[0]) + Math.pow(3, selectedRoom[1]) + Math.pow(5, k)
//   confirmdelete.value[key] = true
//   window.setTimeout(() => delete confirmdelete.value[key], 1000)
// }
//
// function addRepeats() {
//   rooms[selectedRoom[0]].elements[selectedRoom[1]].repeats.push({...fallback});
// }
//
// function deleteRepeats(k: number) {
//   const key = Math.pow(2, selectedRoom[0]) +
//       Math.pow(3, selectedRoom[1]) +
//       Math.pow(5, k)
//   if (!confirmdelete.value[key]) {
//     return
//   }
//
//   delete confirmdelete.value[k]
//   rooms[selectedRoom[0]].elements[selectedRoom[1]].repeats.splice(k, 1);
// }

const confirmdelete = ref<number[]>([]);
const snapToGrid = ref(false);

const shift = useKeyModifier('Shift')
const control = useKeyModifier('Control')

const clamping = computed(() => ((shift.value || snapToGrid.value) ? 5 : 0) * ((control.value || snapToGrid.value) ? 3 : 1))

function resetTranslateRotate() {
  coordinateBeingTranslated.value = null
  coordinateBeingRotated.value = null
  offset.value = {x: 0, y:0}
}

function maybeTranslateRotate(e: MouseEvent) {
  if (coordinateBeingRotated.value !== undefined && coordinateBeingRotated.value !== null) {
    const mouse = toInnerCoordinates(e, true)
    const mid = coordinateBeingRotated.value?.middle
    let angle = new Vector(mouse.x-mid.x, mouse.y - mid.y).asAngle() - coordinateBeingRotated.value?.angle

    let coordAngle = angle + coordinateBeingRotated.value?.initialAngle
    if (clamping.value > 1) {
      // angle + coordinateBeingRotated.value?.initialAngle must be multiple of clamping.value
      coordAngle = Math.round(coordAngle / clamping.value) * clamping.value
      angle = coordAngle - coordinateBeingRotated.value?.initialAngle
    }

    coordinateBeingRotated.value.coord.rotation = coordAngle

    let newCoord = coordinateBeingRotated.value?.toRoot.copy().rotate(angle).add(coordinateBeingRotated.value?.middle)
    coordinateBeingRotated.value.coord.x = newCoord.x
    coordinateBeingRotated.value.coord.y = newCoord.y
  } else if (coordinateBeingTranslated.value !== undefined && coordinateBeingTranslated.value !== null) {
    const innerCoords = toInnerCoordinates(e)
    let newCoord = {
      x: innerCoords.x - offset.value.x,
      y: innerCoords.y - offset.value.y,
    }

    // console.log(offset.value.x, offset.value.y, newCoord)

    if (clamping.value > 1) {
      const cv = clamping.value
      newCoord.x = Math.round(newCoord.x / cv) * cv
      newCoord.y = Math.round(newCoord.y / cv) * cv
    }

    coordinateBeingTranslated.value.x = newCoord.x
    coordinateBeingTranslated.value.y = newCoord.y
  }
}

function dragStart(e: DragStartEvent) {
  coordinateBeingTranslated.value = e.coord
  const innerCoords = toInnerCoordinates(e.event)

  offset.value = {
    x: innerCoords.x - e.coord.x,
    y: innerCoords.y - e.coord.y,
  }
}

interface RotationAroundInterface {
  coord: RotationCoordinateInterface
  middle: Vector
  toRoot: Vector
  angle: number
  initialAngle: number
  root: Vector
}

function rotateStart(e: RotationStartEvent) {
  const rotCoord = e.coord as RotationCoordinateInterface

  // Maybe use these instead of the box size to calculate where someone clicked
  // const handleCoords = toInnerCoordinates(e.event)

  // Derivation starts at the coordinate itself
  //  1. derives vector from coordinate to middle
  //  2. derives angle from middle to handle
  //  3. derives inverse of vector from coordinate to middle
  //     i.e. from middle to coordinate by flipping signs
  const toMiddle = new Vector(e.width / 2 + e.x, e.height / 2 + e.y).rotate(rotCoord.rotation)

  const middle = toMiddle.copy().add(e.coord)
  const toRoot = toMiddle.copy().multiply(-1)

  // calculate angle to mouse. Should be similar!
  const mouse = toInnerCoordinates(e.event, true)
  // mouse.x = mouse.x - topLeft.x;
  // mouse.y = mouse.y - topLeft.y;
  const toMouse = new Vector(mouse.x - middle.x, mouse.y - middle.y)
  const angle = toMouse.asAngle(); // - rotCoord.rotation

  coordinateBeingRotated.value = {
    coord: rotCoord,
    middle: middle,
    angle: angle,
    toRoot: toRoot,
    initialAngle: rotCoord.rotation,
    root: middle.copy().add(toRoot),
  }
}

const topLeft = reactive<CoordinateInterface>({x: 0, y: 0});
const offset = ref<CoordinateInterface>({x: 0, y: 0});
const coordinateBeingTranslated = ref<CoordinateInterface | null>();
const coordinateBeingRotated = ref<RotationAroundInterface | null>();
const scale = ref<number>(1)

const svgRef = ref()
const innerSvgRef = ref()

function toInnerCoordinates(e: MouseEvent, translate: boolean = false): CoordinateInterface {
  const br = svgRef.value.getBoundingClientRect();
  return {
    x: (e.clientX - br.left) / scale.value - (translate ? topLeft.x: 0),
    y: (e.clientY - br.top) / scale.value -  (translate ? topLeft.y: 0),
  }
}

function toCoord(relativeTo: CoordinateInterface): (e: MouseEvent) => CoordinateInterface {
  return (e: MouseEvent) => {
    const br = svgRef.value.getBoundingClientRect();
    return {
      x: (e.clientX - br.left) / scale.value - relativeTo.x,
      y: (e.clientY - br.top) / scale.value - relativeTo.y,
    }
  }
}

// TODO extract into a 'reset panzoom' function
onMounted(rescale)

function rescale() {
  topLeft.x = 0
  topLeft.y = 0
  scale.value = 0

  const padding = 40
  const bbox = svgRef.value.getBBox({
    stroke: true,
  })

  // brect is the size of the entire svg canvas in screen coordinates.
  const brect = svgRef.value.getBoundingClientRect()
  const width = brect.width - 2 * padding
  const height = brect.height - 2 * padding

  const scaleX = width / (bbox.width)
  const scaleY = height / (bbox.height)

  console.log('rescale', bbox, brect)

  topLeft.x -= bbox.x
  topLeft.y -= bbox.y

  // Align the map to the center
  if (scaleX <= scaleY) {
    scale.value = scaleX;
    topLeft.y += (brect.height / scaleX - bbox.height) / 2
    topLeft.x += padding / scale.value
  } else {
    scale.value = scaleY;
    topLeft.x += (brect.width / scaleY - bbox.width) / 2
    topLeft.y += padding / scale.value
  }
}

const selectedRoom = ref<RoomInterface | null>();
const selectedElement = ref<Element | null>();

function maybeShow(e: any) {
  selectedRoom.value = e.room
  selectedElement.value = e.el
}

</script>
