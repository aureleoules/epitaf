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
    visibility?: string
    members?: Array<string>
    completed?: boolean
    completed_at?: Date
    created_by?: string
    updated_by?: string
    created_by_login?: string
    updated_by_login?: string
    due_date?: Date
}