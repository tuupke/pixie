import {defineStore} from 'pinia'
import {
    CoordinateInterface,
    ElementInterface,
    PathInterface,
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
            paths: [
                // {
                //     start: {x: 0, y: 0},
                //     end: {x: 1000, y: 1000}
                // }
            ] as PathInterface[],
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
                            new Repeats(SequenceType.Line, 2, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 500, true),
                            new Repeats(SequenceType.Line, 3, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 1500, true),
                            new Repeats(SequenceType.Line, 4, SequenceAxis.Vertical, SequenceDirection.Negative, 0, 1000, true),
                            // new Repeats(SequenceType.Circle, 4, SequenceAxis.Horizontal, SequenceDirection.Negative , 400, 0, true),
                        ]
                    } as ElementInterface,
                        // {
                        //     base: {x: 100, y: 400, rotation: 0},
                        //     repeats: []
                        // } as ElementInterface
                    ] as ElementInterface[],
                    paths: [] as PathInterface[],
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
