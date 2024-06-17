<template>
  <Toolbar class="flex-initial flex m-2">
    <template #start>
      <ToggleButton on-label="Snapped" off-label="free-form" class="mr-2" v-model="snapToGrid"/>
      <!--      <SelectButton v-model="pathMode" :options="EditTypes" optionLabel="name" optionValue="value" class="mr-2"/>-->
      <Button label="Reset panzoom" class="mr-2" @click="rescale"/>
      <Button v-if="selectedRoom==null" label="New Room" icon="pi pi-plus" class="mr-2" severity="success"/>
      <Button v-else label="New Element" icon="pi pi-plus" class="mr-2" severity="success"/>
    </template>
  </Toolbar>
  <div class="settings-view-wrapper">
    <div class="col-3 left-nav">
      <span v-if="selectedRoomIndex === null">
        <Card>
          <template #title>Info</template>
          <template #content>
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
          </template>
        </Card>
      <Accordion class="mt-2" :active-index="0">
<!--        <AccordionTab>-->
        <!--          <template #header>-->
        <!--            <span class="flex align-items-center gap-2 w-full">-->
        <!--                <span class="font-bold white-space-nowrap">Info</span>-->
        <!--            </span>-->
        <!--          </template>-->
        <!--          <div>-->
        <!--            Number of rooms: {{ map.placements.length }}-->
        <!--          </div>-->
        <!--          <div>-->
        <!--            Number of tables/areas: {{-->
        <!--              map.placements.reduce((totalTables, roomPlacement): number => {-->
        <!--                return totalTables + roomPlacement.room.elements.reduce((previous, element): number => {-->
        <!--                  return previous + element.repeats.reduce((p, r) => {-->
        <!--                    return p * r.num-->
        <!--                  }, 1)-->
        <!--                }, 0)-->
        <!--              }, 0)-->
        <!--            }}-->
        <!--          </div>-->
        <!--          <div v-if="coordinateBeingRotated !== null">-->
        <!--          </div>-->
        <!--        </AccordionTab>-->

        <AccordionTab v-for="(placement, index) in map.placements">
          <template #header>
            <span class="flex align-items-center gap-2 w-full">
                          <span class="flex align-items-center gap-2 w-full">
                <span class="font-bold white-space-nowrap">{{ placement.room.name }}</span>
            </span>
            </span>
          </template>

          <div class="flex flex-column gap-2">
            <label for="placementname">Name</label>
            <InputText v-model="placement.room.name" aria-describedby="placementname-help"/>
            <small id="placementname-help">Enter the name of this room.</small>
          </div>
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
          <div class="mt-2">
            <Button icon="pi pi-pencil" class="mr-2" severity="warning" label="edit"
                    v-on:click.stop="selectedRoomIndex = index; rescale()"/>
            <DeleteButton
                type="button"
                icon="pi pi-trash"
                severity="danger"
                label="Delete"
                :sequenceKey="[KeyCategory.Placements, index]"
            />
          </div>
        </AccordionTab>
      </Accordion>
      </span>
      <span v-else>
        <Card class="gap-2">
          <template #title>{{ selectedRoom!.name }} - info</template>
          <template #content>
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
            <Button
                class="mt-2"
                icon="pi pi-backward"
                label="back to overview"
                severity="warning"
                outlined
                v-on:click.stop="selectedRoomIndex = null; rescale()"
            />
          </template>
        </Card>
        <Card class="mt-2">
            <template #title>
              <span class="flex align-items-center gap-2 w-full">
                  <span class="font-bold white-space-nowrap">ELEMENT</span>
                  <Button
                      class="ml-auto"
                      size="small"
                      icon="pi pi-plus"
                      severity="success"
                      outlined
                      v-on:click.stop="addRepeats"/>
                  <DeleteButton
                      class="ml-1"
                      :sequence-key="highlightedKey"
                      :upto="KeyCategory.Elements"
                      size="small"
                      icon="pi pi-trash"
                      severity="danger"
                  />
              </span>
            </template>
          <template #content>
            <AssignmentOverview :sequenceKey="highlightedKey"></AssignmentOverview>
          </template>
          </Card>
        <Accordion v-if="highlightedKey" class="mt-2">
          <AccordionTab
              v-if="selectedElement!=null"
              v-for="(repeat, k) in selectedElement!.repeats">
            <template #header>
              <span class="flex align-items-center gap-2 w-full">
                  <span class="font-bold white-space-nowrap">R{{ k }}: {{ sequenceToString(repeat) }}</span>
                  <DeleteButton
                      class="ml-auto"
                      size="small"

                      icon="pi pi-trash"
                      severity="danger"
                      :sequenceKey="[KeyCategory.Placements, selectedRoomIndex, KeyCategory.Room, KeyCategory.Elements, selectedElementIndex!, KeyCategory.Repeats, k]"
                  />
              </span>
            </template>
            <div class="flex flex-column">
              <div class="flex flex-row">
                <div class="flex flex-shrink-0 flex-column m-2">
                  <label :for="'type'+k">Shape type</label>
                  <SelectButton :allowEmpty=false
                                v-model="repeat.type" :options="[SequenceType.Line, SequenceType.Circle]"
                                :inputId="'type'+k"/>
                </div>
              </div>

              <div class="flex-row m-2" v-if="repeat.type === SequenceType.Circle">
                <label :for="'repeat'+k">Radius</label>
                <InputNumber :id="'radius'+k" v-model="repeat.radius" showButtons class="p-inputgroup"/>
              </div>

              <div class="flex flex-shrink-0 flex-column m-2">
                <!--              <label v-if="repeat.type === SequenceType.Line" :for="k+'axis'">Orientation</label>-->
                <!--              <label v-else :for="k+'axis'"> clockwise </label>-->
                <label :for="k+'axis'">Axis</label>
                <SelectButton :allowEmpty=false
                              v-model="repeat.axis" :options="[SequenceAxis.Horizontal, SequenceAxis.Vertical]"
                              :option-label="axisName(repeat)"
                              :inputId="k+'axis'"/>
              </div>

              <div class="flex flex-shrink-0 flex-column m-2">
                <label v-if="repeat.type === SequenceType.Line" :for="k+'dir'" class="ml-2">
                  Direction
                </label>
                <label v-else :for="k+'dir'" class="ml-2"> backs facing</label>

                <SelectButton
                    :allowEmpty=false
                    v-model="repeat.dir"
                    :options="[SequenceDirection.Positive, SequenceDirection.Negative]"
                    :option-label="directionName(repeat)"
                    :inputId="k+'dir'"
                />
              </div>

              <div class="flex-row m-2">
                <label :for="'repeat'+k">Repeat for #times</label>
                <InputNumber
                    class="p-inputgroup"
                    id="'repeat'+k"
                    show-buttons
                    v-model="repeat.num"/>
              </div>

              <div class="flex-row m-2">
                <label :for="'separation'+k">Separation</label>

                <div class="p-inputgroup">
                <span v-if="repeat.type ===  SequenceType.Circle" class="p-inputgroup-addon">
                    <Checkbox v-model="repeat.equivalentSpaced" :binary="true"/>
                </span>
                  <InputNumber
                      :inputId="'separation'+k"
                      :disabled="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle"
                      :modelValue="repeat.equivalentSpaced && repeat.type=== SequenceType.Circle ? 360/Math.max(1, repeat.num) : repeat.separation"
                      @update:model-value="(newValue: number) => repeat.separation = newValue"
                      showButtons
                      :suffix="repeat.type ===  SequenceType.Circle ? '°':' cm'"/>
                </div>
              </div>
              <div class="flex-row m-2">
                <label :for="'bitboxes'+k">Flip when repeats increases:</label>
                <div class="p-inputgroup">
                  <BitBoxes
                      :inputId="'bitboxes'+k"
                      :num="selectedElement!.repeats.length"
                      :options="selectedElement!.repeats.map(sequenceToString)"
                      :skip-bit="k"
                      v-model="repeat.snakeGroup"/>
                </div>
              </div>
            </div>
          </AccordionTab>
        </Accordion>
        <div v-else class="mt-2">
          Highlight an element
        </div>
      </span>
    </div>
    <div class="col-9">
      <svg style="border: 1px solid red;" id="layoutsvg" ref="svgRef" @wheel="scroll" width="100%" height="100%"
           @mousemove="maybeTranslateRotate" @mouseup="resetTranslateRotate"
           @mousedown="(e: MouseEvent) => dragStart({coord: topLeft, event: e})">
        <g ref="innerSvgRef" :transform="'scale('+scale+') translate('+topLeft.x+','+topLeft.y+')'">
          <DragNg
              v-if="selectedRoomIndex===null"
              v-for="(placement, i) in map.placements"
              :coord=placement.coord
              :transform=true
              @dragStart="moveStart"
              @rotateStart="rotateStart">
            <Room v-bind="placement.room"
                  @hoverElement="hoveringElement"
                  :sequenceKey="[KeyCategory.Placements, i]"
            />
          </DragNg>
          <Room v-else v-bind="selectedRoom!" :translating=true
                :sequenceKey="[KeyCategory.Placements, selectedRoomIndex]"
                @dragStart="moveStart"
                @rotateStart="rotateStart"
                @hoverElement="hoveringElement"
          />

          <Path
              v-if="selectedRoomIndex===null"
              v-for="p in map.paths"
              v-bind="p"
              :editing="pathMode===EditType.Path"
              :scale="1/scale"
              @dragStart="pathDrag"
          />
        </g>
        <Hatching
            v-for="colors in hatchings"
            :width="settings.areaWidth"
            :colors="colors"
            :id="colors.map((v) => v.replaceAll('#', '')).join('-')+'-hatching'"
        />
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

