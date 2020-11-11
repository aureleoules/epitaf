export interface Cours {
    name: string;
    start_date: Date;
    end_date: Date;
    duration: number;
    staff: string[];
    rooms: string[];
}

export interface Day {
    date: Date;
    courses: Cours[];
}

export interface Calendar {
    days: Day[];
}