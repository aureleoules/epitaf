export type StaffList = {
    Id: number;
    Name: string;
    Type: number;
    ParentId: number;
}

export type RoomList = {
    Id: number;
    Name: string;
    ParentId: number;
    Rooms?: any;
    Type: number;
}

export type GroupList = {
    Groups?: any;
    Id: number;
    Name: string;
    Type: number;
    ParentId: number;
}

export type CourseList = {
    Id: number;
    Name: string;
    BeginDate: Date;
    EndDate: Date;
    Duration: number;
    StaffList: StaffList[];
    RoomList: RoomList[];
    GroupList: GroupList[];
    Code?: any;
    Type?: any;
    Url?: any;
    Info?: any;
}

export type DayList = {
    CourseList: CourseList[];
    DateTime: Date;
}

export type Calendar = {
    Id: number;
    DayList: DayList[];
}