</style>

<script setup lang="ts">

import Hatching from "../components/Layout/Hatching.vue";

import Room from "../components/Layout/Room.vue";


import {computed, onMounted, provide, reactive, ref} from "vue";
import {onKeyStroke, useKeyModifier} from '@vueuse/core';
import {
  CoordinateInterface,
  DragStartEvent,
  ElementEvent,
  Key,
  KeyCategory,
  QualifiedKey,
  Repeats,
  RoomInterface,
  RotationCoordinateInterface,
  RotationStartEvent,
  SequenceAxis,
  SequenceDirection,
  SequenceInterface,
  SequenceType,
  Vector
} from "../types.ts";
import {teamareaStore} from "../stores/teamarea";

import {mapStore} from "../stores/map";
import DragNg from "./DragNg.vue";
import Path from "../components/Layout/Path.vue";
import DeleteButton from "../components/Layout/DeleteButton.vue";
import BitBoxes from "../components/Layout/BitBoxes.vue";
import AssignmentOverview from "../components/AssignmentOverview.vue";

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

  if (highlightedKey.value !== null) {
    const element = map.fromQualifiedKeyUpTo(highlightedKey.value, KeyCategory.Elements);
    const mouse = toInnerCoordinates(state, true)

    rotateAroundBy(element.base, mouse, state.deltaY / 10)
    state.stopImmediatePropagation()
    state.stopPropagation()
    return false
  }

  scale.value *= scrollDt

  // Calculate the new top-left offset
  const newInner = toInnerCoordinates(state)
  topLeft.x += newInner.x - innerCoords.x
  topLeft.y += newInner.y - innerCoords.y
}

