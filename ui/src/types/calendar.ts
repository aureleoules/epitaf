export interface Room {
    id: number;
    capacity: number;
    name: string;
    idRoomType: number;
    idLocation: number;
}

export interface Group {
    id: number;
    idParent?: any;
    name: string;
    path?: any;
    count?: any;
    isReadOnly?: any;
    idSchool: number;
    color: string;
}

export interface Calendar {
    idReservation: number;
    idCourse?: any;
    name: string;
    idType: number;
    startDate: Date;
    endDate: Date;
    isOnline: boolean;
    rooms: Room[];
    groups: Group[];
    teachers?: any;
    idSchool: number;
    schoolName: string;
}