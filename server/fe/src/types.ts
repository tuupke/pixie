export enum SequenceType {
    Line = "line",
    Circle = "circle",
}

export enum SequenceAxis {
    Vertical = "vertical",
    Horizontal = "horizontal",
}

export enum SequenceDirection {
    Positive = "positive",
    Negative = "negative",
}

export interface SequenceInterface {
    type: SequenceType
    num: number
    axis: SequenceAxis
    dir: SequenceDirection
    radius: number,
    separation: number,
    equivalentSpaced: boolean,
    ascending: boolean,
    snakeGroup: number | null,
}

export enum KeyCategory {
    Table = "table",
    Path = "path",
    Wall = "wall",
    Repeats = "repeats",
    Repeat = "repeat",
    Placement = "placement",
    Placements = "placements",
    Elements = "elements",
    Room = "room",
    Abs = "absolute",
    Top = "top",
    Bottom = "bottom",
}

export type Key = KeyCategory | number
export type QualifiedKey = Key[]

export interface ElementInterface {
    base: RotationCoordinateInterface
    repeats: SequenceInterface[]
}

export interface RoomInterface {
    name: string
    outline: CoordinateInterface[]
    elements: ElementInterface[]
    overrides: [QualifiedKey, number | null][],
    paths: PathInterface[]
}

export interface PathCoordinatesInterface {
    start: CoordinateInterface,
    end: CoordinateInterface
}

export interface CoordinateInterface {
    x: number
    y: number
}

export type RotationCoordinateInterface = CoordinateInterface & {
    rotation: number
}

export interface ElementEvent {
    key: QualifiedKey
    event: MouseEvent
}

export interface DragStartEvent {
    coord: CoordinateInterface | RotationCoordinateInterface
    event: MouseEvent
}

export interface BoxInterface {
    width: number
    height: number
    x: number
    y: number
}

export type areaKey = number[]

export type PathPos = areaKey | CoordinateInterface

export interface PathSpec {
    start: PathPos
    end: PathPos
}

export interface PathInterface {
    start: QualifiedKey
    end: QualifiedKey
}

export type RotationStartEvent = DragStartEvent & BoxInterface

export class Vector implements CoordinateInterface {
    x: number = 0
    y: number = 0

    constructor(x: number, y: number, rotate: number = 0) {
        this.x = x
        this.y = y

        if (rotate !== 0) {
            this.rotate(rotate)
        }
    }

    rotate(angle: number): Vector {
        angle = angle * (Math.PI / 180);
        const cos = Math.cos(angle);
        const sin = Math.sin(angle);

        const ox = this.x
        const oy = this.y

        this.x = ox * cos - oy * sin
        this.y = ox * sin + oy * cos

        return this
    }

    magnitude(): number {
        const xx = this.x * this.x
        const yy = this.y * this.y

        return Math.sqrt(xx + yy)
    }

    distance(vect: CoordinateInterface): number {
        const x = vect.x - this.x;
        const y = vect.y - this.y;
        const xx = x * x
        const yy = y * y

        return Math.sqrt(xx + yy)

    }

    add(vect: CoordinateInterface): Vector {
        this.x += vect.x
        this.y += vect.y

        return this
    }

    copy(): Vector {
        return new Vector(this.x, this.y)
    }

    normalize(): Vector {
        return this.multiply(1 / this.magnitude())
    }

    multiply(n: number): Vector {
        this.x *= n
        this.y *= n

        return this
    }

    multiplyVector(n: CoordinateInterface): Vector {
        this.x *= n.x
        this.y *= n.y

        return this
    }

    // asAngle returns the angle of this vector relative to the x-axis in 'SVG degrees' (counterclockwise)
    asAngle(): number {
        return this.angleWith(new Vector(1, 0))
    }

    angleWith(v: CoordinateInterface): number {
        const otherVect = Math.atan2(v.y, v.x)
        return (Math.atan2(this.y, this.x) - otherVect) * 180 / Math.PI
    }
}

export interface RotateVector {
  v: Vector;
  rotateBy: number;
}

export class Repeats implements SequenceInterface {
    constructor(
        public type: SequenceType = SequenceType.Line,
        public num: number = 1,
        public axis: SequenceAxis = SequenceAxis.Horizontal,
        public dir: SequenceDirection = SequenceDirection.Positive,
        public radius: number = 0,
        public separation: number = 0,
        public equivalentSpaced: boolean = true,
        public ascending: boolean = true,
        public snakeGroup: number | null = null,
    ) {}

    trueSeparation(): number {
        return (this.type !== SequenceType.Circle || !this.equivalentSpaced)
            ? this.separation : 360 / Math.max(1, this.num)
    }

    lineVect(): RotateVector {
        // Assume horizontal movement
        let baseVect = new Vector(this.trueSeparation(), 0)
        if (this.axis === SequenceAxis.Vertical) {
            baseVect = new Vector(0, -this.trueSeparation())
        }

        if (this.dir === SequenceDirection.Positive) {
            baseVect = baseVect.multiply(-1)
        }

        // TODO, original returned v as baseVect.rotate(props.rotation). But it should be 0 right? This has been 'pushed down' into calculateCoordinate.
        return {v: baseVect, rotateBy: 0}
    }

    circleVect(): RotateVector {
        const orientation = this.axis == SequenceAxis.Vertical ? 1 : -1
        const direction = this.dir == SequenceDirection.Positive ? 1 : -1

        const angle = orientation*this.trueSeparation()
        const half_angle = angle/2

        const magnitude = 2*this.radius*Math.sin(half_angle*(Math.PI / 180))

        return {
            v: new Vector(0, direction).multiply(magnitude).rotate((90-half_angle)),
            rotateBy: -angle,
        }
    }

    calculateCoordinate(current: RotationCoordinateInterface, index: number): RotationCoordinateInterface{
        return this.calculateCoordinates(current, index, index+1)[0]
    }

    calculateAllCoordinates(current: RotationCoordinateInterface): RotationCoordinateInterface[] {
        return this.calculateCoordinates(current, 0, this.num)
    }

    calculateCoordinates(current: RotationCoordinateInterface, start: number, end: number): RotationCoordinateInterface[] {
        const baseVect: RotateVector = (this.type === SequenceType.Line
            ? this.lineVect()
            : this.circleVect());

        baseVect.v.rotate(current.rotation);

        let currentVector = new Vector(current.x, current.y)
        let currentRotation = current.rotation
        let coordinates = Array<RotationCoordinateInterface>()

        for (let i = 0; i < end; i++) {
            if (i >= start) {
                coordinates.push({
                    x: currentVector.x,
                    y: currentVector.y,
                    rotation: currentRotation
                })
            }

            // Update to where it should point now. Since baseVect might be 'circular' it
            currentVector.add(baseVect.v)
            baseVect.v.rotate(baseVect.rotateBy)
            currentRotation+=baseVect.rotateBy
        }

        return coordinates
    }
}