function rotateAroundBy(toRotate: RotationCoordinateInterface, around: CoordinateInterface, by: number) {
  if (clamping.value > 1) {
    // Determine the difference to a clamping value
    const da = clamping.value - toRotate.rotation

    // Correct for the difference to the
    // by += da
    //   // angle + coordinateBeingRotated.value?.initialAngle must be multiple of clamping.value
    //   coordAngle = Math.round(coordAngle / clamping.value) * clamping.value
    //   angle = coordAngle - coordinateBeingRotated.value?.initialAngle
  }


  const newCoord = new Vector(toRotate.x- around.x, toRotate.y-around.y).rotate(by).add(around)

  toRotate.x = newCoord.x
  toRotate.y = newCoord.y
  toRotate.rotation += by
}

const hatchings = ref<string[][]>([])
provide('hatchingForColors', hatchingForColors)

function hatchingForColors(...colors: string[]): string {
  if (colors.length === 0) {
    colors.push('white')
  }

  colors = colors.sort()
  const exists = hatchings.value.reduce((a, v) => a || (colors.length === v.length && colors.every((value, index) => value === v[index])), false)
  if (!exists) {
    // Append
    hatchings.value.push(colors)
  }
  return `url(#${colors.map((v) => v.replaceAll("#", "")).join("-")}-hatching)`
}

function addRepeats() {
  const element = map.fromQualifiedKeyUpTo(highlightedKey.value!, KeyCategory.Elements)
  element.repeats.push(new Repeats());
}

provide('sequenceToString', sequenceToString)

function sequenceToString(r: SequenceInterface): string {
  let direction: string;
  let separation: string;
  let radius: string = "";

  if (r.type === SequenceType.Circle) {
    radius = " ⌀" + (r.radius * 2) + settings.distanceUnit
    separation = r.separation.toString();
    if (r.equivalentSpaced) {
      separation = (360 / r.num).toString()
      if (r.axis === SequenceAxis.Horizontal) {
        direction = "⥁";
      } else {
        direction = "⥀";
      }
    } else {
      if (r.axis === SequenceAxis.Horizontal) {
        direction = "↻"
      } else {
        direction = "↺"
      }
    }
    separation += "°"
  } else {
    separation = "Δ" + r.separation + settings.distanceUnit;
    if (r.axis === SequenceAxis.Horizontal) {
      if (r.dir === SequenceDirection.Positive) {
        direction = "←"
      } else {
        // Negative
        direction = "→"
      }
    } else {
      // Vertical
      if (r.dir === SequenceDirection.Positive) {
        direction = "↓"
      } else {
        // Negative
        direction = "↑"
      }
    }
  }

  return `${r.type === SequenceType.Line ? '━' : "⭘"} ${separation} ${direction}${r.num}${radius}`
}

