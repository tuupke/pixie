import {CoordinateInterface, PathCoordinatesInterface} from "./types.ts";

const epsilon = 0.0001

function between(a: number, b: number, c: number): boolean {
    const eps = 3;
    return a - eps <= b && b <= c + eps;
}

export interface pathInterface {
    start: CoordinateInterface
    end: CoordinateInterface
    drawStart: CoordinateInterface | null
    drawEnd: CoordinateInterface | null
    startArbitrary: boolean
    endArbitrary: boolean
    startDist: number,
    endDist: number,
    dx: number
    dy: number
    l2: number
    dist: number
}


export function pointOnPath(paths: pathInterface[], p: CoordinateInterface): boolean {
    const threshold = 40;
    for (let i in paths) {
        const path = paths[i];
        if (path.dist < epsilon) {
            return dist(path.start, p) < threshold
        }

        const t = Math.max(0, Math.min(1, ((p.x - path.start.x) * path.dx + (p.y - path.start.y) * path.dy) / path.l2));

        // line.sx + t * dx, line.sy + t * dy
        const d = dist(p, {x: path.start.x + t * path.dx, y: path.start.y + t * path.dy})
        if (d < threshold) {
            return true
        }
    }

    return false;
}

export function pathIntersection(a: PathCoordinatesInterface, b: PathCoordinatesInterface): CoordinateInterface | null {
    let x = ((a.start.x * a.end.y - a.start.y * a.end.x) * (b.start.x - b.end.x) - (a.start.x - a.end.x) * (b.start.x * b.end.y - b.start.y * b.end.x)) /
        ((a.start.x - a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x - b.end.x));
    let y = ((a.start.x * a.end.y - a.start.y * a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x * b.end.y - b.start.y * b.end.x)) /
        ((a.start.x - a.end.x) * (b.start.y - b.end.y) - (a.start.y - a.end.y) * (b.start.x - b.end.x));
    if (isNaN(x) || isNaN(y)) {
        return null;
    } else {
        if (a.start.x >= a.end.x) {
            if (!between(a.end.x, x, a.start.x)) {
                return null;
            }
        } else {
            if (!between(a.start.x, x, a.end.x)) {
                return null;
            }
        }
        if (a.start.y >= a.end.y) {
            if (!between(a.end.y, y, a.start.y)) {
                return null;
            }
        } else {
            if (!between(a.start.y, y, a.end.y)) {
                return null;
            }
        }
        if (b.start.x >= b.end.x) {
            if (!between(b.end.x, x, b.start.x)) {
                return null;
            }
        } else {
            if (!between(b.start.x, x, b.end.x)) {
                return null;
            }
        }
        if (b.start.y >= b.end.y) {
            if (!between(b.end.y, y, b.start.y)) {
                return null;
            }
        } else {
            if (!between(b.start.y, y, b.end.y)) {
                return null;
            }
        }
    }

    return {x: x, y: y};
}

export function pathIntersect(paths: pathInterface[], b: PathCoordinatesInterface): CoordinateInterface | null {
    for (let i in paths) {
        const coord = pathIntersection(paths[i], b)
        if (coord !== null) {
            return coord
        }
    }

    return null
}


export function dist(a: CoordinateInterface, b: CoordinateInterface): number {
    const dx: number = b.x - a.x;
    const dy: number = b.y - a.y;
    return Math.sqrt(dx * dx + dy * dy);
}
