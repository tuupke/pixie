import {defineStore} from 'pinia'

export const teamtableStore = defineStore({
    id: 'teamtable',
    getters: {
        areaX () { return -this.areaWidth / 2 },
        areaY () { return -this.areaHeight / 2 },
        tableX () { return this.areaX+this.marginX + this.offsetX },
        tableY () { return this.areaY+this.marginY + this.offsetY },

        tableWidth() { return this.areaWidth - 2*this.marginX},
        tableHeight() { return this.areaHeight - 2*this.marginY - (this.seatNum > 0 ? 1 : 0) * (this.seatHeight + this.seatDist) },

        seatWidth() { return this.tableWidth / this.seatNum - this.seatSep },
        seatY() { return this.tableY + this.tableHeight+this.seatDist },

        dFontSize() {
            let xWidth = this.tableWidth / this.maxTeamLength;
            if (xWidth <= 0) {
                xWidth = 2000
            }

            return Math.min(this.tableHeight, xWidth)
        }
    },
    state: () => ({
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
    actions: {
        registerTeamId(l) {
            this.maxTeamLength = Math.max(Math.floor(Math.log10(l)), this.maxTeamLength);
        }
    },
})
