import { base } from './base';

export type Task = base & {
    title?: string
    short_id?: string
    subject?: string
    content?: string
    promotion?: number
    class?: string
    region?: string
    semester?: string
    global?: boolean
    members?: Array<string>
    created_by?: string
    updated_by?: string
    due_date?: Date
}