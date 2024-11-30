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

            return Math.min((this.tableHeight-this.strokeWidth), xWidth)
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
            strokeWidth: 1 as number,

            areaOffsetX: 50 as number,
            areaOffsetY: 50 as number,

            areaWidth: 180 as number,
            areaHeight: 100 as number,

            areaPaddingX: 0 as number,
            areaPaddingY: 0 as number,

            tableHeight: 60 as number,
            tableWidth: 180 as number,

            maxTeamLength: 0 as number,

            seatSep: 10 as number,
            seatDist: 10 as number,
            seatNum: 3 as number,
            seatHeight: 15 as number,
            seatPadding: 10 as number,

            distanceUnit: "cm" as string,

            PathAttachDistance: 100,
            PathDetectDistance: 100,
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
                relEdge * Math.sign(distY) * this.areaHeight / 2

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
        },

        tableOutlineAtCoord(coord: RotationCoordinateInterface, includeTop: boolean = false, includeBottom: boolean = false) {
            // Assume the table is axis-aligned.
            const westX: number = -this.areaOffsetX * this.areaWidth / 100
            const eastX: number = this.areaWidth + westX
            const southY: number = this.areaOffsetY * this.areaHeight / 100
            const northY: number = southY - this.areaHeight

            const tOffset = includeTop ? this.PathAttachDistance : 0
            const bOffset = includeBottom ? this.PathAttachDistance : 0

            // Calculate the outlines of the bounding box.
            return {
                northEast: new Vector(eastX, northY - tOffset, coord.rotation).add(coord),
                southEast: new Vector(eastX, southY + bOffset, coord.rotation).add(coord),
                southWest: new Vector(westX, southY + bOffset, coord.rotation).add(coord),
                northWest: new Vector(westX, northY - tOffset, coord.rotation).add(coord),
            }
        },

        resetToIcpcStandard() {
            // ICPC standard draws the areas around the middle.
            // All derived values are thus relative to the middle.
            this.areaOffsetX = 50 as number
            this.areaOffsetY = 50 as number

            // table width (in meters). ICPC standard is 1.8
            // table depth (in meters). ICPC standard is 0.8
            // team area width (in meters). ICPC standard is 3.0
            // team area depth (in meters). ICPC standard is 2.0

            this.areaWidth = 300
            this.areaHeight = 200

            this.tableHeight = 80
            this.tableWidth = 180

            // The table is supposed to be drawn at the middle of the area.
            this.areaPaddingX = (this.areaWidth - this.tableWidth) / 2
            this.areaPaddingY = (this.areaHeight - this.tableHeight) / 2

            // Tables are drawn 90 degrees clockwise, so width is height and vice versa.
            // The table is drawn in the middle of the area. While in pixie the table is drawn at the top.
            // Code taken from https://github.com/icpctools/icpctools/blob/main/ContestModel/src/org/icpc/tools/contest/model/FloorMap.java.
            // c = tableWidth * scale / 8f
            // x = -tableDepth * scale * 0.75f
            // y = -tableWidth * scale / 2f + c * i * 2.5f + c / 2f
            // w = tableDepth * scale * 0.25f
            // h = c * 2f

            // Transformation to pixie values yields:
            //  The seats start tableDepth/4 from the center. The table reaches until tableDepth/2 So:
            //    seatDist = tableHeight / 4
            //  The first ICPCTools-seat y-position starts at c/2 from the 'top'. So:
            //    seatPadding = tableWidth / 16
            //  h: tableHeight / 4
            //    seatHeight = tableHeight / 4
            //  w: tableWidth / 4, and subsequent seats are drawn c*2.5 apart.
            //    seatSep = tableHeight / 8 * 2.5 = tableHeight * 5 / 16
            this.seatSep = this.tableHeight * 5 / 16
            this.seatDist = this.tableHeight / 4
            this.seatNum = 3
            this.seatHeight = this.tableHeight / 4
            this.seatPadding = this.tableWidth / 16

            this.distanceUnit = "cm"
            this.PathAttachDistance = 100
            this.PathDetectDistance = 110
            this.strokeWidth = 1
        },
    },
})
