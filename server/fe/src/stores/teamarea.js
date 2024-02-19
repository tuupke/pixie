import {defineStore} from 'pinia'

export const teamareaStore = defineStore({
    id: 'teamarea',
    getters: {
        areaX() {
            return -this.areaWidth*this.areaOffsetX/100;
        },
        areaY() {
            return -this.areaHeight*this.areaOffsetY/100;
        },
        tableX() {
            return this.areaX + this.areaPaddingX
        },
        tableY() {
            return this.areaY + this.areaPaddingY // - this.tableOffsetY
        },

        seatY() {
            return this.tableY + this.tableHeight + this.seatDist
        },

        dFontSize() {
            let xWidth = this.tableWidth / this.maxTeamLength;
            if (xWidth <= 0) {
                xWidth = 2000
            }

            return Math.min(this.tableHeight-3, xWidth)
        },

        seatWidth() {
            return (
                this.tableWidth -
                (this.seatNum + 1) * this.seatPadding -
                (this.seatNum - 1) * this.seatSep
            ) / this.seatNum
        },
    },
    state: () => ({
        areaOffsetX: 50,
        areaOffsetY: 0,

        areaWidth: 500,
        areaHeight: 300,

        areaPaddingX: 0,
        areaPaddingY: 0,

        tableHeight: 240,
        tableWidth: 500,
        tableOffsetX: 0,
        tableOffsetY: 0,

        maxTeamLength: 0,

        seatSep: 30,
        seatDist: 30,
        seatNum: 3,
        seatHeight: 30,
        seatPadding: 0,

        distanceUnit: "cm"
    }),
    actions: {
        registerTeamId(l) {
            this.maxTeamLength = Math.max(Math.floor(Math.log10(l)), this.maxTeamLength);
        }
    },
})
