import {defineStore} from 'pinia'
import {
    CoordinateInterface,
    ElementInterface,
    Key,
    KeyCategory,
    PathInterface,
    QualifiedKey,
    Repeats,
    RoomInterface,
    RotationCoordinateInterface,
    SequenceAxis,
    SequenceDirection,
    SequenceInterface,
    SequenceType,
} from "../types.ts";
import {Trie} from "../trie.ts";

export interface tableAssignment {
    key: QualifiedKey
    num: number
    ignored: boolean
    duplicate: boolean
}

export const mapStore = defineStore({
    id: 'map',
    state: () => {
        return {
            paths: [] as PathInterface[],
            placements: [{
                coord: {
                    x: 0,
                    y: 0,
                    rotation: 0,
                } as RotationCoordinateInterface,
                room: {
                    name: "Main room",
                    outline: [
                        {x: -100, y: -100},
                        {x: 100, y: -100},
                        {x: 100, y: 100},
                        {x: -100, y: 100}
                    ] as CoordinateInterface[],
                    paths: [
                        {
                            end: [KeyCategory.Elements, 1, KeyCategory.Repeats, 1, KeyCategory.Top],
                            start: [KeyCategory.Elements, 0, KeyCategory.Repeats, 0, 1, 2, KeyCategory.Top]
                        }, {
                            start: [KeyCategory.Elements, 0, KeyCategory.Repeats, 1, 2, 3, KeyCategory.Bottom],
                            end: [KeyCategory.Elements, 0, KeyCategory.Repeats, 0, 2, 3, KeyCategory.Top]
                        }
                    ],
                    elements: [{
                        base: {x: 0, y: 0, rotation: 0},
                        repeats: [
                            new Repeats(SequenceType.Line, 2, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 500, true, true, 1 << 2),
                            new Repeats(SequenceType.Line, 3, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 1500, true, true, 1 << 2),
                            new Repeats(SequenceType.Line, 4, SequenceAxis.Vertical, SequenceDirection.Negative, 0, 1000, true, true),
                        ],
                    } as ElementInterface,
                    {
                        base: {x: -1500, y: 800, rotation: -30},
                        repeats: [
                            new Repeats(SequenceType.Line, 2, SequenceAxis.Horizontal, SequenceDirection.Negative, 0, 500, true),
                        ],
                    } as ElementInterface
                    ] as ElementInterface[],
                    overrides: [[[ KeyCategory.Elements, 0, KeyCategory.Repeats, 0, 0, 2] as QualifiedKey, null]],
                } as RoomInterface,
            }] as RoomPlacement[]
        }
    },
    getters: {
        exceptions(): Trie<Key, number | null>{
            // Build a trie for the exceptions
            const exceptions: Trie<Key, number | null> = new Trie<Key, number | null>()

            this.placements.forEach((room, roomIndex) => {
                let baseKey: QualifiedKey = [KeyCategory.Placements, roomIndex, KeyCategory.Room];
                room.room.overrides.forEach(e => exceptions.addKey(baseKey.concat(...e[0]), e[1]))
            })
            return exceptions
        },
        assignments(): Trie<Key, tableAssignment> {
            function increment(r: SequenceInterface[], key: number[], snake: number, at: number = 0): number {
                if (at >= r.length) {
                    return snake
                }

                const isAscending = (snake: number, snakeGroup: number | null, ascending: boolean): boolean => {
                    // Calculate the number of set bits
                    let count = 0
                    let v = snake & (snakeGroup || 0)
                    while (v > 0) {
                        v &= (v - 1);
                        count++
                    }

                    // 'flip' once for every 1 in the snake.
                    return (count & 1) == 0 ? ascending : !ascending;
                }

                const repeats: SequenceInterface = r[at]

                let overflow: boolean = false
                let ascending = isAscending(snake, repeats.snakeGroup, repeats.ascending)

                if ((ascending && key[at] == repeats.num - 1) || (!ascending && key[at] == 0)) {
                    snake = increment(r, key, snake, at+1)
                    overflow = true

                    // Recalculate ascending based on new `snake` value
                    ascending = isAscending(snake, repeats.snakeGroup, repeats.ascending)
                }

                if (overflow) {
                    key[at] = ascending
                        ? 0
                        : (repeats.num - 1)
                } else if (ascending) {
                    key[at]++
                } else {
                    key[at]--
                }

                // Flip the 'snake bit'
                return snake ^ (1 << at)
            }

            const t: Trie<Key, tableAssignment> = new Trie<Key, tableAssignment>()
            let start: number = 0;
            let lastValues: QualifiedKey[] = [];
            this.placements.forEach((room, roomIndex) => {
                room.room.elements.forEach((element, elementIndex) => {
                    let baseKey: QualifiedKey = [KeyCategory.Placements, roomIndex, KeyCategory.Room, KeyCategory.Elements, elementIndex];
                    if (element.repeats.length == 0) {
                        t.addKey(baseKey, {
                            key: baseKey,
                            num: ++start,
                            ignored: false,
                            duplicate: lastValues[start-1] === undefined
                        });
                        lastValues[start-1] = baseKey

                        return
                    }

                    baseKey.push(KeyCategory.Repeats)
                    const base: number[] = element.repeats.map((n) => n.ascending ? 0 : n.num - 1);
                    let snake: number = 0;

                    const iterations = element.repeats.reduce((a, r) => a * r.num, 1)
                    for (let i = 0; i < iterations; i++) {
                        // Try and find an exception
                        let value: number | null = ++start

                        const key = baseKey.concat(...base)
                        const overrides = this.exceptions.getValue(key)
                        if (overrides !== undefined) {
                            if (overrides === null) {
                                value = overrides
                                start--;
                            } else {
                                value = overrides
                                start = overrides
                            }
                        }

                        const duplicate = lastValues[start] !== undefined
                        if (duplicate && value !== null) {
                            // Set the duplicate to false as well
                            const obj = t.getValue(lastValues[start])
                            if (obj !== undefined) {
                                obj.duplicate = true
                                t.addKey(lastValues[start], obj)
                            }
                        }

                        t.addKey(key, {
                            key: key,
                            num: start,
                            ignored: value == null,
                            duplicate: duplicate
                        });

                        if (!duplicate) {
                            lastValues[value] = key
                        }

                        snake = increment(element.repeats, base, snake)
                    }
                })
            })

            // // Try and find duplicates
            // t.forEach((value) => {
            //
            //     // value.num+=5
            // })

            return t
        }
    },
    actions: {
        keyCompare(a: QualifiedKey, b: QualifiedKey): number | null {
            if (!Array.isArray(a) || !Array.isArray(b)) {
                return null
            }

            for (let i: number = Math.min(a.length, b.length); i >= 0; i--) {
                if (a[i] === b[i]) {
                    // Nothing to do
                } else if (a[i] < b[i]) {
                    return 1;
                } else {
                    return -1;
                }
            }

            return b.length - a.length
        },
        fromQualifiedKeyUpTo(k: QualifiedKey, category: Key): any {
            const subKey = k.indexOf(category)
            if ((subKey) < 0) {
                return null
            }

            return this.fromQualifiedKey(k, k.length - subKey - 2)
        },
        fromQualifiedKeyUpToIncluding(k: QualifiedKey, category: Key): any {
            const subKey = k.indexOf(category)
            if (subKey < 0) {
                return null
            }

            return this.fromQualifiedKey(k, k.length - subKey - 1)
        },
        fromQualifiedKey(k: QualifiedKey, skipLast: number = 0): any {
            let obj: any = this
            for (let i = 0; i < k.length - skipLast; i++) {
                if (obj === undefined || !(k[i] in obj)) {
                    return null
                }

                obj = obj[k[i]]
            }

            return obj
        },

        doDelete: function (k: QualifiedKey): void {
            // Retrieve the parent of what needs to be removed.
            let deletingObject: any = this.fromQualifiedKey(k, 1)
            if (deletingObject === null || !(k[k.length - 1] in deletingObject)) {
                console.log("Object to delete cannot be found")
            }

            // Remove other relevant items.
            let deletingCategory: Key = -1;
            for (let i = k.length - 1; i >= 0; i--) {
                if (typeof (k[i]) !== "number") {
                    deletingCategory = k[i]
                    break
                }
            }

            console.log("Deleting", deletingCategory)
            if (deletingCategory == KeyCategory.Elements || deletingCategory == KeyCategory.Repeats) {
                const roomKey = k.toSpliced(k.indexOf(KeyCategory.Room) + 1)
                const paths = this.fromQualifiedKey(roomKey).paths;

                const elementIndex = k.indexOf(KeyCategory.Elements)
                const element = k[elementIndex + 1]
                for (let i = paths.length - 1; i >= 0; i--) {
                    const startIndex = paths[i].start.indexOf(KeyCategory.Elements)
                    const endIndex = paths[i].end.indexOf(KeyCategory.Elements)

                    if (element == paths[i].start[startIndex + 1] || element == paths[i].end[endIndex + 1]) {
                        paths.splice(i, 1)
                    }
                }
            }

            deletingObject.splice(k[k.length - 1], 1)
        }
    },
})

interface RoomPlacement {
    room: RoomInterface
    coord: RotationCoordinateInterface
}