provide("elementFromQualifiedKey", map.fromQualifiedKey)

const snapToGrid = ref(false);

enum EditType {
  Drag,
  Path,
  Select,
}

const pathMode = ref<EditType>(EditType.Drag);

const shift = useKeyModifier('Shift')
const control = useKeyModifier('Control')

const clamping = computed(() => ((shift.value || snapToGrid.value) ? 5 : 0) * ((control.value || snapToGrid.value) ? 30 : 1))

function resetTranslateRotate() {
  coordinateBeingTranslated.value = null
  coordinateBeingRotated.value = null
  offset.value = {x: 0, y: 0}
}

function maybeTranslateRotate(e: MouseEvent) {
  if (coordinateBeingRotated.value !== undefined && coordinateBeingRotated.value !== null) {
    const mouse = toInnerCoordinates(e, true)
    const mid = coordinateBeingRotated.value?.middle
    let angle = new Vector(mouse.x - mid.x, mouse.y - mid.y).asAngle() - coordinateBeingRotated.value?.angle

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

function pathDrag(e: DragStartEvent) {
  if (pathMode.value !== EditType.Path) {
    return
  }

  dragStart(e)
}

function moveStart(e: DragStartEvent) {
  if (pathMode.value !== EditType.Drag) {
    return
  }

  dragStart(e)
}

function dragStart(e: DragStartEvent) {
  const innerCoords = toInnerCoordinates(e.event)
  if (pathMode.value == EditType.Path && e.event.ctrlKey) {

    // startCoord
    const startCoord = toInnerCoordinates(e.event, true)

    const path = reactive({
      start: {x: startCoord.x, y: startCoord.y},
      end: {x: startCoord.x, y: startCoord.y}
    })

    if (selectedRoom.value) {
      selectedRoom.value?.paths.push(path)
    } else {
      map.paths.push(path)
    }

    coordinateBeingTranslated.value = path.end;
    offset.value = {x: e.coord.x, y: e.coord.y}

    return
  } else {
    coordinateBeingTranslated.value = e.coord
  }

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
  if (pathMode.value !== EditType.Drag) {
    return
  }

  const rotCoord = e.coord as RotationCoordinateInterface

  // Maybe use these instead of the box size to calculate where someone clicked
  // const handleCoords = toInnerCoordinates(e.event)

  // Derivation starts at the coordinate itself
  //  1. derives vector from coordinate to middle
  //  2. derives angle from middle to handle
  //  3. derives inverse of vector from coordinate to middle
  //     i.e. from middle to coordinate by flipping signs
  // const toMiddle = new Vector(e.width / 2 + e.x, e.height / 2 + e.y).rotate(rotCoord.rotation)
  const toMiddle = new Vector(e.width / 2, e.height / 2).rotate(rotCoord.rotation)

  const middle = new Vector(e.width / 2 + e.x, e.height / 2 + e.y);//toMiddle.copy().add(e.coord)
  const toRoot = toMiddle.copy().multiply(-1)

  randomCoord.value = middle

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

const randomCoord = ref<CoordinateInterface | null>(null)

const topLeft = reactive<CoordinateInterface>({x: 0, y: 0});
const offset = ref<CoordinateInterface>({x: 0, y: 0});
const coordinateBeingTranslated = ref<CoordinateInterface | null>();
const coordinateBeingRotated = ref<RotationAroundInterface | null>();
const scale = ref<number>(1)
provide('scale', scale)


const svgRef = ref()
const innerSvgRef = ref()

function toInnerCoordinates(e: MouseEvent, translate: boolean = false): CoordinateInterface {
  const br = svgRef.value.getBoundingClientRect();
  return {
    x: (e.clientX - br.left) / scale.value - (translate ? topLeft.x : 0),
    y: (e.clientY - br.top) / scale.value - (translate ? topLeft.y : 0),
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

onMounted(rescale);
onMounted(() => {
  window.setTimeout(() => {
    selectedRoomIndex.value = 0
    // selectedElementIndex.value = 0
  }, 200)
})

function rescale() {
  window.setTimeout(rescaleT, 0)
}

function rescaleT() {
  const padding = 40
  const contentBox = svgRef.value.getBBox({
    stroke: true,
  })

  const svgBox = svgRef.value.getBoundingClientRect()
  const width = svgBox.width - 2 * padding
  const height = svgBox.height - 2 * padding
  const oldScale = scale.value

  const scaleX = width / (contentBox.width / oldScale)
  const scaleY = height / (contentBox.height / oldScale)

  topLeft.x = -contentBox.x / scale.value + topLeft.x
  topLeft.y = -contentBox.y / scale.value + topLeft.y

  // Align the map to the center
  if (scaleX <= scaleY) {
    scale.value = scaleX;
    topLeft.x += padding / scaleX
    topLeft.y += (svgBox.height / scaleX - contentBox.height / oldScale) / 2;
  } else {
    scale.value = scaleY;
    topLeft.x += (svgBox.width / scaleY - contentBox.width / oldScale) / 2;
    topLeft.y += padding / scaleY
  }
}

const highlightedKey = ref<QualifiedKey | null>(null)

function resetHover() {
  highlightedKey.value = null
}

onKeyStroke('Escape', resetHover)
onKeyStroke('Tab', (e) => {
  if (nextElement(e.shiftKey ? -1 : 1)) {
    e.preventDefault()
  }
});

function nextElement(dir: number) {
  if (highlightedKey.value === null) {
    return false;
  }

  const k = [...highlightedKey.value];

  // TODO select next element
  let cat: Key = -1;
  let i: int;
  for (i = k.length - 1; i >= 0; i--) {
    if (typeof (k[i]) !== "number") {
      cat = k[i]
      break
    }
  }

  switch (cat) {
    case KeyCategory.Repeats:
      const elementKey = k.toSpliced(k.indexOf(KeyCategory.Elements) + 2)
      const element = map.fromQualifiedKey(elementKey);
      if (!element) {
        return false;
      }

      for (let atRepeat = 0; atRepeat < element.repeats.length; atRepeat++) {
        k[i + 1 + atRepeat] += dir

        // Stop when withing bounds
        if (k[i + 1 + atRepeat] >= 0 && k[i + 1 + atRepeat] < element.repeats[atRepeat].num) {
          break
        }

        k[i + 1 + atRepeat] = dir > 0 ? 0 : element.repeats[atRepeat].num - 1;
      }

      highlightedKey.value = k;

      // TODO move to the next element

      return true;
  }
}

provide("highlightedKey", highlightedKey)


const selectedRoomIndex = ref<number | null>(null);
const selectedRoom = computed<RoomInterface | null>(() => {
  return selectedRoomIndex.value === null
      ? null
      : map.placements[selectedRoomIndex.value].room
})

// const selectedElementIndex = ref<number | null>(null);
const selectedElementIndex = computed<number | null>(() => {
  if (highlightedKey.value === null) {
    return null
  }

  const index = highlightedKey.value.indexOf(KeyCategory.Elements)
  if (index < 0 || highlightedKey.value?.length < index + 1) {
    return null;
  }

  return highlightedKey.value[index + 1];
});
const selectedElement = computed(() => {
  return selectedElementIndex.value === null || selectedRoomIndex.value === null
      ? null
      : map.placements[selectedRoomIndex.value].room.elements[selectedElementIndex.value]
})

function hoveringElement(e: ElementEvent) {
  highlightedKey.value = e.key
}

function axisName(s: SequenceInterface): (t: SequenceAxis) => string {
  return (t: SequenceAxis) => {
    switch (s.type) {
      case SequenceType.Line:
        switch (t) {
          case SequenceAxis.Horizontal:
            return "Horizontal"
          case SequenceAxis.Vertical:
            return "Vertical"
        }
        break
      case SequenceType.Circle:
        switch (t) {
          case SequenceAxis.Horizontal:
            return "Clockwise"
          case SequenceAxis.Vertical:
            return "Counterclockwise"
        }
        break
    }

    return 'unknown'
  }
}

function directionName(s: SequenceInterface): (t: SequenceDirection) => string {
  return (t: SequenceDirection) => {
    switch (s.type) {
      case SequenceType.Line:
        switch (s.axis) {
          case SequenceAxis.Horizontal:
            switch (t) {
              case SequenceDirection.Positive:
                return "Left"
              case SequenceDirection.Negative:
                return "Right"
            }
          case SequenceAxis.Vertical:
            switch (t) {
              case SequenceDirection.Positive:
                return "Behind"
              case SequenceDirection.Negative:
                return "In front"
            }
        }
      case SequenceType.Circle:
        switch (t) {
          case SequenceDirection.Positive:
            return "Backs"
          case SequenceDirection.Negative:
            return "Fronts"
        }
    }
  }

  return 'unknown'
}

</script>
