import {defineStore} from 'pinia'

export const teamareaStore = defineStore('teamarea', {
    getters: {
        areaX(): number {
            return -this.areaWidth * this.areaOffsetX / 100;
        },
        areaY(): number {
            return -this.areaHeight * this.areaOffsetY / 100;
        },
        tableX(): number {
            return this.areaX + this.areaPaddingX
        },
        tableY(): number {
            return this.areaY + this.areaPaddingY
        },

        seatY(): number {
            return this.tableY + this.tableHeight + this.seatDist
        },

        dFontSize(): number {
            let xWidth = this.tableWidth / this.maxTeamLength;
            if (xWidth <= 0) {
                xWidth = 2000
            }

            return Math.min(this.tableHeight - 70+this.strokeWidth, xWidth)
        },

        seatWidth(): number {
            return (
                this.tableWidth -
                (this.seatNum + 1) * this.seatPadding -
                (this.seatNum - 1) * this.seatSep
            ) / this.seatNum
        }
    },
    state: () => {
        return {
            strokeWidth: 3 as number,

            areaOffsetX: 50 as number,
            areaOffsetY: 0 as number,

            areaWidth: 500 as number,
            areaHeight: 300 as number,

            areaPaddingX: 0 as number,
            areaPaddingY: 0 as number,

            tableHeight: 240 as number,
            tableWidth: 500 as number,
            tableOffsetX: 0 as number,
            tableOffsetY: 0 as number,

            maxTeamLength: 0 as number,

            seatSep: 30 as number,
            seatDist: 30 as number,
            seatNum: 3 as number,
            seatHeight: 30 as number,
            seatPadding: 0 as number,

            distanceUnit: "cm" as string
        }
    },
    actions: {
        registerTeamId(l: string): void {
            this.maxTeamLength = Math.max(Math.floor(l.length), this.maxTeamLength);
        }
    },
})
