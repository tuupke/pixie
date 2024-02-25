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
      <Accordion>
        <AccordionTab>
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">Info</span>
                <Button class="ml-auto" size="small" icon="pi pi-plus" severity="success" rounded
                        outlined v-on:click.stop=""/>
            </span>
          </template>
          <div>
            Number of rooms: {{ rooms.length }}
          </div>
          <div>
            Number of tables/areas: {{
              rooms.reduce((totalTables, room): number => {
                return totalTables + room.elements.reduce((previous, element): number => {
                  return previous + element.repeats.reduce((p, r) => {
                    return p * r.num
                  }, 1)
                }, 0)
              }, 0)
            }}

          </div>
        </AccordionTab>

        <AccordionTab v-for="room in rooms">
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                          <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">Room: {{ room.name }}</span>
                <Button class="ml-auto" size="small" icon="pi pi-pencil" severity="warning" rounded
                        outlined v-on:click.stop=""/>
            </span>
            </span>
          </template>
          <div>
            Number of placed elements: {{ room.elements.length }}
          </div>
          <div>
            Number of tables/areas: {{
              room.elements.reduce((previous, element): number => {
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
           @mousedown="(e: MouseEvent) => dragStart({coord: topLeft, event: e})">
        <g ref="innerSvgRef" :transform="'scale('+scale+') translate('+topLeft.x+','+topLeft.y+')'">
          <DragNg v-for="room in rooms" @dragStart="(e: any) => dragStart({coord: room, event: e})">
            <g class="roomcontainer" :transform="'translate('+room.x+','+room.y+')'">
             <Room  v-bind="(<RoomInterface>room)"  />
          </g>
          </DragNg>
        </g>
<!--        <Hatching :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15" id="ratrst"/>-->
<!--        <Hatching colorA="orange" colorB="lightgreen" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"-->
<!--                  id="noTeamHatching"/>-->
<!--        <Hatching colorA="orange" colorB="white" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"-->
<!--                  id="noTeamHatching"/>-->
<!--        <Hatching colorA="lightgray" colorB="white" :dist="settings.areaWidth/20" :stroke="settings.areaWidth/15"-->
<!--                  id="selectedHatching"/>-->
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
import {computed, onMounted, reactive, ref} from "vue";
import {useKeyModifier} from '@vueuse/core';
import {
  Coordinate,
  DragStartEvent,
  RoomInterface,
  RotationCoordinate,
  SequenceAxis,
  SequenceDirection,
  SequenceInterface,
  SequenceType
} from "../types.ts";

// import Hatching from "../components/Layout/Hatching.vue";
import {teamareaStore} from "../stores/teamarea";
import DragNg from "./DragNg.vue";

const settings = teamareaStore()

function axisMap(type: SequenceType, axis: SequenceAxis): string {
  switch (type) {
    case SequenceType.Line:
      switch (axis) {
        case SequenceAxis.Horizontal:
          return "Horizontal"
        case SequenceAxis.Vertical:
          return "Vertical"
      }
    case SequenceType.Circle:
      switch (axis) {
        case SequenceAxis.Horizontal:
          return "Clockwise"
        case SequenceAxis.Vertical:
          return "Counterclockwise"
      }
  }

  return 'unknown'
}

function directionMap(type: SequenceType, axis: SequenceAxis, direction: SequenceDirection): string {
  switch (type) {
    case SequenceType.Line:
      switch (axis) {
        case SequenceAxis.Horizontal:
          switch (direction) {
            case SequenceDirection.Positive:
              return "left"
            case SequenceDirection.Negative:
              return "right"
          }
          break
        case SequenceAxis.Vertical:
          switch (direction) {
            case SequenceDirection.Positive:
              return "behind"
            case SequenceDirection.Negative:
              return "in front"
          }
          break
      }
      break
    case SequenceType.Circle:
      switch (direction) {
        case SequenceDirection.Positive:
          return "backs"
        case SequenceDirection.Negative:
          return "fronts"
      }
  }

  return 'unknown'
}

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

let data: RoomPlacement[] = [
  {
    x: 0,
    y: 0,
    rotation: 0,
    name: "Main room",
    outline: [
      {x: 0, y: 0},
      {x: 10, y: 0},
      {x: 10, y: 10},
      {x: 0, y: 10}
    ],
    elements: [
      {
        base: {x: 100, y: 100, rotation: 0},
        repeats: [
          {
            type: SequenceType.Line,
            num: 3,
            radius: 0,
            axis: SequenceAxis.Horizontal,
            dir: SequenceDirection.Negative,
            separation: 1000,
            equivalentSpaced: true,
          },
          {
            type: SequenceType.Circle,
            num: 2,
            axis: SequenceAxis.Horizontal,
            radius: 500,
            separation: 0,
            equivalentSpaced: true,
            dir: SequenceDirection.Positive,
          },
        ]
      },
      {
        base: {x: 100, y: 400, rotation: 0},
        repeats: []
      }
    ]
  }
];

// Ensure all repeats have correct props
const fallback: SequenceInterface = {
  type: SequenceType.Line,
  num: 1,
  axis: SequenceAxis.Horizontal,
  dir: SequenceDirection.Positive,
  radius: 100,
  separation: 100,
  equivalentSpaced: true,
}

for (let ri = 0; ri < data.length; ri++) {
  for (let li = 0; li < data[ri].elements.length; li++) {
    data[ri].elements[li].repeats = data[ri].elements[li].repeats.map(e => {
      return {...fallback, ...e}
    })
  }
}

type RoomPlacement = RoomInterface & RotationCoordinate
const rooms = reactive<RoomPlacement[]>(data)
const confirmdelete = ref<number[]>([]);
const snapToGrid = ref(false);

const shift = useKeyModifier('Shift')
const control = useKeyModifier('Control')
const alt = useKeyModifier('Alt')

const clamping = computed(() => ((shift.value || snapToGrid.value) ? 5 : 0) * ((control.value || snapToGrid.value) ? 3 : 1))

// function maybeTranslateRotate(e: MouseEvent) {
//   if (translateRotate.translatingMap) {
//     // coord is in svg space
//     const coord = [
//       (e.clientX + translateRotate.offset[0]) / scale,
//       (e.clientY + translateRotate.offset[1]) / scale,
//     ];
//
//     if (clamping.value > 1) {
//       const cv = clamping.value
//       coord[0] = Math.round(coord[0] / cv) * cv
//       coord[1] = Math.round(coord[1] / cv) * cv
//     }
//
//     topLeft.x = coord[0]
//     topLeft.y = coord[1]
//
//     console.log(coord)
//   } else if (translateRotate.translatingRoom) {
//     const room = translateRotate.translatingRoom[0]
//     const el = translateRotate.translatingRoom[1]
//
//     // coord is in svg space
//     const coord = [
//       (e.clientX + translateRotate.offset[0]) / scale,
//       (e.clientY + translateRotate.offset[1]) / scale,
//     ];
//
//     if (clamping.value > 1) {
//       const cv = clamping.value
//       coord[0] = Math.round(coord[0] / cv) * cv
//       coord[1] = Math.round(coord[1] / cv) * cv
//     }
//
//     rooms[room].elements[el].base.x = coord[0]
//     rooms[room].elements[el].base.y = coord[1]
//   } else if (translateRotate.rotatingRoom) {
//     const room: number = translateRotate.rotatingRoom[0]
//     const el: number = translateRotate.rotatingRoom[1]
//
//     // coord is in svg space
//     const coord = [
//       e.clientX + translateRotate.offset[0],
//       e.clientY + translateRotate.offset[1],
//     ];
//
//     const xDiff = coord[0] - rooms[room].elements[el].base.x
//     const yDiff = coord[1] - rooms[room].elements[el].base.y
//
//     const upVect = Math.atan2(-1, 0)
//     let angle = (Math.atan2(yDiff, xDiff) - upVect) * 180 / Math.PI
//
//     if (clamping.value > 1) {
//       angle = Math.round(angle / clamping.value) * clamping.value
//     }
//
//     rooms[room].elements[el].base[2] = angle
//   }
// }

function resetTranslateRotate() {
  coordinateBeingTranslated.value = null
  coordinateBeingRotated.value = null
}

function maybeTranslateRotate(e: MouseEvent) {
  if (coordinateBeingTranslated.value === undefined || coordinateBeingTranslated.value === null) {
    return
  }

  const innerCoords = toInnerCoordinates(e)
  let newCoord = {
    x: innerCoords.x - offset.value.x,
    y: innerCoords.y - offset.value.y,
  }

  if (clamping.value > 1) {
    const cv = clamping.value
    newCoord.x = Math.round(newCoord.x / cv) * cv
    newCoord.y = Math.round(newCoord.y / cv) * cv
  }

  coordinateBeingTranslated.value.x = newCoord.x
  coordinateBeingTranslated.value.y = newCoord.y
}

function dragStart(e: DragStartEvent) {
  if ("rotation" in (e.coord) && alt.value) {
    coordinateBeingRotated.value = {
      active: e.coord,
      original: {
        x: e.coord.x.valueOf(),
        y: e.coord.y.valueOf(),
        rotation: e.coord.rotation.valueOf()
      }
    }

    return
  }

  coordinateBeingTranslated.value = e.coord
  const innerCoords = toInnerCoordinates(e.event)
  offset.value = {
    x: innerCoords.x - e.coord.x,
    y: innerCoords.y - e.coord.y,
  }
}

const topLeft = reactive<Coordinate>({x: 0, y: 0});
const offset = ref<Coordinate>({x: 0, y: 0});
const coordinateBeingTranslated = ref<Coordinate | null>();
const coordinateBeingRotated = ref<{ active: RotationCoordinate, original: RotationCoordinate } | null>();
const scale = ref<number>(1)

const svgRef = ref()
const innerSvgRef = ref()

function toInnerCoordinates(e: MouseEvent): Coordinate {
  const br = svgRef.value.getBoundingClientRect();

  return {
    x: (e.clientX - br.left) / scale.value,
    y: (e.clientY - br.top) / scale.value,
  }
}

// TODO extract into a 'reset panzoom' function
onMounted(() => {
  const padding = 20
  const bbox = svgRef.value.getBBox({
    stroke: true,
  })

  const brect = svgRef.value.getBoundingClientRect()
  const width = brect.width - 2 * padding
  const height = brect.height - 2 * padding

  const scaleX = width / (bbox.width)
  const scaleY = height / (bbox.height)

  topLeft.x = -bbox.x
  topLeft.y = -bbox.y

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
})


const selectedRoom = ref<Room | null>();
const selectedElement = ref<Element | null>();

function maybeShow(e: any) {
  selectedRoom.value = e.room
  selectedElement.value = e.el
}

</script>
