
enum RepeatTypes {
    Unknown,
    Line,
    Circle,
}

type Position = {x: number, y: number, rot: number};

type Repeat = {
    Type: RepeatTypes;
    Position: Position;
    Num: number;

    Axes: [boolean, boolean];
    Separation: number;
    Repeats: number;
    EquivalentSpacing: boolean;
}

type Placement = Position | {
    Repeats: Repeat[];
}

type Table = {
    Seats: {
        Sep: Number;
        Dist: Number;
        Num: Number;
        Height: Number;
    }

    Margin: {
        X: Number;
        Y: Number;
    }

    OffSet: {
        X: Number;
        Y: Number;
    }

    Area: {
        Width: Number;
        Height: Number;
    }
}

class Room {
    Name: string = "";
    Shape: RepeatTypes = RepeatTypes.Unknown;
    Placements: Placement[] = [];
}