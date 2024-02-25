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
    base: RotationCoordinate
    repeats: SequenceInterface[]
}

export interface RoomInterface {
    name: string
    outline: Coordinate[]
    elements: ElementInterface[]
}

export interface Coordinate {
    x: number
    y: number
}

export type RotationCoordinate = Coordinate & {
    rotation: number
}

export interface DragStartEvent {
    coord: Coordinate | RotationCoordinate
    event: MouseEvent
}
