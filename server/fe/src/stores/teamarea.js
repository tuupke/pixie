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
            return this.areaX + this.areaPaddingX - this.tableOffsetX
        },
        tableY() {
            return this.areaY + this.areaPaddingY - this.tableOffsetY
        },

        tableWidth() {
            return this.areaWidth - 2 * this.areaPaddingX
        },
        tableHeight() {
            return this.areaHeight - 2 * this.areaPaddingY - (this.seatNum > 0 ? 1 : 0) * (this.seatHeight + this.seatDist)
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

        areaWidth: 50,
        areaHeight: 120,

        areaPaddingX: 0,
        areaPaddingY: 0,

        tableOffsetX: 0,
        tableOffsetY: 0,

        maxTeamLength: 0,

        seatSep: 3,
        seatDist: 3,
        seatNum: 3,
        seatHeight: 3,
        seatPadding: 0,

        distanceUnit: "dm"
    }),
    actions: {
        registerTeamId(l) {
            this.maxTeamLength = Math.max(Math.floor(Math.log10(l)), this.maxTeamLength);
        }
    },
})
