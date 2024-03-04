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
}

export interface ElementInterface {
    base: RotationCoordinateInterface
    repeats: SequenceInterface[]
}

export interface RoomInterface {
    name: string
    outline: CoordinateInterface[]
    elements: ElementInterface[]
}

export interface CoordinateInterface {
    x: number
    y: number
}

export type RotationCoordinateInterface = CoordinateInterface & {
    rotation: number
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

export type RotationStartEvent = DragStartEvent & BoxInterface

export class Vector implements CoordinateInterface {
    x: number = 0
    y: number = 0

    constructor(x: number, y: number) {
        this.x = x
        this.y = y
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

export class Repeats implements SequenceInterface {
    constructor(
        public type: SequenceType,
        public num: number,
        public axis: SequenceAxis,
        public dir: SequenceDirection,
        public radius: number,
        public separation: number,
        public equivalentSpaced: boolean
    ) {
    }

    axisName(): string {
        switch (this.type) {
            case SequenceType.Line:
                switch (this.axis) {
                    case SequenceAxis.Horizontal:
                        return "Horizontal"
                    case SequenceAxis.Vertical:
                        return "Vertical"
                }
                break
            case SequenceType.Circle:
                switch (this.axis) {
                    case SequenceAxis.Horizontal:
                        return "Clockwise"
                    case SequenceAxis.Vertical:
                        return "Counterclockwise"
                }
                break
        }

        return 'unknown'
    }


    directionName(): string {
        switch (this.type) {
            case SequenceType.Line:
                switch (this.axis) {
                    case SequenceAxis.Horizontal:
                        switch (this.dir) {
                            case SequenceDirection.Positive:
                                return "left"
                            case SequenceDirection.Negative:
                                return "right"
                        }
                        break
                    case SequenceAxis.Vertical:
                        switch (this.dir) {
                            case SequenceDirection.Positive:
                                return "behind"
                            case SequenceDirection.Negative:
                                return "in front"
                        }
                        break
                }
                break
            case SequenceType.Circle:
                switch (this.dir) {
                    case SequenceDirection.Positive:
                        return "backs"
                    case SequenceDirection.Negative:
                        return "fronts"
                }
        }

        return 'unknown'
    }


}
