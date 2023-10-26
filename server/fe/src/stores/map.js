import {defineStore} from 'pinia'

export const mapStore = defineStore({
    id: 'map',
    getters: {},
    state: () => ({
        rooms: [
            {
                name: "Room",
                type: "Rect",
                coords: [[0, 0], [10, 0], [10, 10], [0, 10]],
                sequences: [
                    {
                        name: "Column",
                        base: [0,0],
                        repeats: [
                            {
                                type: "Line",
                                axis: true,
                                dir: false,
                            },
                            {
                                type: "Circle",
                                axis: true,
                                dir: false,
                            }
                        ]
                    }
                ]
            },

        ]


        maxTeamLength: 0,

        seatSep: 3,
        seatDist: 3,
        seatNum: 3,

        marginX: 5,
        marginY: 5,

        areaWidth: 50,
        areaHeight: 30,

        height: 15,
        seatHeight: 3,
        rotation: 0,
        offsetX: 0,
        offsetY: 0,
    }),
    actions: {},
})
