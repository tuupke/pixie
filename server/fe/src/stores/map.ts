import {defineStore} from 'pinia'
import {
    CoordinateInterface,
    ElementInterface,
    Repeats,
    RoomInterface,
    RotationCoordinateInterface,
    SequenceAxis,
    SequenceDirection,
    SequenceType,
} from "../types.ts";

export const mapStore = defineStore({
    id: 'map',
    getters: {},
    state: () => {
        return {
            placements: [{
                coord: {
                    x: 0,
                    y: 0,
                    rotation: 0,
                } as RotationCoordinateInterface,
                room: {
                    name: "Main room",
                    outline: [
                        {x: -100, y: -100},
                        {x: 100, y: -100},
                        {x: 100, y: 100},
                        {x: -100, y: 100}
                    ] as CoordinateInterface[],
                    elements: [{
                        base: {x: 0, y: 0, rotation: 0},
                        repeats: [
                            new Repeats(SequenceType.Line, 3, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 1000, true),
                            new Repeats(SequenceType.Circle, 2, SequenceAxis.Horizontal, SequenceDirection.Positive, 500, 0, true),
                        ]
                    } as ElementInterface,
                        {
                            base: {x: 100, y: 400, rotation: 0},
                            repeats: []
                        } as ElementInterface
                    ] as ElementInterface[]
                } as RoomInterface,
            }] as RoomPlacement[]
        }
    },
    actions: {},
})

interface RoomPlacement {
    room: RoomInterface
    coord: RotationCoordinateInterface
}
