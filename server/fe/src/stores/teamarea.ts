import {defineStore} from 'pinia'
import {RotationCoordinateInterface, Vector} from "../types.ts";

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

            return Math.min(this.tableHeight - 70 + this.strokeWidth, xWidth)
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
            areaOffsetY: 50 as number,

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

            distanceUnit: "cm" as string,

            PathAttachDistance: 100,
            PathDetectDistance: 110,
        }
    },
    actions: {
        registerTeamId(l: string): void {
            this.maxTeamLength = Math.max(Math.floor(l.length), this.maxTeamLength);
        },

        offset(coord: RotationCoordinateInterface, distX: number, distY: number, relativeToEdge: boolean = false): RotationCoordinateInterface {

            const relEdge = relativeToEdge ? 1 : 0
            const coordX: number = (50 - this.areaOffsetX) / 100 * this.areaWidth +
                distX +
                relEdge * Math.sign(distX) * this.areaWidth / 2

            const coordY: number =
                (50 - this.areaOffsetY) / 100 * this.areaHeight +
                distY +
                relEdge *  Math.sign(distY) * this.areaHeight / 2

            const vect = new Vector(coordX, coordY, coord.rotation).add(coord)
            return {
                x: vect.x,
                y: vect.y,
                rotation: coord.rotation,
            }
        },

        toMiddle(coord: RotationCoordinateInterface): RotationCoordinateInterface {
            const x: number = (50 - this.areaOffsetX) / 100 * this.areaWidth
            const y: number = (50 - this.areaOffsetY) / 100 * this.areaHeight
            const vect = new Vector(x, y, coord.rotation).add(coord)
            return {
                x: vect.x,
                y: vect.y,
                rotation: coord.rotation,
            }
        }
    },
})
